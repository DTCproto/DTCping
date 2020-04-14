package sort

import (
	ping "DTCping/base/ping/base"
	"sort"
)

type Base struct {
	Ip   string
	Data *ping.Statistics
}

type StatisticsSort []Base

func (s StatisticsSort) Len() int {
	return len(s)
}

func (s StatisticsSort) Less(i, j int) bool {
	if s[i].Data.PacketLoss != s[j].Data.PacketLoss {
		return s[i].Data.PacketLoss < s[j].Data.PacketLoss
	} else {
		return s[i].Data.AvgRtt.Nanoseconds() <= s[j].Data.AvgRtt.Nanoseconds()
	}
}

func (s StatisticsSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

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
