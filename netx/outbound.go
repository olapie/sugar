package netx

import (
	"fmt"
	"net"
)

// GetOutboundIP returns preferred outbound ip of this machine
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	if err = conn.Close(); err != nil {
		return nil, err
	}
	return localAddr.IP, nil
}

// GetOutboundIPString returns preferred outbound ip of this machine
func GetOutboundIPString() string {
	ip, err := GetOutboundIP()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return ip.String()
}
