package colo

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

func FetchColo(ips []string) map[string]*IpColo {
	ipColoMap, ipErrorArr, err := FetchColoSingle(defaultLimiterNumber, ips)
	if err != nil {
		// 存在失败数据重试一次
		log.Printf("COLO首次请求存在失败情况!总量为: [%d]...重试中...\r\n", len(ipErrorArr))
		ipColoMapRetry, _, _ := FetchColoSingle(defaultLimiterNumber, ipErrorArr)
		for key, value := range ipColoMapRetry {
			ipColoMap[key] = value
		}
	}
	log.Printf("COLO处理完毕!\r\n")
	return ipColoMap
}

func FetchColoSingle(limiterNum int, ips []string) (map[string]*IpColo, []string, error) {
	log.Printf("COLO开始处理...总量为: [%d]\r\n", len(ips))
	// 数据处理缓存
	resBodyChan := make(chan *HttpColoByte, defaultBodyCache)
	// 处理完成数据接收
	ipColoChan := make(chan *IpColo, defaultBodyCache)
	// 是否处理完毕接受数据
	done := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(len(ips))
	// 建立response处理
	go DataReceiveControl(done, resBodyChan, ipColoChan)
	// 发起请求
	go DataSenderControl(resBodyChan, limiterNum, ips)
	log.Printf("COLO本次处理总量为: [%d]\r\n", len(ips))
	go func() {
		wg.Wait()
		close(done)
	}()

	var ipErrorArr []string
	ipColoMap := map[string]*IpColo{}
	resSuccessNum := 0
	ipNum := 0
	for {
		select {
		case <-done:
			var ResError error
			if len(ipErrorArr) > 0 {
				ResError = errors.New("失败数据量为: " + strconv.Itoa(len(ipErrorArr)))
			}
			return ipColoMap, ipErrorArr, ResError
		case ipColo := <-ipColoChan:
			wg.Done()
			ipColoMap[ipColo.Ip] = ipColo
			if ipColo.Ip != "" {
				if ipColo.Error == nil {
					resSuccessNum = resSuccessNum + 1
				} else {
					ipErrorArr = append(ipErrorArr, ipColo.Ip)
				}
			}
			ipNum++
			if ipNum != 1 && (ipNum)%500 == 0 {
				log.Printf("COLO处理中....已处理[%d]次,有效处理数据量[%d]\r\n", ipNum, resSuccessNum)
			}
		}
	}
}
