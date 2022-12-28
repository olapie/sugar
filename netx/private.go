package netx

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"code.olapie.com/sugar/conv"
)

func GetPrivateIP4String() string {
	_, ip, err := GetPrivateIP4()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return ip.String()
}

func GetPrivateIP4() (*net.Interface, net.IP, error) {
	ifi, addr, err := GetPrivateIP4Addr()
	if err != nil {
		return ifi, nil, err
	}

	if ipNet, ok := addr.(*net.IPNet); ok {
		if ip := ipNet.IP.To4(); ip != nil {
			if ip.IsPrivate() {
				return ifi, ip, nil
			}
		}
	}
	return nil, nil, errors.New("not found")
}

func GetPrivateIP4Addr() (*net.Interface, net.Addr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, nil, err
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip := ipNet.IP.To4(); ip != nil {
					if ip.IsPrivate() {
						return &i, addr, nil
					}
				}
			}
		}
	}
	return nil, nil, errors.New("not found")
}

/*
IsPrivateIP4 tells if ip is address of local area network
24-bit block	10.0.0.0 – 10.255.255.255	16777216	10.0.0.0/8 (255.0.0.0)	24 bits	8 bits	single class A network
20-bit block	172.16.0.0 – 172.31.255.255	1048576	172.16.0.0/12 (255.240.0.0)	20 bits	12 bits	16 contiguous class B networks
16-bit block	192.168.0.0 – 192.168.255.255	65536	192.168.0.0/16 (255.255.0.0)	16 bits	16 bits	256 contiguous class C networks
*/
func IsPrivateIP4[T string | net.IP](ipOrString T) bool {
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
		return IsPrivateIP4(net.IPv4(a, b, c, d))
	}

	ip := any(ipOrString).(net.IP)
	return ip.IsPrivate()
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
