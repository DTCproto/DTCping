package colo

import (
	DTCHttp "DTCping/base/http"
	"errors"
	"log"
	"strings"
	"sync"
)

type IpColo struct {
	Ip    string
	Colo  string
	Error error
}

func FetchColo(ips []string) map[string]IpColo {
	log.Printf("COLO开始处理...总量为: [%d]\r\n", len(ips))
	var scMap sync.Map
	wg := &sync.WaitGroup{}
	limiter := make(chan bool, 256)
	for i := range ips {
		wg.Add(1)
		limiter <- true
		// start a goroutine
		go fetchColoSingle(ips[i], &scMap, limiter, wg)
	}
	wg.Wait()
	ipColoMap := make(map[string]IpColo)
	resultNum := 0
	scMap.Range(func(key, value interface{}) bool {
		resultNum = resultNum + 1
		if ip, keyOk := key.(string); keyOk {
			if ipColo, valueOk := value.(IpColo); valueOk {
				ipColoMap[ip] = ipColo
			}
		}
		return true
	})
	log.Printf("COLO处理完毕[%d]!\r\n", resultNum)
	return ipColoMap
}

func fetchColoSingle(ip string, scMap *sync.Map, limiter chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var ipColo IpColo
	ipColo.Ip = ip
	getUrl := "http://" + ip + "/cdn-cgi/trace"
	resBodyArr, err := DTCHttp.GetArrayString(getUrl)
	if err != nil {
		// 失败重试一次
		resBodyArr, err = DTCHttp.GetArrayString(getUrl)
	}
	if err != nil {
		// log.Printf("[%s]请求失败...\r\n", ip)
		// log.Print(err)
		ipColo.Error = err
		scMap.Store(ipColo.Ip, ipColo)
		<-limiter
		return
	}
	coloValueStr, err := getColo(resBodyArr)
	if err != nil {
		log.Printf("[%s]获取COLO失败...\r\n", ip)
		ipColo.Error = err
		scMap.Store(ipColo.Ip, ipColo)
		<-limiter
		return
	}
	ipColo.Colo = coloValueStr
	scMap.Store(ipColo.Ip, ipColo)
	<-limiter
}

//fl=4f197
//h=104.24.71.22
//ip=160.120.134.135
//ts=1584887600.227
//visit_scheme=http
//uag=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36
//colo=SJC
//http=http/1.1
//loc=CN
//tls=off
//sni=off
//warp=off

func getColo(bodyArr []string) (string, error) {
	for i := range bodyArr {
		ipSegTail := strings.Split(bodyArr[i], "=")
		if ipSegTail[0] == "colo" {
			return ipSegTail[1], nil
		}
		if ipSegTail[1] == "colo" {
			return ipSegTail[0], nil
		}
	}
	return "", errors.New("获取colo失败")
}
