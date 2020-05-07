package client

import (
	"DTCping/base/colo"
	"DTCping/base/iata"
	"DTCping/base/sort"
	"encoding/csv"
	"log"
	"os"
	"strconv"

	ping "DTCping/base/ping/base"
)

// 经测试超过IP过多会导致该库测试不准确，轻则飙高，严重导致断网。
func Pings(addrs []string, number, singleNumber, limiterNumber int, coloOpenFlag bool, filePath, iataFilePath string) {

	stMaps, err := SplitPings(addrs, number, singleNumber)
	if err != nil {
		log.Fatalln(err)
	}
	coloMaps := map[string]*colo.IpColo{}
	iataMaps := map[string]iata.Icao{}
	if coloOpenFlag {
		coloMaps = colo.FetchColo(limiterNumber, addrs)
		iataMaps, err = iata.LocalFirstGetIatas(iataFilePath)
		if err != nil {
			log.Print(err)
		}
	}

	newFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("创建文件失败\r\n")
	}
	defer newFile.Close()

	// 写入UTF-8
	_, _ = newFile.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
	// 写数据到csv文件
	writeFd := csv.NewWriter(newFile)
	header := []string{"IP", "loss%", "发送的包数", "返回的包数", "平均延迟", "MIN延迟", "MAX延迟", "COLO", "Country", "City"} //标题
	data := [][]string{
		header,
	}

	statisticsSort := sort.ProcessStatisticsSlice(stMaps)

	for i := range statisticsSort {
		coloData := func() string {
			ipColo := coloMaps[statisticsSort[i].IPAddr.String()]
			coloData := ""
			if ipColo != nil {
				coloData = ipColo.Colo
			}
			return coloData
		}()
		context := []string{
			statisticsSort[i].IPAddr.String(),
			strconv.FormatFloat(statisticsSort[i].PacketLoss, 'f', 2, 64) + "%",
			strconv.Itoa(statisticsSort[i].PacketsSent),
			strconv.Itoa(statisticsSort[i].PacketsRecv),
			statisticsSort[i].AvgRtt.String(),
			statisticsSort[i].MinRtt.String(),
			statisticsSort[i].MaxRtt.String(),
			coloData,
			iataMaps[coloData].Country,
			iataMaps[coloData].City,
		}
		data = append(data, context)
	}
	_ = writeFd.WriteAll(data)
	writeFd.Flush()
}

func SplitPings(addrs []string, number, singleNumber int) ([]*ping.Statistics, error) {
	// 预估长度 || make([]int, 0, 10) & length 0 and capacity 10 & use to append
	stMaps := make([]*ping.Statistics, 0, len(addrs))
	groupNum := len(addrs) / singleNumber
	allGroupNum := groupNum
	if singleNumber*groupNum < len(addrs) {
		allGroupNum = allGroupNum + 1
	}
	log.Printf("分次查询总次数: [%d]\r\n", allGroupNum)
	var err error
	for i := 0; i < groupNum; i++ {
		log.Printf("[%d]次号查询处理...\r\n", i+1)
		stMaps, err = pingSingle(stMaps, addrs[singleNumber*i:singleNumber*(i+1)-1], number)
		if err != nil {
			return nil, err
		}
	}
	if singleNumber*groupNum < len(addrs) {
		stMaps, err = pingSingle(stMaps, addrs[singleNumber*groupNum:], number)
		log.Printf("末次号查询处理...\r\n")
		if err != nil {
			return nil, err
		}
	}
	log.Printf("查询完毕!\r\n")
	return stMaps, nil
}

func pingSingle(stMaps []*ping.Statistics, standardAddrs []string, number int) ([]*ping.Statistics, error) {
	bp, err := ping.NewBatchPing(standardAddrs, true) // true will need to be root, false may be permission denied
	if err != nil {
		log.Fatalf("new batch ping err %v\r\n", err)
	}
	bp.Debug = false  // debug == true will fmt debug log
	bp.Count = number // if hava multi source ip, can use one isp
	bp.OnFinish = func(stMapTemp []*ping.Statistics) {
		log.Printf("%d %d %d \n", len(stMaps), len(stMapTemp), cap(stMaps))
		stMaps = append(stMaps, stMapTemp...)
	}
	err = bp.Run()
	if err != nil {
		return nil, err
	}
	bp.OnFinish(ping.BatchStatistics(bp))
	return stMaps, nil
}
