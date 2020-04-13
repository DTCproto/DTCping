package base

import (
	"bytes"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"syscall"
	"time"
)

func (p *PacketCore) SendICMP(seqID int) {
	var typ icmp.Type
	if p.ipv4 {
		typ = ipv4.ICMPTypeEcho
	} else {
		typ = ipv6.ICMPTypeEchoRequest
	}

	var dst net.Addr = p.ipAddr
	if p.network == NetworkUDP {
		dst = &net.UDPAddr{IP: p.ipAddr.IP, Zone: p.ipAddr.Zone}
	}

	t := append(timeToBytes(time.Now()), intToBytes(p.Tracker)...)
	if remainSize := p.Size - timeSliceLength - trackerLength; remainSize > 0 {
		t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}

	body := &icmp.Echo{
		ID:   p.id,
		Seq:  seqID,
		Data: t,
	}

	msg := &icmp.Message{
		Type: typ,
		Code: 0,
		Body: body,
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return
	}

	for {
		if p.ipv4 {
			if _, err := p.conn4.WriteTo(msgBytes, dst); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Err == syscall.ENOBUFS {
						continue
					}
				}
			}
		} else {
			if _, err := p.conn6.WriteTo(msgBytes, dst); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Err == syscall.ENOBUFS {
						continue
					}
				}
			}
		}
		break
	}

	return
}
