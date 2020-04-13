package base

import (
	"log"
	"net"
	"sync"
	"time"
)

func recvIpv4(bp *BatchPacketCore, wg *sync.WaitGroup) {
	defer wg.Done()
	var ttl int

	for {
		select {
		case <-bp.done:
			return
		default:
			bytes := make([]byte, 512)
			_ = bp.conn4.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
			n, cm, addr, err := bp.conn4.IPv4PacketConn().ReadFrom(bytes)
			if cm != nil {
				ttl = cm.TTL
			}

			if err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						continue
					} else {
						if bp.Debug {
							log.Printf("read err %s ", err)
						}
						return
					}
				}
			}

			recvPkg := &packet{bytes: bytes, nBytes: n, ttl: ttl, proto: protoIpv4, addr: addr}
			if bp.Debug {
				log.Printf("recv addr %v \n", recvPkg.addr.String())
			}
			err = processPacket(bp, recvPkg)
			if err != nil && bp.Debug {
				log.Printf("processPacket err %v, recvpkg %v \n", err, recvPkg)
			}
		}
	}
}

func recvIpv6(bp *BatchPacketCore, wg *sync.WaitGroup) {
	defer wg.Done()
	var ttl int
	for {
		select {
		case <-bp.done:
			return
		default:
			bytes := make([]byte, 512)
			_ = bp.conn6.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
			n, cm, addr, err := bp.conn6.IPv6PacketConn().ReadFrom(bytes)
			if cm != nil {
				ttl = cm.HopLimit
			}
			if err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						continue
					}
				}
			}

			recvPkg := &packet{bytes: bytes, nBytes: n, ttl: ttl, proto: protoIpv6, addr: addr}
			if bp.Debug {
				log.Printf("recv addr %v \n", recvPkg.addr.String())
			}
			err = processPacket(bp, recvPkg)
			if err != nil && bp.Debug {
				log.Printf("processPacket err %v, recvpkg %v \n", err, recvPkg)
			}
		}
	}
}
