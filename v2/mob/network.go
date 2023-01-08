package mob

import "code.olapie.com/sugar/v2/xnetwork"

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
