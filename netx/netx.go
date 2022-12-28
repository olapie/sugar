package netx

import (
	"fmt"
	"net"
)

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

func GetIP4Info() (ifaceToIPList map[string][]string, err error) {
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

func GetIP6Info() (ifaceToIPList map[string][]string, err error) {
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

func GetIP4(ifaceName string) string {
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

func GetIP6(ifaceName string) string {
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
				if ip := ipNet.IP.To16(); ip != nil {
					return ip.String()
				}
			}
		}
	}
	return ""
}
