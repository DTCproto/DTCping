package http

import "net"

func LookupIps(host string) ([]string, error) {
	var result []string
	ns, err := net.LookupIP(host)
	if err != nil {
		return result, err
	}
	size := len(ns)
	for i := 0; i < size; i++ {
		result = append(result, ns[i].String())
	}
	return result, nil
}
