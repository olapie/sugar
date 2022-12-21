package netx

import (
	"fmt"
	"net"
	"time"
)

// BroadcastUDP sends packet to all devices in the same LAN
// UDP packet is carried by one IP packet
// IP packet is limited by MTU(Maximum Transmission Unit)
// MTU is around 1400-1500
// so one UDP packet should be less than 1400 with UDP header, IP header, ...
// Least MTU is 576, so UDP packet around 500 is pretty safe
func BroadcastUDP(port int, packet []byte) error {
	addr := fmt.Sprintf("255.255.255.255:%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	_, err = udpConn.Write(packet)
	return err
}

func RepeatBroadcastUDP(port int, packet []byte, interval time.Duration) {
	for {
		err := BroadcastUDP(port, packet)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(interval)
	}
}

func ReceiveUDP(port int, timeout time.Duration) ([]byte, net.Addr, error) {
	conn, err := net.ListenPacket("udp", ":"+fmt.Sprint(port))
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, nil, err
	}

	buf := make([]byte, 1500)
	nRead, addr, err := conn.ReadFrom(buf)
	if err != nil {
		return nil, nil, err
	}
	return buf[:nRead], addr, nil
}
