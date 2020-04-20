package colo

import (
	DTCHttp "DTCping/base/http"
	"errors"
	"strings"
)

func DataReceiveControl(done chan bool, resBodyChan chan *HttpColoByte, ipColoChan chan *IpColo) {
	for {
		select {
		case <-done:
			return
		case httpColoByte := <-resBodyChan:
			go dataReceiveSingle(httpColoByte, ipColoChan)
		}
	}
}

func dataReceiveSingle(httpColoByte *HttpColoByte, ipColoChan chan *IpColo) {
	// [指针初始化 new() ]
	ipColo := &IpColo{
		Ip:    httpColoByte.Ip,
		Error: httpColoByte.Error,
	}
	if httpColoByte.Error == nil {
		coloResStr, err := parsingRespToMaps(httpColoByte.respBody)
		ipColo.Error = err
		ipColo.Colo = coloResStr
	}
	ipColoChan <- ipColo
}

func parsingRespToMaps(respBody []byte) (string, error) {
	resBodyArr, err := DTCHttp.ByteToArrayString(respBody)
	if err != nil {
		return "", err
	}
	return getColo(resBodyArr)
}

func getColo(bodyArr []string) (string, error) {
	for i := range bodyArr {
		ipSegTail := strings.Split(bodyArr[i], "=")
		if len(ipSegTail) == 2 {
			if ipSegTail[0] == "colo" {
				return ipSegTail[1], nil
			}
			if ipSegTail[1] == "colo" {
				return ipSegTail[0], nil
			}
		}
	}
	return "", errors.New("获取colo失败")
}
