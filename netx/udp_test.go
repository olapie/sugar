package netx_test

import (
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"code.olapie.com/sugar/netx"
	"code.olapie.com/sugar/testx"
)

func TestBroadcastUDP(t *testing.T) {
	conn, err := net.ListenPacket("udp", ":9988")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	packet := testx.RandomBytes(10)
	buf := make([]byte, 200)
	received := make(chan error)

	go func() {
		nRead, addr, err := conn.ReadFrom(buf)
		if err == nil {
			t.Log(nRead, addr)
			buf = buf[:nRead]
		}
		received <- err
	}()

	err = netx.BroadcastUDP(9988, packet)
	testx.NoError(t, err)

	select {
	case err := <-received:
		testx.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("failed to receive udp packet")
	}

	testx.Equal(t, packet, buf)
	t.Log(packet)
	t.Log(buf)
}

func TestBroadcastUDPMaxPayload(t *testing.T) {
	t.Run("ExceedsLimitOfMTU", func(t *testing.T) {
		packet := testx.RandomBytes(1500)
		err := netx.BroadcastUDP(9988, packet)
		testx.Error(t, err)
		testx.True(t, strings.Contains(err.Error(), "write: message too long"))
	})
	t.Run("Success", func(t *testing.T) {
		packet := testx.RandomBytes(1300)
		err := netx.BroadcastUDP(9988, packet)
		testx.NoError(t, err)
	})
}
