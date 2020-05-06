package sort

import (
	ping "DTCping/base/ping/base"
	"sort"
)

// 方法2:
// 默认Slice [相同元素保持原排序SliceStable]
func ProcessStatisticsSlice(stMaps []map[string]*ping.Statistics) []Base {
	var sSortSlice []Base
	for i := range stMaps {
		for ip, st := range stMaps[i] {
			sSortSlice = appendBase(sSortSlice, Base{
				Ip:   ip,
				Data: st,
			})
		}
	}
	sort.Slice(sSortSlice, func(i, j int) bool {
		if sSortSlice[i].Data.PacketLoss != sSortSlice[j].Data.PacketLoss {
			return sSortSlice[i].Data.PacketLoss < sSortSlice[j].Data.PacketLoss
		} else {
			// 平均延迟
			// return sSortSlice[i].Data.AvgRtt.Nanoseconds() <= sSortSlice[j].Data.AvgRtt.Nanoseconds()
			// 最大延迟
			return sSortSlice[i].Data.MaxRtt.Nanoseconds() <= sSortSlice[j].Data.MaxRtt.Nanoseconds()
		}
	})
	return sSortSlice
}

// 默认append函数在for循环里容易out of memory
func appendBase(slice []Base, data ...Base) []Base {
	m := len(slice)
	n := m + len(slice)
	if n > cap(slice) {
		slice = append(slice, data...)
		newSlice := make([]Base, len(slice))
		copy(newSlice, slice)
		return newSlice
	} else {
		slice = append(slice, data...)
		return slice
	}
}
