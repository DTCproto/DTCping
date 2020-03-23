package iata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const cfIataUrl = "https://raw.githubusercontent.com/mwgg/Airports/master/airports.json"

func GetCloudflareIatas() (map[string]Icao, error) {
	resBody, err := GetDataFromUrl(cfIataUrl)
	if err != nil {
		return nil, err
	}
	resMaps := make(map[string]Icao)
	_ = json.Unmarshal(resBody, &resMaps)
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
