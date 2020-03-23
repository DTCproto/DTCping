package http

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const cfUrl = "https://www.cloudflare.com/ips-v4"

func GetCloudflareIps() ([]string, error) {
	return GetArrayString(cfUrl)
}

func GetArrayString(getUrl string) ([]string, error) {
	res, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return byteToArrayString(body)
}

func byteToArrayString(body []byte) ([]string, error) {
	bodyAllStr := string(body)
	if bodyAllStr == "" {
		return nil, errors.New("response body 为空！\r\n")
	}
	ArrayTemp := strings.Split(bodyAllStr, "\n")
	var ArrayResult []string
	for i := 0; i < len(ArrayTemp); i++ {
		// 去除空格与换行符
		str := strings.Replace(ArrayTemp[i], " ", "", -1)
		str = strings.Replace(str, "\n", "", -1)
		if str != "" {
			ArrayResult = append(ArrayResult, str)
		}
	}
	return ArrayResult, nil
}
