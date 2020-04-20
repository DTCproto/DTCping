package file

import (
	"encoding/json"
	"io/ioutil"
)

type PathConfig struct {
	// 遍历次数
	Number int `json:"number"`
	// 保存文件名
	CsvFileName string `json:"csv_file_name"`
	// IP段列表
	IpSegments []string `json:"ip_segments"`
	// 是否合并结果到一个CSV
	MergeFlag bool `json:"merge_flag"`
	// 抽取D段地址值列表(1-254)
	RandomValues []int `json:"random_values"`
	// 是否抽取D段地址进行查询[优先度高于MergeFlag]
	ExtractFlag bool `json:"extract_flag"`
	// 是否查询Colo信息
	ColoFlag bool `json:"colo_flag"`
	// 优先读取本地iata对照文件
	IataSrcPath string `json:"iata_src_path"`
	// 单批次ping的ip数量
	EachSingleNum int `json:"each_single_num"`
	// COLO并发数
	ColoLimiterNum int `json:"colo_limiter_num"`
}

func (c *PathConfig) ParseConfigFile(filePath string) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileData, c)
	return err
}
