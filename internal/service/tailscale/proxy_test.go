package tailscale

import (
	"context"
	"net"
	"slices"
	"testing"
	"time"
)

func TestNormalizePortForwardSpecs(t *testing.T) {
	normalized, err := normalizePortForwardSpecs([]string{
		" 123/UDP ",
		"80/TCP",
		"22/tcp",
		"080/tcp",
		"123/udp",
	})
	if err != nil {
		t.Fatalf("failed to normalize forward ports: %v", err)
	}

	expected := []string{"22/tcp", "80/tcp", "123/udp"}
	if !slices.Equal(normalized, expected) {
		t.Fatalf("unexpected normalized ports: got %#v, want %#v", normalized, expected)
	}
}

func TestParsePortForwardSpecRejectsInvalidValues(t *testing.T) {
	invalid := []string{
		"",
		"80",
		"tcp/80",
		"0/tcp",
		"65536/udp",
		"80/http",
		"80/tcp/foo",
		"80 /tcp",
		"80/ udp",
	}
	for _, value := range invalid {
		t.Run(value, func(t *testing.T) {
			if _, err := parsePortForwardSpec(value); err == nil {
				t.Fatalf("expected %q to be rejected", value)
			}
		})
	}
}

func TestBuildForwardsIncludesHTTPAndConfiguredProtocols(t *testing.T) {
	forwards, err := buildForwards("0.0.0.0:8073", []string{
		"22/tcp",
		"123/udp",
		"8073/tcp",
		"8073/udp",
	})
	if err != nil {
		t.Fatalf("failed to build forwards: %v", err)
	}

	expectedTCP := []tcpForward{
		{listenAddr: ":8073", targetAddr: "127.0.0.1:8073"},
		{listenAddr: ":22", targetAddr: "localhost:22"},
	}
	if !slices.Equal(forwards.tcp, expectedTCP) {
		t.Fatalf("unexpected TCP forwards: got %#v, want %#v", forwards.tcp, expectedTCP)
	}

	expectedUDP := []udpForward{
		{port: 123, targetAddr: "localhost:123"},
		{port: 8073, targetAddr: "localhost:8073"},
	}
	if !slices.Equal(forwards.udp, expectedUDP) {
		t.Fatalf("unexpected UDP forwards: got %#v, want %#v", forwards.udp, expectedUDP)
	}
}

func TestServeUDPProxyForwardsRequestsAndResponses(t *testing.T) {
	target, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create UDP target: %v", err)
	}
	defer target.Close()

	echoDone := make(chan struct{})
	go func() {
		defer close(echoDone)
		buffer := make([]byte, 1024)
		n, addr, readErr := target.ReadFrom(buffer)
		if readErr == nil {
			_, _ = target.WriteTo(buffer[:n], addr)
		}
	}()

	listener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create UDP proxy listener: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	proxyDone := make(chan error, 1)
	go func() {
		proxyDone <- serveUDPProxy(ctx, listener, target.LocalAddr().String(), &testLogger{})
	}()

	client, err := net.Dial("udp", listener.LocalAddr().String())
	if err != nil {
		cancel()
		t.Fatalf("failed to create UDP client: %v", err)
	}
	defer client.Close()
	if err := client.SetDeadline(time.Now().Add(time.Second)); err != nil {
		cancel()
		t.Fatalf("failed to set UDP client deadline: %v", err)
	}

	payload := []byte("udp forwarding test")
	if _, err := client.Write(payload); err != nil {
		cancel()
		t.Fatalf("failed to send UDP request: %v", err)
	}
	response := make([]byte, len(payload))
	n, err := client.Read(response)
	if err != nil {
		cancel()
		t.Fatalf("failed to read UDP response: %v", err)
	}
	if string(response[:n]) != string(payload) {
		cancel()
		t.Fatalf("unexpected UDP response: got %q, want %q", response[:n], payload)
	}

	cancel()
	if err := <-proxyDone; err != nil {
		t.Fatalf("UDP proxy returned an error during shutdown: %v", err)
	}
	<-echoDone
}
