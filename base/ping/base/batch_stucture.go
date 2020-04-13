package base

import (
	"golang.org/x/net/icmp"
	"time"
)

// BatchPacketCore is PacketCore manager
type BatchPacketCore struct {
	done chan bool

	// mapIpAddr is ip addr map
	mapIpAddr map[string]string

	// mapIpPacketCore is ip PacketCore map
	mapIpPacketCore map[string]*PacketCore

	// interval是每个数据包发送之间的等待时间[默认值为1s]
	interval time.Duration

	// 超时指定基准退出之前的超时，无论已接收多少个数据包
	timeout time.Duration

	//count是每个用户的基数
	Count int

	//sendCount是已发送的数字
	sendCount int

	//网络是网络模式，可以是ip或udp，并且ip具有特权
	network string

	//id是[唯一识别码]
	id int

	//conn4 is ipv4 icmp PacketConn
	conn4 *icmp.PacketConn

	//conn6 is ipv6 icmp PacketConn
	conn6 *icmp.PacketConn

	//所有的addr
	addrs []string

	// 当PacketCore退出时可以调用OnFinish
	OnFinish func(map[string]*Statistics)

	Debug bool

	// SeqID是ICMP序列号
	seqID int
}
