package netx_test

import (
	"code.olapie.com/sugar/netx"
	"testing"
)

func TestGetOutboundIPString(t *testing.T) {
	t.Log(netx.GetOutboundIPString())
}
