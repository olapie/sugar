package mob

import (
	"code.olapie.com/sugar/v2/xnetwork"
	"net"
	"time"
)

const (
	NoNetwork = 0
	Cellular  = 1
	WIFI      = 2
)

const (
	Idle       = 0
	Connecting = 1
	Connected  = 2
)

func GetOutboundIP() string {
	return xnetwork.GetOutboundIPString()
}

const (
	DNSGoogle1    = "8.8.8.8"
	DNSGoogle2    = "8.8.8.8"
	DNSCloudflare = "1.1.1.1"
	DNS114A       = "114.114.114.114"
	DNS114B       = "114.114.115.115"
	DNSAlibaba1   = "223.5.5.5"
	DNSAlibaba2   = "223.6.6.6"
	DNSBaidu      = "180.76.76.76"
)

var cnDNSList = []string{DNS114A, DNSAlibaba1, DNSBaidu, DNS114B, DNSAlibaba2}
var otherDNSList = []string{DNSGoogle1, DNSCloudflare, DNSGoogle2}

func IsNetworkReachable() bool {
	return checkNetwork(cnDNSList...) || checkNetwork(otherDNSList...)
}

func checkNetwork(ips ...string) bool {
	for _, ip := range ips {
		conn, err := net.DialTimeout("tcp", ip+":80", time.Second*2)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}
