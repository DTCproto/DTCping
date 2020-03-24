package colo

import (
	"net/http"
)

func DataSender(ips []string, resBodyChan chan HttpColoByte) {
	// HTTP并发量控制
	limiter := make(chan bool, 150)
	for i := range ips {
		limiter <- true
		go dataSenderSingle(ips[i], resBodyChan, limiter)
	}
}

func dataSenderSingle(ip string, resBodyChan chan HttpColoByte, limiter chan bool) {
	getUrl := "http://" + ip + "/cdn-cgi/trace"
	httpColoByte := HttpColoByte{
		Ip: ip,
	}
	res, err := http.Get(getUrl)
	if err != nil {
		// 失败重试一次
		res, err = http.Get(getUrl)
	}
	httpColoByte.resp = res
	httpColoByte.Error = err
	resBodyChan <- httpColoByte

	<-limiter
}
