package netx

import (
	"code.olapie.com/sugar/conv"
	"encoding/binary"
	"fmt"
	"net"
	"sort"
	"strings"
)

func ListMacAddresses() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces:", err)
		return nil
	}

	addrs := make([]string, len(ifaces))
	for i, ifa := range ifaces {
		addrs[i] = ifa.HardwareAddr.String()
	}

	return addrs
}

func ListInterfaceNames() []string {
	var a []string
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces:", err)
		return nil
	}
	for _, ifa := range ifaces {
		a = append(a, ifa.Name)
	}
	return a
}

func ListIPNets(ifi *net.Interface) []*net.IPNet {
	addrs, err := ifi.Addrs()
	if err != nil {
		fmt.Println("net.Interfaces:", err)
		return nil
	}

	var a []*net.IPNet
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ip := ipNet.IP.To4(); ip != nil {
				a = append(a, ipNet)
			}
		}
	}
	return a
}

func NewIPv4FromString(s string) net.IP {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return nil
	}

	b := make([]byte, 4)
	for i := range parts {
		k, err := conv.ToInt(parts[i])
		if err != nil {
			return nil
		}

		if k < 0 || k > 255 {
			return nil
		}

		b[i] = byte(k)
	}
	return net.IPv4(b[0], b[1], b[2], b[3])
}

func GetPrivateIPv4(ifi *net.Interface) net.IP {
	ipNet := GetPrivateIPv4Net(ifi)
	if ipNet != nil {
		return ipNet.IP.To4()
	}
	return nil
}

func GetPrivateIPv4Net(ifi *net.Interface) *net.IPNet {
	addrs, err := ifi.Addrs()
	if err != nil {
		fmt.Println("net.Interface.Addrs:", err)
		return nil
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ip := ipNet.IP.To4(); ip != nil {
				if ip.IsPrivate() {
					return ipNet
				}
			}
		}
	}
	return nil
}

/*
GetPrivateIPv4Interface returns an interface with private ip, a LAN interface
24-bit block	10.0.0.0 – 10.255.255.255	16777216	10.0.0.0/8 (255.0.0.0)	24 bits	8 bits	single class A network
20-bit block	172.16.0.0 – 172.31.255.255	1048576	172.16.0.0/12 (255.240.0.0)	20 bits	12 bits	16 contiguous class B networks
16-bit block	192.168.0.0 – 192.168.255.255	65536	192.168.0.0/16 (255.255.0.0)	16 bits	16 bits	256 contiguous class C networks
*/
func GetPrivateIPv4Interface() *net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip := ipNet.IP.To4(); ip != nil {
					if ip.IsPrivate() {
						return &i
					}
				}
			}
		}
	}
	return nil
}

// GetBroadcastIPv4 returns broadcast ip
// Class A, B, and C networks have natural masks, or default subnet masks:
// Class A: 255.0.0.0
// Class B: 255.255.0.0
// Class C: 255.255.255.0
func GetBroadcastIPv4(ifi *net.Interface) net.IP {
	if ifi == nil {
		return nil
	}
	ipNet := GetPrivateIPv4Net(ifi)
	ip4 := ipNet.IP.To4()
	if ip4 == nil {
		return nil
	}
	ones, _ := ipNet.Mask.Size()
	return CalculateBroadcastIPv4(ip4, ones)
}

func CalculateBroadcastIPv4(ip4 []byte, maskOnes int) []byte {
	if len(ip4) != 4 {
		panic("ip4 length is not 4")
	}
	zeros := 32 - maskOnes
	if zeros == 0 {
		return []byte{255, 255, 255, 255}
	}
	n := binary.BigEndian.Uint32(ip4)
	mask := uint32(1) << (zeros - 1)
	n &= ^mask
	n |= (1 << zeros) - 1
	return binary.BigEndian.AppendUint32(nil, n)
}

func GetMulticastIPv4String(ifi *net.Interface) string {
	addrs := GetMulticastIPv4Addrs(ifi)
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

func GetMulticastIPv4Addrs(ifi *net.Interface) []net.Addr {
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

// BroadcastUDP sends packet to all devices in the same LAN
// UDP packet is carried by one IP packet
// IP packet is limited by MTU(Maximum Transmission Unit)
// MTU is around 1400-1500
// so one UDP packet should be less than 1400 with UDP header, IP header, ...
// Least MTU is 576, so UDP packet around 500 is pretty safe
//func BroadcastUDP(addr *net.UDPAddr, packet []byte) error {
//	udpConn, err := net.DialUDP("udp", nil, addr)
//	if err != nil {
//		return err
//	}
//
//	_, err = udpConn.Write(packet)
//	return err
//}

//func ReceiveUDP(port int, timeout time.Duration) ([]byte, net.Addr, error) {
//	conn, err := net.ListenPacket("udp", ":"+fmt.Sprint(port))
//	if err != nil {
//		return nil, nil, err
//	}
//	defer conn.Close()
//	err = conn.SetReadDeadline(time.Now().Add(timeout))
//	if err != nil {
//		return nil, nil, err
//	}
//
//	buf := make([]byte, 1500)
//	nRead, addr, err := conn.ReadFrom(buf)
//	if err != nil {
//		return nil, nil, err
//	}
//	return buf[:nRead], addr, nil
//}
