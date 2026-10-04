package transport

import (
	"testing"
	"time"
)

func TestNewTCPTransport(t *testing.T) {
	t.Parallel()
	tests := []struct {
		endpoint string
		timeout  int
		wantHost string
		wantTTL  time.Duration
	}{
		{"tcp://127.0.0.1:1234", 2, "127.0.0.1:1234", 2 * time.Second},
		{"tcp://example.test:5000", 0, "example.test:5000", 5 * time.Second},
		{"tcp://[2001:db8::1]:5000", 7, "[2001:db8::1]:5000", 7 * time.Second},
	}
	for _, tt := range tests {
		transport, err := New(tt.endpoint, tt.timeout)
		if err != nil {
			t.Fatalf("New(%q) error = %v", tt.endpoint, err)
		}
		tcp, ok := transport.(*TcpTransportImpl)
		if !ok || tcp.host != tt.wantHost || tcp.timeout != tt.wantTTL {
			t.Fatalf("New(%q) = %#v, want host %q and timeout %v", tt.endpoint, transport, tt.wantHost, tt.wantTTL)
		}
	}
}

func TestNewSerialTransport(t *testing.T) {
	t.Parallel()
	tests := []struct {
		endpoint string
		wantPort string
		wantBaud int
	}{
		{"serial:///dev/ttyUSB0?baudrate=115200", "/dev/ttyUSB0", 115200},
		{"serial://COM3?baudrate=9600", "COM3", 9600},
	}
	for _, tt := range tests {
		transport, err := New(tt.endpoint, 3)
		if err != nil {
			t.Fatalf("New(%q) error = %v", tt.endpoint, err)
		}
		serial, ok := transport.(*SerialTransportImpl)
		if !ok || serial.port != tt.wantPort || serial.baudrate != tt.wantBaud || serial.timeout != 3*time.Second {
			t.Fatalf("New(%q) = %#v", tt.endpoint, transport)
		}
	}
}

func TestNewTransportRejectsInvalidEndpoints(t *testing.T) {
	t.Parallel()
	for _, endpoint := range []string{"", "://bad", "udp://host:1", "serial:///dev/ttyUSB0", "serial:///dev/ttyUSB0?baudrate=bad"} {
		if transport, err := New(endpoint, 1); err == nil || transport != nil {
			t.Errorf("New(%q) = %#v, %v, want error", endpoint, transport, err)
		}
	}
}

func TestClosedTCPTransportOperations(t *testing.T) {
	t.Parallel()
	transport := &TcpTransportImpl{}
	if err := transport.Close(); err == nil {
		t.Error("Close() succeeded before Open()")
	}
	if n, err := transport.Read(make([]byte, 1)); err == nil || n != 0 {
		t.Errorf("Read() = %d, %v before Open()", n, err)
	}
	if n, err := transport.Write([]byte{1}); err == nil || n != 0 {
		t.Errorf("Write() = %d, %v before Open()", n, err)
	}
	if _, _, _, err := transport.ReadUntil(t.Context(), 1, func(*[]byte, *time.Time) bool { return true }, time.Second); err == nil {
		t.Error("ReadUntil() succeeded before Open()")
	}
	if err := transport.SetTimeout(2 * time.Second); err != nil || transport.timeout != 2*time.Second {
		t.Fatalf("SetTimeout() = %v, timeout = %v", err, transport.timeout)
	}
	if transport.GetLatency(100) != 0 || transport.Flush() != nil {
		t.Fatal("closed TCP transport returned latency or flush error")
	}
}
