// Package local checks if an addresse is on this local machine
package local

import (
	"log"
	"net"
	"strings"
)

func localAddrs() []string {
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		log.Println(e)
		return []string{}
	}

	ret := make([]string, 0, len(addrs))

	for _, addr := range addrs {
		ret = append(ret, addr.String())
	}

	return ret
}

func Check(addr string) bool {
	a, e := net.ResolveTCPAddr("tcp", addr)
	if e != nil {
		return false
	}

	ip := a.IP.String()

	addrs := localAddrs()
	for _, addr := range addrs {
		if strings.HasPrefix(addr, ip) {
			return true
		}
	}

	return false
}
