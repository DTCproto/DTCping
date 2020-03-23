package cidr

import (
	"errors"
	"fmt"
	"net"
)

func ParseCidrSingle(cidrStr string) ([]string, error) {
	ipAddr, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ipAddr.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if ips == nil {
		return ips, errors.New("CIDR Parse ips 数量异常！")
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for ipsNum := len(ip) - 1; ipsNum >= 0; ipsNum-- {
		ip[ipsNum]++
		if ip[ipsNum] > 0 {
			break
		}
	}
}

type IpSgt struct {
	NodeId string
	Node   []string
}

func ParseCidr(IpSegments []string, rdValues []int, mergeFlag, extractFlag bool) ([]IpSgt, error) {
	var results []IpSgt
	// 优先度设置
	if extractFlag {
		mergeFlag = true
	}
	// 是否合并
	if mergeFlag {
		var ipsAllTemp []string
		for i := 0; i < len(IpSegments); i++ {
			ipsTemp, err := ParseCidrSingle(IpSegments[i])
			if err != nil || len(ipsTemp) <= 0 {
				break
			}
			ipsAllTemp = append(ipsAllTemp, ipsTemp...)
		}
		if extractFlag {
			ipsAllTemp = ParseArrayToRandomValue(ipsAllTemp, rdValues)
			ipsAllTemp = RemoveRepeatedElement(ipsAllTemp)
		}
		var node IpSgt
		node.NodeId = fmt.Sprint(IpSegments[0]) + " ..."
		node.Node = ipsAllTemp
		results = append(results, node)
	} else {
		for i := 0; i < len(IpSegments); i++ {
			ipsTemp, err := ParseCidrSingle(IpSegments[i])
			if err != nil || len(ipsTemp) <= 0 {
				break
			}
			var node IpSgt
			node.NodeId = IpSegments[i]
			node.Node = ipsTemp
			results = append(results, node)
		}
	}
	if len(results) <= 0 {
		return nil, errors.New("无有效Ip Segment,请确认！\r\n")
	}
	return results, nil
}
