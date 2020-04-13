package base

// Statistics is all addr data Statistic
func BatchStatistics(bp *BatchPacketCore) map[string]*Statistics {
	stMap := map[string]*Statistics{}
	for ip, packetCore := range bp.mapIpPacketCore {
		addr := bp.mapIpAddr[ip]
		stMap[addr] = packetCore.Statistics()
	}
	return stMap
}

// Finish will call OnFinish
func Finish(bp *BatchPacketCore) {
	handler := bp.OnFinish
	if bp.OnFinish != nil {
		handler(BatchStatistics(bp))
	}
}
