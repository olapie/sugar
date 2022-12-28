package netx

import (
	"fmt"
	"net"
)

// GetOutboundIP returns preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("net.Dial:", err)
		return nil
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	if err = conn.Close(); err != nil {
		fmt.Println("conn.Close:", err)
		return nil
	}
	return localAddr.IP
}

// GetOutboundIPString returns preferred outbound ip of this machine
func GetOutboundIPString() string {
	ip := GetOutboundIP()
	if ip != nil {
		return ip.String()
	}
	return ""
}
