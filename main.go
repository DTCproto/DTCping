package main

import (
	"DTCping/base/cidr"
	"DTCping/base/file"
	"DTCping/base/http"
	pingClient "DTCping/base/ping/client"
	"flag"
	"fmt"
	"log"
	"strconv"
)

var config struct {
	ConfigFilePath string
	StdIp          string
	StdIpSgt       string
	StdAddrName    string
	SaveFilePath   string
	PingNumber     int
	ColoOpenFlag   *bool
	IataSrcPath    string
	SingleNumber   int
	LimiterNumber  int
}

const (
	DTCPingVersion  = "v2.3.6-20200420"
	PingReference01 = "https://github.com/caucy/batch_ping"
	PingReference02 = "https://github.com/sparrc/go-ping"
	PingReference03 = "https://www.cloudflare.com/ips-v4"
	PingReference04 = "http://{ip}/cdn-cgi/trace"
	PingReference05 = "https://github.com/mwgg/Airports"
	PingReference06 = "https://raw.githubusercontent.com/mwgg/Airports/master/airports.json"
)

func init() {
	//相关说明
	fmt.Printf("当前DTCPing版本[%s]\r\n", DTCPingVersion)
	fmt.Printf("当前参考库如下: \r\n")
	fmt.Printf(PingReference01 + "\r\n")
	fmt.Printf(PingReference03 + "\r\n")
	fmt.Printf(PingReference05 + "\r\n")
	fmt.Printf("\r\n\n")

	flag.StringVar(&config.ConfigFilePath, "path", "", "config file path(./config.json)")
	flag.StringVar(&config.StdIp, "ip", "", "IP (1.0.0.1)")
	flag.StringVar(&config.StdIpSgt, "ips", "", "IP Sgt (1.0.0.0/24)")
	flag.StringVar(&config.StdAddrName, "name", "", "Addr Name (cloudflare.com)")
	flag.StringVar(&config.SaveFilePath, "s", "pingIpv4", "Save File Path (pingIpv4)")
	flag.IntVar(&config.PingNumber, "n", 10, "Ping Number (10)")
	config.ColoOpenFlag = flag.Bool("colo", false, "Colo Open Flag (default false)")
	flag.StringVar(&config.IataSrcPath, "iata", "iatas.json", "iata src file path(default iatas.json)")
	flag.IntVar(&config.SingleNumber, "esn", 256, "Each Single Number (256)")
	flag.IntVar(&config.LimiterNumber, "cln", 256, "COLO Limiter Number (256)")
}

func main() {
	flag.Parse()

	number := config.PingNumber
	filePath := config.SaveFilePath
	coloOpenFlag := *config.ColoOpenFlag
	iataSrcPath := config.IataSrcPath
	singleNumber := config.SingleNumber
	limiterNumber := config.LimiterNumber

	disTypeFlag := true
	if disTypeFlag && config.ConfigFilePath != "" {
		disTypeFlag = false
		disReadConfigFilePath(filePath, iataSrcPath, number, singleNumber, limiterNumber, coloOpenFlag)
	}
	if disTypeFlag && config.StdIp != "" {
		disTypeFlag = false
		disStdAddr(config.StdIp, filePath, iataSrcPath, number, singleNumber, limiterNumber, coloOpenFlag)
	}
	if disTypeFlag && config.StdAddrName != "" {
		disTypeFlag = false
		disStdAddr(config.StdAddrName, filePath, iataSrcPath, number, singleNumber, limiterNumber, coloOpenFlag)
	}
	if disTypeFlag && config.StdIpSgt != "" {
		disTypeFlag = false
		disStdIpSgt(config.StdIpSgt, filePath, iataSrcPath, number, singleNumber, limiterNumber, coloOpenFlag)
	}
	if disTypeFlag {
		disDefault(number, singleNumber, limiterNumber, coloOpenFlag, iataSrcPath)
	}
	// 打开文件
	//cmd := exec.Command("cmd", "/k", "start", filePath+".csv")
	//_ = cmd.Start()
	log.Print("处理完毕!\r\n")
}

