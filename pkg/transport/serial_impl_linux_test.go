package transport

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func openPseudoTerminal(t *testing.T) (*os.File, string) {
	t.Helper()
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Skipf("pseudo terminal is unavailable: %v", err)
	}
	t.Cleanup(func() { _ = master.Close() })
	if err := unix.IoctlSetPointerInt(int(master.Fd()), unix.TIOCSPTLCK, 0); err != nil {
		t.Fatalf("failed to unlock pseudo terminal: %v", err)
	}
	number, err := unix.IoctlGetInt(int(master.Fd()), unix.TIOCGPTN)
	if err != nil {
		t.Fatalf("failed to get pseudo terminal number: %v", err)
	}
	return master, fmt.Sprintf("/dev/pts/%d", number)
}

func TestSerialTransportLifecycle(t *testing.T) {
	master, port := openPseudoTerminal(t)
	transport := &SerialTransportImpl{port: port, baudrate: 115200, timeout: time.Second}

	if err := transport.SetTimeout(50 * time.Millisecond); err != nil || transport.timeout != 50*time.Millisecond {
		t.Fatalf("SetTimeout() before Open() = %v, timeout = %v", err, transport.timeout)
	}
	if _, err := transport.Read(make([]byte, 1)); err == nil {
		t.Error("Read() succeeded before Open()")
	}
	if _, err := transport.Write([]byte{1}); err == nil {
		t.Error("Write() succeeded before Open()")
	}
	if err := transport.Flush(); err == nil {
		t.Error("Flush() succeeded before Open()")
	}
	if err := transport.Close(); err == nil {
		t.Error("Close() succeeded before Open()")
	}

	if err := transport.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := transport.Open(); err == nil {
		t.Fatal("second Open() succeeded")
	}
	if _, err := master.Write([]byte("A")); err != nil {
		t.Fatal(err)
	}
	data, timedOut, _, err := transport.ReadUntil(context.Background(), 4, func(buf *[]byte, _ *time.Time) bool {
		return len(*buf) == 1
	}, time.Second)
	if err != nil || timedOut || string(data) != "A" {
		t.Fatalf("ReadUntil() = %q, timedOut=%v, err=%v", data, timedOut, err)
	}
	if err := transport.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := transport.Close(); err == nil {
		t.Fatal("second Close() succeeded")
	}
	if _, err := transport.Read(make([]byte, 1)); err == nil {
		t.Fatal("Read() succeeded after Close()")
	}
	if err := transport.Open(); err != nil {
		t.Fatalf("Open() after Close() error = %v", err)
	}
	if err := transport.Close(); err != nil {
		t.Fatalf("Close() after reopen error = %v", err)
	}
}

func TestSerialTransportCloseDuringReadUntil(t *testing.T) {
	master, port := openPseudoTerminal(t)
	transport := &SerialTransportImpl{port: port, baudrate: 115200, timeout: 10 * time.Millisecond}
	if err := transport.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if _, err := master.Write([]byte("A")); err != nil {
		t.Fatal(err)
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
