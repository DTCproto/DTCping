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
	limiter := make(chan bool, 120)
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
	timer := time.AfterFunc(5*time.Second, func() {
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
	// 延时读取Response Body 定时器未执行时，才有延时的意义
	// 停止原定时器成功[返回true]，才做延时操作
	if !timer.Stop() {
		if timer.C != nil {
			<-timer.C
		}
	} else {
		timer.Reset(5 * time.Second)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	httpColoByte.respBody = respBody
	httpColoByte.Error = err
	resBodyChan <- httpColoByte
}
