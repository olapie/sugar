package xnetwork_test

import (
	"testing"

	"code.olapie.com/sugar/xnetwork"
)

func TestGetOutboundIPString(t *testing.T) {
	t.Log(xnetwork.GetOutboundIPString())
}
