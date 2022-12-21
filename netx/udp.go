package netx

import (
	"fmt"
	"net"
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

//func ReceiveUDP(port int) {
//	conn, err := net.ListenPacket("udp", ":1053")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
//	buf := make([]byte, 14*1024*1024)
//	for {
//		nRead, addr, err := conn.ReadFrom(buf)
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			fmt.Println(addr, buf[:nRead])
//		}
//	}
//}
