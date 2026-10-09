package transport

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestTCPTransportReadsAndWrites(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4)
		_, _ = io.ReadFull(conn, buf)
		_, _ = conn.Write([]byte("OKZ"))
	}()

	client, err := New("tcp://"+listener.Addr().String(), 1)
	if err != nil {
		t.Fatal(err)
	}
	tcp := client.(*TcpTransportImpl)
	if _, _, _, err := tcp.ReadUntil(context.Background(), 1, nil, time.Second); err == nil {
		t.Fatal("nil callback accepted before open")
	}
	if err := client.Open(); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := client.ReadUntil(context.Background(), 1, nil, time.Second); err == nil {
		t.Fatal("nil callback accepted")
	}
	if _, err := client.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 2)
	if _, err := client.Read(buf); err != nil || string(buf) != "OK" {
		t.Fatalf("read = %q, %v", buf, err)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if data, timedOut, _, err := client.ReadUntil(canceled, 4, func(*[]byte, *time.Time) bool { return true }, time.Second); err != nil || timedOut || data != nil {
		t.Fatalf("canceled read = %q timedOut=%v err=%v", data, timedOut, err)
	}

	data, timedOut, elapsed, err := client.ReadUntil(context.Background(), 4, func(buf *[]byte, updated *time.Time) bool {
		*updated = time.Now()
		return len(*buf) == 1
	}, time.Second)
	if err != nil || timedOut || string(data) != "Z" || elapsed < 0 {
		t.Fatalf("ReadUntil = %q timedOut=%v elapsed=%v err=%v", data, timedOut, elapsed, err)
	}
	if err := client.Flush(); err != nil {
		t.Fatalf("Flush() error = %v", err)
	}
	// Linux reports the measured TCP RTT; loopback is not necessarily zero.
	if latency := client.GetLatency(8); latency < 0 {
		t.Fatalf("GetLatency() = %v, want a non-negative duration", latency)
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}

	overflowListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer overflowListener.Close()
	go func() {
		conn, acceptErr := overflowListener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write([]byte("XY"))
		time.Sleep(100 * time.Millisecond)
	}()
	overflow, err := New("tcp://"+overflowListener.Addr().String(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := overflow.Open(); err != nil {
		t.Fatal(err)
	}
	defer overflow.Close()
	_ = overflow.SetTimeout(20 * time.Millisecond)
	if _, _, _, err := overflow.ReadUntil(context.Background(), 0, func(*[]byte, *time.Time) bool { return false }, time.Second); err == nil {
		t.Fatal("maxBytes overflow was accepted")
	}
	if _, timedOut, _, err := overflow.ReadUntil(context.Background(), 4, func(*[]byte, *time.Time) bool { return false }, 30*time.Millisecond); err != nil || !timedOut {
		t.Fatalf("timeout read timedOut=%v err=%v", timedOut, err)
	}

	refused, err := New("tcp://127.0.0.1:1", 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := refused.Open(); err == nil {
		t.Fatal("dial to a closed port succeeded")
	}
}
