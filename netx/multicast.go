package netx

import (
	"net"
	"sort"
	"strings"
)

func GetMulticastIP4String(ifi *net.Interface) string {
	addrs := GetMulticastIP4Addrs(ifi)
	if len(addrs) == 0 {
		return ""
	}
	sort.Slice(addrs, func(i, j int) bool {
		si, sj := addrs[i].String(), addrs[j].String()
		if len(si) == len(sj) {
			return si < sj
		}
		return len(si) < len(sj)
	})
	return addrs[0].String()
}

func GetMulticastIP4Addrs(ifi *net.Interface) []net.Addr {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	for _, i := range interfaces {
		if i.Name != ifi.Name {
			continue
		}
		addrs, err := i.MulticastAddrs()
		if err != nil {
			return nil
		}

		for j := len(addrs) - 1; j >= 0; j-- {
			if strings.HasPrefix(addrs[j].String(), "224") {
				continue
			}
			addrs = append(addrs[0:j], addrs[j+1:]...)
		}
		return addrs
	}
	return nil
}
