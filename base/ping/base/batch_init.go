package base

import "time"

func NewBatchPing(addrs []string, privileged bool) (*BatchPacketCore, error) {
	network := NetworkUDP
	if privileged {
		network = NetworkIP
	}
	batchPacketCore := &BatchPacketCore{
		interval:        time.Second,
		timeout:         time.Second * 100000,
		Count:           10,
		network:         network,
		id:              GetPId(),
		done:            make(chan bool),
		addrs:           addrs,
		mapIpPacketCore: make(map[string]*PacketCore),
		mapIpAddr:       make(map[string]string),
	}
	return batchPacketCore, nil
}
