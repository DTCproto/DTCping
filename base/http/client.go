package http

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	cfUrl    = "https://www.cloudflare.com/ips-v4"
	cfIp1111 = "1.1.1.1"
)

func GetCloudflareIps() ([]string, error) {
	return GetArrayString(cfUrl)
}

func GetArrayString(getUrl string) ([]string, error) {
	res, err := httpGet(getUrl)
	if err != nil {
		log.Println(err)
		log.Println("HTTP GET Again...")
		res, err = httpGetByCf1111(getUrl)
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return ByteToArrayString(body)
}

// 按行读取分割数据
func ByteToArrayString(body []byte) ([]string, error) {
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

func httpGet(getUrl string) (*http.Response, error) {
	return http.Get(getUrl)
}

// 临时解决方案，后续优化为作用局域内不影响全局设置[现为全局设置]
func httpGetByCf1111(getUrl string) (*http.Response, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == "www.cloudflare.com:443" {
			addr = cfIp1111 + ":443"
		}
		return dialer.DialContext(ctx, network, addr)
	}
	return http.Get(getUrl)
}