func disReadConfigFilePath(filePath, iataFilePath string, number, singleNumber, limiterNumber int, coloOpenFlag bool) {
	log.Printf("READ config.ConfigFilePath\r\n")
	flags := file.PathConfig{}
	err := flags.ParseConfigFile(config.ConfigFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	if flags.CsvFileName != "" {
		filePath = flags.CsvFileName
	}
	if flags.Number != 0 {
		number = flags.Number
	}
	if flags.ColoFlag {
		coloOpenFlag = flags.ColoFlag
	}
	if flags.EachSingleNum != 0 {
		singleNumber = flags.EachSingleNum
	}
	if flags.ColoLimiterNum != 0 {
		limiterNumber = flags.ColoLimiterNum
	}
	ipSgts, err := cidr.ParseCidr(flags.IpSegments, flags.RandomValues, flags.MergeFlag, flags.ExtractFlag)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(ipSgts); i++ {
		log.Printf("[%s]开始处理[循环次数为%d]...\r\n", ipSgts[i].NodeId, number)
		log.Printf("本次ip数量为: [%d]\r\n", len(ipSgts[i].Node))
		pingClient.Pings(ipSgts[i].Node, number, singleNumber, limiterNumber, coloOpenFlag, filePath+"["+strconv.Itoa(i+1)+"].csv", iataFilePath)
		log.Printf("第[%d]批次已处理完毕![%s]\n\n", i+1, ipSgts[i].NodeId)
	}
}

func disStdAddr(stdAddr, filePath, iataFilePath string, number, singleNumber, limiterNumber int, coloOpenFlag bool) {
	log.Printf("READ STD.Ip&Addr\r\n")
	log.Printf("[%s]开始处理[循环次数为%d]...\r\n", stdAddr, number)
	ips, err := http.LookupIps(stdAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("本次ip数量为: [%d]\r\n", len(ips))
	pingClient.Pings(ips, number, singleNumber, limiterNumber, coloOpenFlag, filePath+"[ip_addr].csv", iataFilePath)
	log.Printf("单批次已处理完毕!\r\n\n")
}

func disStdIpSgt(stdIpSgt, filePath, iataFilePath string, number, singleNumber, limiterNumber int, coloOpenFlag bool) {
	log.Printf("READ STD.IpSgt\r\n")
	log.Printf("[%s]开始处理[循环次数为%d]...\r\n", stdIpSgt, number)
	ips, err := cidr.ParseCidrSingle(stdIpSgt)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("本次ip数量为: [%d]\r\n", len(ips))
	pingClient.Pings(ips, number, singleNumber, limiterNumber, coloOpenFlag, filePath+"[ip_sgt].csv", iataFilePath)
	log.Printf("单批次已处理完毕!\r\n\n")
}

func disDefault(number, singleNumber, limiterNumber int, coloOpenFlag bool, iataFilePath string) {
	log.Printf("READ Default\r\n")
	ArrayIps, err := http.GetCloudflareIps()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("All cloudflare IP Ranges/Ipv4 Segments:\r\n")
	for i := 0; i < len(ArrayIps); i++ {
		fmt.Printf("%s\r\n", ArrayIps[i])
	}
	ipSgts, err := cidr.ParseCidr(ArrayIps, []int{100, 22}, true, true)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(ipSgts); i++ {
		log.Printf("[%s]开始处理[循环次数为%d]...\r\n", ipSgts[i].NodeId, number)
		log.Printf("本次ip数量为: [%d]\r\n", len(ipSgts[i].Node))
		pingClient.Pings(ipSgts[i].Node, number, singleNumber, limiterNumber, coloOpenFlag, "AllCfIpv4PingTo["+strconv.Itoa(i+1)+"].csv", iataFilePath)
		log.Printf("第[%d]批次已处理完毕![%s]\r\n\n", i+1, ipSgts[i].NodeId)
	}
}
