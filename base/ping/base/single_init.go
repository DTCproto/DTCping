package base

import (
	"math"
	"math/rand"
	"net"
	"time"
)

func NewPing(addr string, pid int, network string) (*PacketCore, error) {
	ipAddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	var ipv4 bool
	if isIPv4(ipAddr.IP) {
		ipv4 = true
	} else if isIPv6(ipAddr.IP) {
		ipv4 = false
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &PacketCore{
		ipAddr:   ipAddr,
		addr:     addr,
		Interval: time.Second,
		Timeout:  time.Second * 100000,
		Count:    -1,
		id:       pid,
		network:  network,
		ipv4:     ipv4,
		Size:     timeSliceLength,
		Tracker:  r.Int63n(math.MaxInt64),
		lastSeq:  make(map[int]int),
	}, nil
}
