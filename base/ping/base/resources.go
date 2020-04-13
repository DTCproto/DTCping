package base

const (
	SockAddrV4       = "0.0.0.0"
	SockAddrV6       = "::"
	NetworkIP        = "ip"
	NetworkUDP       = "udp"
	NetworkTCP       = "tcp"
	timeSliceLength  = 8
	trackerLength    = 8
	protocolICMP     = 1
	protocolIPv6ICMP = 58
	protoIpv4        = "ipv4"
	protoIpv6        = "ipv6"
)

var (
	ipv4Proto = map[string]string{"ip": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"ip": "ip6:ipv6-icmp", "udp": "udp6"}
)
