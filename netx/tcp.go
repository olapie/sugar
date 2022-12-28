package netx

import (
	"fmt"
	"net"
)

func FindTCPPort(ip string, minPort, maxPort int) int {
	for port := minPort; port <= maxPort; port++ {
		addr := fmt.Sprintf("%s:%d", ip, port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			l.Close()
			return port
		}
	}
	return 0
}
