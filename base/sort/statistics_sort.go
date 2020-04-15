package sort

import (
	ping "DTCping/base/ping/base"
	"sort"
)

type StatisticsSort []Base

func (s StatisticsSort) Len() int {
	return len(s)
}

func (s StatisticsSort) Less(i, j int) bool {
	if s[i].Data.PacketLoss != s[j].Data.PacketLoss {
		return s[i].Data.PacketLoss < s[j].Data.PacketLoss
	} else {
		// 平均延迟
		// return s[i].Data.AvgRtt.Nanoseconds() <= s[j].Data.AvgRtt.Nanoseconds()
		// 最大延迟
		return s[i].Data.MaxRtt.Nanoseconds() <= s[j].Data.MaxRtt.Nanoseconds()
	}
}

func (s StatisticsSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 方法1:
// 默认Sort [相同元素保持原排序Stable]
func ProcessStatistics(stMaps []map[string]*ping.Statistics) []Base {
	var statisticsSort []Base
	for i := range stMaps {
		for ip, st := range stMaps[i] {
			statisticsSort = append(statisticsSort, Base{
				Ip:   ip,
				Data: st,
			})
		}
	}
	sort.Sort(StatisticsSort(statisticsSort))
	return statisticsSort
}
