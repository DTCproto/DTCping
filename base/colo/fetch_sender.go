package colo

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

func DataSenderControl(ips []string, resBodyChan chan HttpColoByte) {
	// HTTP并发量控制
	limiter := make(chan bool, 100)
	for i := range ips {
		limiter <- true
		go dataSenderSingle(ips[i], resBodyChan, limiter)
	}
}

func dataSenderSingle(ip string, resBodyChan chan HttpColoByte, limiter chan bool) {
	// Func执行完毕需释放占用信号
	defer func() { <-limiter }()
	// 构建返回数据包
	httpColoByte := HttpColoByte{
		Ip: ip,
	}
	// 请求超时设置
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(20*time.Second, func() {
		cancel()
	})
	//构建URL
	var urlBuf bytes.Buffer
	urlBuf.WriteString("http://")
	urlBuf.WriteString(ip)
	urlBuf.WriteString("/cdn-cgi/trace")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlBuf.String(), nil)
	if err != nil {
		httpColoByte.Error = err
		resBodyChan <- httpColoByte
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		httpColoByte.Error = err
		resBodyChan <- httpColoByte
		return
	}
	defer resp.Body.Close()
	// 延时读取Response Body
	timer.Reset(15 * time.Second)

	respBody, err := ioutil.ReadAll(resp.Body)
	httpColoByte.respBody = respBody
	httpColoByte.Error = err
	resBodyChan <- httpColoByte
}
