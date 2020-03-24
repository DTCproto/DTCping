package colo

import (
	"log"
)

func FetchColo(ips []string) map[string]IpColo {
	log.Printf("COLO开始处理...总量为: [%d]\r\n", len(ips))
	// 数据处理缓存
	resBodyChan := make(chan HttpColoByte, 4096)
	// 处理完成数据接收
	ipColoChan := make(chan IpColo, 4096)
	// 建立response处理
	go DataReceive(len(ips), resBodyChan, ipColoChan)
	// 发起请求
	go DataSender(ips, resBodyChan)

	ipColoMap := make(map[string]IpColo)
	resSuccessNum := 0
	for i := range ips {
		ipColo := <-ipColoChan
		ipColoMap[ipColo.Ip] = ipColo
		if ipColo.Ip != "" && ipColo.Error == nil {
			resSuccessNum = resSuccessNum + 1
		}
		if i%199 == 0 {
			log.Printf("COLO处理中....已处理[%d]次,有效处理数据量[%d]\r\n", i+1, resSuccessNum)
		}
	}
	log.Printf("COLO处理完毕!总量为: [%d],有效处理数据量[%d]\r\n", len(ips), resSuccessNum)
	return ipColoMap
}
