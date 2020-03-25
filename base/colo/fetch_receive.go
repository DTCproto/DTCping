package colo

import (
	"errors"
	"strings"

	DTCHttp "DTCping/base/http"
)

func DataReceiveControl(LenNumber int, resBodyChan chan HttpColoByte, ipColoChan chan IpColo) {
	for i := 0; i < LenNumber; i++ {
		httpColoByte := <-resBodyChan
		go DataReceiveSingle(httpColoByte, ipColoChan)
	}
}

func DataReceiveSingle(httpColoByte HttpColoByte, ipColoChan chan IpColo) {
	ipColo := IpColo{
		Ip:    httpColoByte.Ip,
		Error: httpColoByte.Error,
	}
	if httpColoByte.Error != nil {
		ipColoChan <- ipColo
	} else {
		coloResStr, err := parsingRespToMaps(httpColoByte.respBody)
		ipColo.Error = err
		ipColo.Colo = coloResStr
		ipColoChan <- ipColo
	}
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
