package base

import (
	"golang.org/x/net/icmp"
	"net"
	"time"
)

// 统计信息表示当前正在运行或已完成的ping操作的统计信息。
type Statistics struct {

	// PacketsRecv是接收到的数据包数。
	PacketsRecv int

	// PacketsSent是发送的数据包数。
	PacketsSent int

	// 收到重复数据包的数量
	PacketsRecvDup int

	// 数据包丢失是数据包丢失的百分比。
	PacketLoss float64

	// IPAddr是要被ping通的主机的地址。
	IPAddr *net.IPAddr

	// Addr是要被ping通的主机的字符串地址。
	Addr string

	// Rtts是通过此ping发送的所有往返时间
	Rtts []time.Duration

	// MinRtt是通过此ping发送的最短往返时间
	MinRtt time.Duration

	// MaxRtt是通过此ping发送的最大往返时间
	MaxRtt time.Duration

	// AvgRtt是通过此ping发送的平均往返时间
	AvgRtt time.Duration

	// StdDevRtt是通过此ping发送的往返时间的标准偏差
	StdDevRtt time.Duration
}

// Packet代表已接收和已处理的ICMP回显数据包
type Packet struct {

	// Rtt是达到基准所需的往返时间
	Rtt time.Duration

	// IPAddr是要被ping通的主机的地址
	IPAddr *net.IPAddr

	// Addr是要被ping通的主机的字符串地址
	Addr string

	// NumBytes是消息中的字节数
	NumBytes int

	// Seq是ICMP序列号
	Seq int

	// TTL是数据包上的生存时间
	Ttl int
}

// PacketCore代表ICMP数据包sender/receiver
type PacketCore struct {

	// 间隔是每个数据包发送之间的等待时间[默认值为1s]
	Interval time.Duration

	// 超时指定基准退出之前的超时，无论已接收多少个数据包
	Timeout time.Duration

	// Count告诉PacketCore在 发送/接收 Count回显数据包后停止
	// 如果未指定此选项，则PacketCore将一直运行直到被中断
	Count int

	// 发送的包数
	PacketsSent int

	// 收到的包数
	PacketsRecv int

	// 收到重复数据包的数量
	PacketsRecvDup int

	// 所有的Rtts的时间集合
	rtts []time.Duration

	// 当PacketCore接收并处理数据包时调用OnRecv
	OnRecv func(*Packet)

	// 当PacketCore退出时调用OnFinish
	OnFinish func(*Statistics)

	// 发送的数据包大小
	Size int

	// Tracker 跟踪器，用于在非特权时唯一标识数据包
	Tracker int64

	// Source是源IP地址
	Source string

	ipAddr *net.IPAddr
	addr   string

	ipv4    bool
	size    int
	id      int
	network string
	lastSeq map[int]int

	// conn4 is ipv4 icmp PacketConn
	conn4 *icmp.PacketConn

	// conn6 is ipv6 icmp PacketConn
	conn6 *icmp.PacketConn
}

type packet struct {
	bytes  []byte
	nBytes int
	ttl    int
	proto  string
	addr   net.Addr
}
