package sort

import (
	ping "DTCping/base/ping/base"
	"sort"
)

// 方法2:
// 默认Slice [相同元素保持原排序SliceStable]
func ProcessStatisticsSlice(stMaps []*ping.Statistics) []*ping.Statistics {
	sort.Slice(stMaps, func(i, j int) bool {
		if stMaps[i].PacketLoss != stMaps[j].PacketLoss {
			return stMaps[i].PacketLoss < stMaps[j].PacketLoss
		} else {
			// 平均延迟
			// return sSortSlice[i].Data.AvgRtt.Nanoseconds() <= sSortSlice[j].Data.AvgRtt.Nanoseconds()
			// 最大延迟
			return stMaps[i].MaxRtt.Nanoseconds() <= stMaps[j].MaxRtt.Nanoseconds()
		}
	})
	return stMaps
}

// 默认append函数在for循环里容易out of memory [等待后续优化]
