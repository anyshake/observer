//go:build !darwin
// +build !darwin

package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bclswl0827/go-serial"
)

type SerialTransportImpl struct {
	baudrate int
	port     string
	timeout  time.Duration

	conn  *serial.Port
	mutex sync.Mutex
}

func (t *SerialTransportImpl) Open() error {
	conn, err := serial.Open(
		t.port,
		serial.WithHUPCL(true),
		serial.WithDataBits(8),
		serial.WithBaudrate(t.baudrate),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(int(t.timeout.Milliseconds())),
		serial.WithWriteTimeout(int(t.timeout.Milliseconds())),
	)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %w", t.port, err)
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

func (t *SerialTransportImpl) SetTimeout(timeout time.Duration) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.conn.SetReadTimeout(int(timeout.Milliseconds())); err != nil {
		return fmt.Errorf("failed to set read timeout: %w", err)
	}
	if err := t.conn.SetWriteTimeout(int(timeout.Milliseconds())); err != nil {
		return fmt.Errorf("failed to set write timeout: %w", err)
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

func (t *SerialTransportImpl) ReadUntil(ctx context.Context, delim []byte, maxBytes int, timeout time.Duration) ([]byte, bool, error) {
	if t.conn == nil {
		return nil, false, errors.New("connection is not opened")
	}

	deadline := time.Now().Add(timeout)
	buffer := make([]byte, 0, maxBytes)
	temp := make([]byte, 1)

	for {
		select {
		case <-ctx.Done():
			return nil, false, nil
		default:
		}

		if time.Now().After(deadline) {
			return nil, true, nil
		}

		t.mutex.Lock()
		n, err := t.conn.Read(temp)
		t.mutex.Unlock()

		if err != nil {
			return nil, false, fmt.Errorf("read error: %w", err)
		}
		if n == 0 {
			continue
		}

		buffer = append(buffer, temp[0])
		if len(buffer) > maxBytes {
			return nil, false, fmt.Errorf("read exceeded maxBytes (%d) before delimiter", maxBytes)
		}

		if len(buffer) >= len(delim) && bytes.HasSuffix(buffer, delim) {
			return buffer, false, nil
		}
	}
}
