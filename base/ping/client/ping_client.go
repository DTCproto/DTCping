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
func Pings(addrs []string, number, singleNumber int, coloOpenFlag bool, filePath, iataFilePath string) {

	stMaps, err := SplitPings(addrs, number, singleNumber)
	if err != nil {
		log.Fatalln(err)
	}
	coloMaps := map[string]colo.IpColo{}
	iataMaps := map[string]iata.Icao{}
	if coloOpenFlag {
		coloMaps = colo.FetchColo(addrs)
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

	statisticsSort := sort.ProcessStatistics(stMaps)

	for i := range statisticsSort {
		context := []string{
			statisticsSort[i].Ip,
			strconv.FormatFloat(statisticsSort[i].Data.PacketLoss, 'f', 2, 64) + "%",
			strconv.Itoa(statisticsSort[i].Data.PacketsSent),
			strconv.Itoa(statisticsSort[i].Data.PacketsRecv),
			statisticsSort[i].Data.AvgRtt.String(),
			statisticsSort[i].Data.MinRtt.String(),
			statisticsSort[i].Data.MaxRtt.String(),
			coloMaps[statisticsSort[i].Ip].Colo,
			iataMaps[coloMaps[statisticsSort[i].Ip].Colo].Country,
			iataMaps[coloMaps[statisticsSort[i].Ip].Colo].City,
		}
		data = append(data, context)
	}
	/*for i := range stMaps {
		for ip, st := range stMaps[i] {
			//log.Printf("\n--- %s ping statistics ---\n", st.Addr)
			//log.Printf("ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n", ip,
			//	st.PacketsSent, st.PacketsRecv, st.PacketLoss)
			//log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			//	st.MinRtt, st.AvgRtt, st.MaxRtt, st.StdDevRtt)
			//log.Printf("rtts is %v \n", st.Rtts)
			context := []string{
				ip,
				strconv.FormatFloat(st.PacketLoss, 'f', 2, 64) + "%",
				strconv.Itoa(st.PacketsSent),
				strconv.Itoa(st.PacketsRecv),
				st.AvgRtt.String(),
				st.MinRtt.String(),
				st.MaxRtt.String(),
				coloMaps[ip].Colo,
				iataMaps[coloMaps[ip].Colo].Country,
				iataMaps[coloMaps[ip].Colo].City,
			}
			data = append(data, context)
		}
	}*/
	_ = writeFd.WriteAll(data)
	writeFd.Flush()
}

func SplitPings(addrs []string, number, singleNumber int) ([]map[string]*ping.Statistics, error) {
	var stMaps []map[string]*ping.Statistics
	groupNum := len(addrs) / singleNumber
	allGroupNum := groupNum
	if 256*groupNum < len(addrs) {
		allGroupNum = allGroupNum + 1
	}
	log.Printf("分次查询总次数: [%d]\r\n", allGroupNum)
	for i := 0; i < groupNum; i++ {
		log.Printf("[%d]次号查询处理...\r\n", i+1)
		stMap, err := pingSingle(addrs[singleNumber*i:singleNumber*(i+1)-1], number)
		if err != nil {
			return nil, err
		}
		stMaps = append(stMaps, stMap)
	}
	if 150*groupNum < len(addrs) {
		stMap, err := pingSingle(addrs[singleNumber*groupNum:], number)
		log.Printf("末次号查询处理...\r\n")
		if err != nil {
			return nil, err
		}
		stMaps = append(stMaps, stMap)
	}
	return stMaps, nil
}

func pingSingle(standardAddrs []string, number int) (map[string]*ping.Statistics, error) {
	bp, err := ping.NewBatchPing(standardAddrs, true) // true will need to be root, false may be permission denied
	if err != nil {
		log.Fatalf("new batch ping err %v\r\n", err)
	}
	bp.Debug = false  // debug == true will fmt debug log
	bp.Count = number // if hava multi source ip, can use one isp

	var stMapResult map[string]*ping.Statistics

	bp.OnFinish = func(stMap map[string]*ping.Statistics) {
		stMapResult = stMap
	}
	err = bp.Run()
	if err != nil {
		return nil, err
	}
	bp.OnFinish(ping.BatchStatistics(bp))
	return stMapResult, nil
}
