package netkit_test

import (
	"net"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/netkit"
	"code.olapie.com/sugar/v2/testutil"
)

func TestMulticast(t *testing.T) {
	ifi := netkit.GetPrivateIPv4Interface()
	if ifi == nil {
		t.Log("No PrivateIPv4Interface")
		t.FailNow()
	}
	multiIP := netkit.GetMulticastIPv4String(ifi)
	if multiIP == "" {
		t.Fatal("no multi ip")
	}
	udpAddr, err := net.ResolveUDPAddr("udp", multiIP+":9999")
	if err != nil {
		t.Fatal(err, multiIP)
	}
	t.Log(udpAddr.String())
	conn, err := net.ListenMulticastUDP("udp", ifi, udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	packet := testutil.RandomBytes(10)
	buf := make([]byte, 2000)
	received := make(chan error)

	go func() {
		nRead, addr, err := conn.ReadFrom(buf)
		if err == nil {
			t.Log(nRead, addr)
			buf = buf[:nRead]
		}
		received <- err
	}()

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Fatal(err)
	}

	_, err = udpConn.Write(packet)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
	select {
	case err := <-received:
		testutil.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("failed to receive udp packet")
	}

	testutil.Equal(t, packet, buf)
	t.Log(packet)
	t.Log(buf)
}

func TestGetBroadcastIPv4(t *testing.T) {
	ifi := netkit.GetPrivateIPv4Interface()
	if ifi == nil {
		t.Log("No PrivateIPv4Interface")
		t.FailNow()
	}
	ip := netkit.GetBroadcastIPv4(ifi)
	t.Log(ip.String())

	udpAddr, err := net.ResolveUDPAddr("udp", ip.String()+":7819")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write(testutil.RandomBytes(10))
	if err != nil {
		t.Fatal(err)
	}
}
