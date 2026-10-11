package ntp_server

import (
	"context"
	"net"
	"net/netip"
	"testing"
	"time"
)

func TestPacketValidationAndEncoding(t *testing.T) {
	fixed := time.Unix(1_700_000_000, 123456789).UTC()
	server := &ntpServer{timeFn: func() time.Time { return fixed }}
	if server.workerCount() < 4 || server.workerCount() > 32 {
		t.Fatalf("worker count = %d", server.workerCount())
	}

	loopback := netip.MustParseAddr("127.0.0.1")
	public := netip.MustParseAddr("8.8.8.8")
	request := make([]byte, 48)
	request[0] = 0x23
	request[2] = 6
	copy(request[40:], []byte{1, 2, 3, 4, 5, 6, 7, 8})
	if !server.isDataValid(loopback, request) || !server.isDataValid(public, request) {
		t.Fatal("client packet was rejected")
	}
	request[0] = 0xE3
	if !server.isDataValid(public, request) {
		t.Fatal("alarm client packet was rejected")
	}
	request[0] = 0x21
	if !server.isDataValid(loopback, request) {
		t.Fatal("symmetric packet from loopback was rejected")
	}
	if server.isDataValid(public, request) {
		t.Fatal("symmetric packet from a public address was accepted")
	}
	request[0] = 0x23
	if server.isDataValid(loopback, request[:10]) {
		t.Fatal("short packet was accepted")
	}
	request[0] = 0x2B
	if server.isDataValid(loopback, request) {
		t.Fatal("unsupported version was accepted")
	}
	request[0] = 0x20
	if server.isDataValid(loopback, request) {
		t.Fatal("unsupported mode was accepted")
	}
	request[0] = 0x63
	if server.isDataValid(loopback, request) {
		t.Fatal("unsupported leap indicator was accepted")
	}

	request[0] = 0x23
	response, err := server.encodePacket(fixed, request)
	if err != nil || len(response) != 48 {
		t.Fatalf("encoded %d bytes, %v", len(response), err)
	}
	if response[0] != 0x24 || response[1] != 1 || response[2] != request[2] || string(response[12:16]) != "GNSS" {
		t.Fatalf("header = %x", response[:16])
	}
	if string(response[24:32]) != string(request[40:48]) {
		t.Fatal("originate timestamp was not copied")
	}

	if err := server.handleInvalidAddress(); err != nil {
		t.Fatal(err)
	}
}

func (p *ntpServer) handleInvalidAddress() error {
	p.handleData(p.timeFn(), &net.UDPAddr{Port: 1}, make([]byte, 48))
	return nil
}

func TestNTPServerExchange(t *testing.T) {
	fixed := time.Unix(1_700_000_000, 0).UTC()
	opened, err := newNtpServer("127.0.0.1", 0, nil)
	if err != nil || opened.timeFn == nil {
		t.Fatal(err)
	}
	_ = opened.conn.Close()
	if _, err := newNtpServer("127.0.0.1:1", 1, func() time.Time { return fixed }); err == nil {
		t.Fatal("invalid listen address was accepted")
	}

	server, err := newNtpServer("127.0.0.1", 0, func() time.Time { return fixed })
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- server.Run(ctx) }()

	client, err := net.DialUDP("udp", nil, server.conn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	request := make([]byte, 48)
	request[0] = 0x23
	if _, err := client.Write(request); err != nil {
		t.Fatal(err)
	}
	_ = client.SetReadDeadline(time.Now().Add(2 * time.Second))
	response := make([]byte, 48)
	if _, err := client.Read(response); err != nil || response[0] != 0x24 {
		t.Fatalf("response = %x, %v", response[:4], err)
	}
	request[0] = 0x20
	if _, err := client.Write(request); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Write([]byte{0x23}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	_ = server.write([]byte{1}, nil)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("NTP server did not stop")
	}
}
