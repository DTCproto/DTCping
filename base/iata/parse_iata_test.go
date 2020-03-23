package iata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestGetCloudflareIatas(t *testing.T) {
	iataMaps, err := GetCloudflareIatas()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s->%s\r\n", iataMaps["SJC"].Iata, iataMaps["SJC"].Country+"/"+iataMaps["SJC"].City)
	log.Printf("%s->%s\r\n", iataMaps["LAX"].Iata, iataMaps["LAX"].Country+"/"+iataMaps["LAX"].City)
	log.Printf("%s->%s\r\n", iataMaps["AMS"].Iata, iataMaps["AMS"].Country+"/"+iataMaps["AMS"].City)
	iataNum := 0
	for range iataMaps {
		iataNum = iataNum + 1
	}
	fmt.Println("IATA number is : " + strconv.Itoa(iataNum))
}

func TestGetCloudflareIatas2(t *testing.T) {
	iataMaps, err := GetCloudflareIatas()
	if err != nil {
		log.Fatal(err)
	}
	newFile, err := os.OpenFile("iatas.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("创建文件失败\r\n")
	}
	defer newFile.Close()
	resDataJson, err := json.MarshalIndent(&iataMaps, "", "  ")
	if err != nil {
		log.Fatal("序列化失败\r\n")
	}
	_, _ = newFile.Write(resDataJson)
}
