package transport

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func startSilentTCPServer(t *testing.T, greeting []byte) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var (
		mu    sync.Mutex
		conns []net.Conn
	)
	accepted := make(chan struct{})
	go func() {
		defer close(accepted)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			conns = append(conns, conn)
			mu.Unlock()
			if len(greeting) > 0 {
				_, _ = conn.Write(greeting)
			}
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		<-accepted
		mu.Lock()
		defer mu.Unlock()
		for _, conn := range conns {
			_ = conn.Close()
		}
	})
	return listener.Addr().String()
}

func TestTCPTransportRejectsSecondOpenAndClosesOnce(t *testing.T) {
	transport := &TcpTransportImpl{host: startSilentTCPServer(t, nil), timeout: time.Second}
	if err := transport.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := transport.Open(); err == nil {
		t.Fatal("second Open() succeeded and replaced the active connection")
	}
	if err := transport.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := transport.Close(); err == nil {
		t.Fatal("second Close() succeeded")
	}
	if err := transport.Open(); err != nil {
		t.Fatalf("Open() after Close() error = %v", err)
	}
	if err := transport.Close(); err != nil {
		t.Fatalf("Close() after reopen error = %v", err)
	}
}

func TestTCPTransportCloseDuringReadUntil(t *testing.T) {
	transport := &TcpTransportImpl{host: startSilentTCPServer(t, []byte{0x01}), timeout: 10 * time.Millisecond}
	if err := transport.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	reading := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		var once sync.Once
		_, _, _, err := transport.ReadUntil(context.Background(), 16, func(*[]byte, *time.Time) bool {
			once.Do(func() { close(reading) })
			return false
		}, 5*time.Second)
		result <- err
	}()
	<-reading
	if err := transport.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("ReadUntil() returned nil error after Close()")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ReadUntil() kept running after Close()")
	}
}

func TestTCPTransportGetLatencyDuringClose(t *testing.T) {
	transport := &TcpTransportImpl{host: startSilentTCPServer(t, nil), timeout: time.Second}
	if err := transport.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	done := make(chan struct{})
	started := make(chan struct{})
	var readers sync.WaitGroup
	readers.Go(func() {
		_ = transport.GetLatency(1)
		close(started)
		for {
			select {
			case <-done:
				return
			default:
				_ = transport.GetLatency(1)
			}
		}
	})
	<-started
	if err := transport.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	close(done)
	readers.Wait()
}
