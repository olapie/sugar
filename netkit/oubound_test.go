package netkit_test

import (
	"testing"

	"code.olapie.com/sugar/v2/netkit"
)

func TestGetOutboundIPString(t *testing.T) {
	t.Log(netkit.GetOutboundIPString())
}
