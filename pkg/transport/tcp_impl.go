//go:build !linux
// +build !linux

package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type TcpTransportImpl struct {
	host    string
	conn    net.Conn
	mutex   sync.Mutex
	timeout time.Duration
}

func (t *TcpTransportImpl) Open() error {
	conn, err := net.Dial("tcp", t.host)
	if err != nil {
		return fmt.Errorf("failed to open TCP connection to %s: %w", t.host, err)
	}

	t.conn = conn
	return nil
}

func (t *TcpTransportImpl) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.conn.Close()
}

func (t *TcpTransportImpl) Read(buf []byte) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.conn.SetReadDeadline(time.Now().Add(t.timeout)); err != nil {
		return 0, fmt.Errorf("failed to set read timeout: %w", err)
	}

	return t.conn.Read(buf)
}

func (t *TcpTransportImpl) Write(buf []byte) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.conn.SetWriteDeadline(time.Now().Add(t.timeout)); err != nil {
		return 0, fmt.Errorf("failed to set write timeout: %w", err)
	}

	return t.conn.Write(buf)
}

func (t *TcpTransportImpl) Flush() error {
	return nil
}

func (t *TcpTransportImpl) GetLatency(packetSize int) time.Duration {
	return 0
}

func (t *TcpTransportImpl) SetTimeout(timeout time.Duration) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.timeout = timeout
	return nil
}

func (t *TcpTransportImpl) ReadUntil(ctx context.Context, delim []byte, maxBytes int, timeout time.Duration) ([]byte, bool, time.Duration, error) {
	if t.conn == nil {
		return nil, false, 0, errors.New("connection is not opened")
	}

	deadline := time.Now().Add(timeout)
	buffer := make([]byte, 0, maxBytes)
	temp := make([]byte, 1)

	var firstByteTime time.Time
	var lastByteTime time.Time

	for {
		select {
		case <-ctx.Done():
			return nil, false, 0, nil
		default:
		}

		if time.Now().After(deadline) {
			return nil, true, 0, nil
		}

		t.mutex.Lock()
		n, err := t.conn.Read(temp)
		currentTime := time.Now()
		t.mutex.Unlock()

		if err != nil {
			return nil, false, 0, fmt.Errorf("read error: %w", err)
		}
		if n == 0 {
			continue
		}

		if firstByteTime.IsZero() {
			firstByteTime = currentTime
		}
		lastByteTime = currentTime

		buffer = append(buffer, temp[0])
		if len(buffer) > maxBytes {
			return nil, false, 0, fmt.Errorf("read exceeded maxBytes (%d) before delimiter", maxBytes)
		}

		if len(buffer) >= len(delim) && bytes.HasSuffix(buffer, delim) {
			duration := lastByteTime.Sub(firstByteTime)
			return buffer, false, duration, nil
		}
	}
}
