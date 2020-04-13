package base

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"log"
	"net"
	"time"
)

func processPacket(bp *BatchPacketCore, recv *packet) error {
	receivedAt := time.Now()
	var proto int
	if recv.proto == protoIpv4 {
		proto = protocolICMP
	} else {
		proto = protocolIPv6ICMP
	}

	var m *icmp.Message
	var err error

	if m, err = icmp.ParseMessage(proto, recv.bytes); err != nil {
		return fmt.Errorf("error parsing icmp message: %s", err.Error())
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// Not an echo reply, ignore it
		if bp.Debug {
			log.Printf("pkg drop %v \n", m)
		}
		return nil
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:
		// If we are privileged, we can match icmp.ID
		if pkt.ID != bp.id {
			if bp.Debug {
				log.Printf("drop pkg %+v id %v addr %s \n", pkt, bp.id, recv.addr)
			}
			return nil
		}

		if len(pkt.Data) < timeSliceLength+trackerLength {
			return fmt.Errorf("insufficient data received; got: %d %v",
				len(pkt.Data), pkt.Data)
		}

		timestamp := bytesToTime(pkt.Data[:timeSliceLength])

		var ip string
		if bp.network == NetworkUDP {
			if ip, _, err = net.SplitHostPort(recv.addr.String()); err != nil {
				return fmt.Errorf("err ip : %v, err %v", recv.addr, err)
			}
		} else {
			ip = recv.addr.String()
		}

		// // 重复包舍弃
		if packetCore, ok := bp.mapIpPacketCore[ip]; ok {
			if packetCore.lastSeq[pkt.Seq+1] == pkt.Seq+1 {
				packetCore.PacketsRecvDup++
			} else {
				packetCore.PacketsRecv++
				packetCore.rtts = append(packetCore.rtts, receivedAt.Sub(timestamp))
			}
			packetCore.lastSeq[pkt.Seq+1] = pkt.Seq + 1
		}
	default:
		// Very bad, not sure how this can happen
		return fmt.Errorf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt)
	}

	return nil

}
