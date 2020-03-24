package iata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const cfIataUrl = "https://raw.githubusercontent.com/mwgg/Airports/master/airports.json"

func LocalFirstGetIatas(filePath string) (map[string]Icao, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Print("读取本地Iatas文件失败，开始请求在线文件...\r\n")
		return GetCloudflareIatas()
	}
	log.Print("读取本地Iatas文件成功...\r\n")
	return formatToMaps(fileData)
}

func GetCloudflareIatas() (map[string]Icao, error) {
	resBody, err := GetDataFromUrl(cfIataUrl)
	if err != nil {
		return nil, err
	}
	return formatToMaps(resBody)
}

func GetDataFromUrl(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		// 增加一次失败重试的机会
		res, err = http.Get(url)
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func formatToMaps(resData []byte) (map[string]Icao, error) {
	resMaps := make(map[string]Icao)
	err := json.Unmarshal(resData, &resMaps)
	if err != nil {
		return nil, err
	}
	iataMaps := make(map[string]Icao)
	for _, value := range resMaps {
		// 去空格
		value.Iata = strings.Replace(value.Iata, " ", "", -1)
		if value.Iata != "" {
			iataMaps[value.Iata] = value
		}
	}
	return iataMaps, nil
}
