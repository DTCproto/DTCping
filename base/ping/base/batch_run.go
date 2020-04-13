package base

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"log"
	"sync"
)

func (bp *BatchPacketCore) Run() (err error) {
	if bp.conn4, err = icmp.ListenPacket(ipv4Proto[bp.network], SockAddrV4); err != nil {
		return err
	}
	if bp.conn6, err = icmp.ListenPacket(ipv6Proto[bp.network], SockAddrV6); err != nil {
		return err
	}
	_ = bp.conn4.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
	_ = bp.conn6.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)

	for _, addr := range bp.addrs {
		packetCore, err := NewPing(addr, bp.id, bp.network)
		if err != nil {
			return err
		}
		bp.mapIpPacketCore[packetCore.ipAddr.String()] = packetCore
		bp.mapIpAddr[packetCore.ipAddr.String()] = addr
		packetCore.conn4 = bp.conn4
		packetCore.conn6 = bp.conn6
	}

	if bp.Debug {
		log.Printf("[debug] pid %d \n", bp.id)
	}

	defer bp.conn4.Close()
	defer bp.conn6.Close()

	var wg sync.WaitGroup
	wg.Add(3)
	go recvIpv4(bp, &wg)
	go recvIpv6(bp, &wg)
	go sendICMP(bp, &wg)
	wg.Wait()
	return nil
}
