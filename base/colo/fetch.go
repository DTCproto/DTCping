package colo

import (
	"errors"
	"log"
	"strconv"
)

func FetchColo(ips []string) map[string]IpColo {
	ipColoMap, ipErrorArr, err := FetchColoSingle(ips)
	if err != nil {
		// 存在失败数据重试一次
		log.Printf("COLO首次请求存在失败情况!总量为: [%d]...重试中...\r\n", len(ipErrorArr))
		ipColoMapRetry, _, _ := FetchColoSingle(ipErrorArr)
		for key, value := range ipColoMapRetry {
			ipColoMap[key] = value
		}
	}
	log.Printf("COLO处理完毕!\r\n")
	return ipColoMap
}

func FetchColoSingle(ips []string) (map[string]IpColo, []string, error) {
	log.Printf("COLO开始处理...总量为: [%d]\r\n", len(ips))
	// 数据处理缓存
	resBodyChan := make(chan HttpColoByte, 4096)
	// 处理完成数据接收
	ipColoChan := make(chan IpColo, 4096)
	// 建立response处理
	go DataReceiveControl(len(ips), resBodyChan, ipColoChan)
	// 发起请求
	go DataSenderControl(ips, resBodyChan)
	log.Printf("COLO本次处理总量为: [%d]\r\n", len(ips))

	var ipErrorArr []string
	ipColoMap := map[string]IpColo{}
	resSuccessNum := 0
	for i := range ips {
		ipColo := <-ipColoChan
		ipColoMap[ipColo.Ip] = ipColo
		if ipColo.Ip != "" {
			if ipColo.Error == nil {
				resSuccessNum = resSuccessNum + 1
			} else {
				ipErrorArr = append(ipErrorArr, ipColo.Ip)
			}
		}
		if i != 1 && (i+1)%500 == 0 {
			log.Printf("COLO处理中....已处理[%d]次,有效处理数据量[%d]\r\n", i+1, resSuccessNum)
		}
	}
	var ResError error
	if len(ipErrorArr) > 0 {
		ResError = errors.New("失败数据量为: " + strconv.Itoa(len(ipErrorArr)))
	}
	return ipColoMap, ipErrorArr, ResError
}
