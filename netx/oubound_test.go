package netx_test

import (
	"testing"

	"code.olapie.com/sugar/netx"
)

func TestGetOutboundIPString(t *testing.T) {
	t.Log(netx.GetOutboundIPString())
}
