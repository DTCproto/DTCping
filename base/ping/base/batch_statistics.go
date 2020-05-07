package base

// Statistics is all addr data Statistic
func BatchStatistics(bp *BatchPacketCore) []*Statistics {
	// 预估长度
	stMap := make([]*Statistics, 0, len(bp.mapIpPacketCore))
	// log.Printf("START %d %d %d \n", len(stMap), len(bp.mapIpPacketCore), cap(stMap))
	for ip := range bp.mapIpPacketCore {
		stMap = append(stMap, bp.mapIpPacketCore[ip].Statistics())
	}
	// log.Printf("END %d %d %d \n", len(stMap), len(bp.mapIpPacketCore), cap(stMap))
	return stMap
}

// Finish will call OnFinish
func Finish(bp *BatchPacketCore) {
	handler := bp.OnFinish
	if bp.OnFinish != nil {
		handler(BatchStatistics(bp))
	}
}
