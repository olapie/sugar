package netx

import (
	"code.olapie.com/sugar/conv"
	"fmt"
	"net"
	"strings"
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

func GetMacAddresses() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	addrs := make([]string, len(ifaces))
	for i, ifa := range ifaces {
		addrs[i] = ifa.HardwareAddr.String()
	}

	return addrs, nil
}

func GetIPv4Info() (ifaceToIPList map[string][]string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("get net interfaces: %w", err)
	}
	ifaceToIPList = make(map[string][]string)
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("get addrs %s: %w", i.Name, err)
		}

		for _, addr := range addrs {
			//var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				if ip4 := v.IP.To4(); ip4 != nil {
					ifaceToIPList[i.Name] = append(ifaceToIPList[i.Name], ip4.String())
				}
			case *net.IPAddr:
				//	if ip4 := v.IP.To4(); ip4 != nil {
				//		res[i.Name] = append(res[i.Name], ip4.String())
				//	}
			}
		}
	}
	return ifaceToIPList, nil
}

func GetIPv6Info() (ifaceToIPList map[string][]string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("get net interfaces: %w", err)
	}
	ifaceToIPList = make(map[string][]string)
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("get addrs %s: %w", i.Name, err)
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if ip6 := v.IP.To16(); ip6 != nil {
					ifaceToIPList[i.Name] = append(ifaceToIPList[i.Name], ip6.String())
				}
			case *net.IPAddr:
				//	if ip4 := v.IP.To4(); ip4 != nil {
				//		res[i.Name] = append(res[i.Name], ip4.String())
				//	}
			}
		}
	}
	return ifaceToIPList, nil
}

func GetIFaceNames() []string {
	var a []string
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return a
	}
	for _, ifa := range ifaces {
		a = append(a, ifa.Name)
	}
	return a
}

func GetIPv4(ifaceName string) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, i := range ifaces {
		if i.Name != ifaceName {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return ""
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip := ipNet.IP.To4(); ip != nil {
					return ip.String()
				}
			}
		}
	}
	return ""
}

func GetIPv6(ifaceName string) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return ""
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip := ipNet.IP.To16(); ip != nil {
					return ip.String()
				}
			}
		}
	}
	return ""
}

func GetLocalIP4String() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return ""
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip := ipNet.IP.To4(); ip != nil {
					fmt.Println(ip.String())
					if IsLocalIP4(ip) {
						return ip.String()
					}
				}
			}
		}
	}
	return ""
}

/*
IsLocalIP4 tells if ip is address of local area network
24-bit block	10.0.0.0 – 10.255.255.255	16777216	10.0.0.0/8 (255.0.0.0)	24 bits	8 bits	single class A network
20-bit block	172.16.0.0 – 172.31.255.255	1048576	172.16.0.0/12 (255.240.0.0)	20 bits	12 bits	16 contiguous class B networks
16-bit block	192.168.0.0 – 192.168.255.255	65536	192.168.0.0/16 (255.255.0.0)	16 bits	16 bits	256 contiguous class C networks
*/
func IsLocalIP4[T string | net.IP](ipOrString T) bool {
	if ipStr, ok := any(ipOrString).(string); ok {
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}
		a, ok := parseIP4Part(arr[0])
		if !ok {
			return false
		}
		b, ok := parseIP4Part(arr[0])
		if !ok {
			return false
		}
		c, ok := parseIP4Part(arr[0])
		if !ok {
			return false
		}
		d, ok := parseIP4Part(arr[0])
		if !ok {
			return false
		}
		return IsLocalIP4(net.IPv4(a, b, c, d))
	}

	ip := any(ipOrString).(net.IP)
	a0, a1 := net.IPv4(10, 0, 0, 0), net.IPv4(10, 255, 255, 255)
	if compareIP4(ip, a0) >= 0 && compareIP4(ip, a1) <= 0 {
		return true
	}
	b0, b1 := net.IPv4(172, 16, 0, 0), net.IPv4(172, 31, 255, 255)
	if compareIP4(ip, b0) >= 0 && compareIP4(ip, b1) <= 0 {
		return true
	}
	c0, c1 := net.IPv4(192, 168, 0, 0), net.IPv4(192, 168, 255, 255)
	if compareIP4(ip, c0) >= 0 && compareIP4(ip, c1) <= 0 {
		return true
	}
	return false
}

func parseIP4Part(s string) (byte, bool) {
	a, err := conv.ToInt(s)
	if err != nil {
		return 0, false
	}

	if a < 0 || a > 255 {
		return 0, false
	}

	return byte(a), true
}

func compareIP4(a, b net.IP) int {
	a = a.To4()
	b = b.To4()
	for i := 0; i < net.IPv4len; i++ {
		if a[i] == b[i] {
			continue
		}
		return int(a[i]) - int(b[i])
	}
	return 0
}
