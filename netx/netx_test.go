package netx_test

import (
	"testing"

	"code.olapie.com/sugar/netx"
)

func TestGetOutboundIPString(t *testing.T) {
	t.Log(netx.GetOutboundIPString())
}

func TestGetIFaceNames(t *testing.T) {
	t.Log(netx.GetIFaceNames())
	names := netx.GetIFaceNames()
	for _, name := range names {
		t.Log(name, netx.GetIPv4(name))
	}
}

func TestGetLocalIP4String(t *testing.T) {
	t.Log(netx.GetLocalIP4String())
}
