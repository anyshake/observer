package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
)

type SerialTransportImpl struct {
	baudrate int
	port     string
	timeout  time.Duration

	conn  serial.Port
	mutex sync.Mutex
}

func (t *SerialTransportImpl) Open() error {
	conn, err := serial.Open(
		t.port,
		&serial.Mode{
			BaudRate: t.baudrate,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %w", t.port, err)
	}
	if err = conn.SetReadTimeout(t.timeout); err != nil {
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	if err = conn.ResetInputBuffer(); err != nil {
		return fmt.Errorf("failed to reset input buffer: %w", err)
	}
	if err = conn.ResetOutputBuffer(); err != nil {
		return fmt.Errorf("failed to reset output buffer: %w", err)
	}

	t.conn = conn
	return nil
}

func (t *SerialTransportImpl) Close() error {
	if t.conn == nil {
		return errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.conn.Close()
}

func (t *SerialTransportImpl) GetLatency(packetSize int) time.Duration {
	totalBits := packetSize * (8 + 1) // 8 bits data + 1 bit stop bit
	return time.Duration(float64(totalBits) * float64(time.Second) / float64(t.baudrate))
}

func (t *SerialTransportImpl) SetTimeout(timeout time.Duration) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.conn.SetReadTimeout(t.timeout); err != nil {
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	t.timeout = timeout
	return nil
}

func (t *SerialTransportImpl) Read(buf []byte) (int, error) {
	if t.conn == nil {
		return 0, errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.conn.Read(buf)
}

func (t *SerialTransportImpl) Write(buf []byte) (int, error) {
	if t.conn == nil {
		return 0, errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.conn.Write(buf)
}

func (t *SerialTransportImpl) Flush() error {
	if t.conn == nil {
		return errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.conn.ResetInputBuffer(); err != nil {
		return err
	}
	if err := t.conn.ResetOutputBuffer(); err != nil {
		return err
	}

	return nil
}

func (t *SerialTransportImpl) ReadUntil(ctx context.Context, delim []byte, maxBytes int, timeout time.Duration) ([]byte, bool, time.Duration, error) {
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
