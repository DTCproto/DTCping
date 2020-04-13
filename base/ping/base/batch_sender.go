package base

import (
	"log"
	"sync"
	"time"
)

func sendICMP(bp *BatchPacketCore, wg *sync.WaitGroup) {
	defer wg.Done()
	timeout := time.NewTicker(bp.timeout)
	interval := time.NewTicker(bp.interval)

	for {
		select {
		case <-bp.done:
			return

		case <-timeout.C:
			close(bp.done)
			return

		case <-interval.C:
			batchSendICMP(bp)
			bp.sendCount++
			if bp.sendCount >= bp.Count {
				time.Sleep(bp.interval)
				close(bp.done)
				if bp.Debug {
					log.Printf("send end sendcout %d, count %d \n", bp.sendCount, bp.Count)
				}
				return
			}
		}
	}
}

// batchSendICMP let all addr send pkg once
func batchSendICMP(bp *BatchPacketCore) {
	for _, packetCore := range bp.mapIpPacketCore {
		packetCore.SendICMP(bp.seqID)
		packetCore.PacketsSent++
	}
	bp.seqID = (bp.seqID + 1) & 0xffff
}
