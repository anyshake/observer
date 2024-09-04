package transport

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"go.bug.st/serial"
)

type TransportDriverSerialImpl struct {
	conn  serial.Port
	mutex sync.Mutex
}

func (t *TransportDriverSerialImpl) Open(deps *TransportDependency) error {
	u, err := url.Parse(deps.DSN)
	if err != nil {
		return err
	}

	deviceName := u.Host
	if len(deviceName) == 0 {
		deviceName = u.Path
	}

	baudrate, err := strconv.Atoi(u.Query().Get("baudrate"))
	if err != nil {
		return err
	}

	conn, err := serial.Open(
		deviceName,
		&serial.Mode{
			BaudRate: baudrate,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		},
	)
	if err != nil {
		return fmt.Errorf("%v %s", err, deviceName)
	}

	conn.SetDTR(false)
	conn.SetRTS(false)

	t.conn = conn
	return nil
}

func (t *TransportDriverSerialImpl) Close() error {
	if t.conn == nil {
		return errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.conn.Close()
}

func (t *TransportDriverSerialImpl) Read(buf []byte, timeout time.Duration, flush bool) (int, error) {
	if t.conn == nil {
		return 0, errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if flush {
		t.conn.ResetInputBuffer()
	}
	t.conn.SetReadTimeout(timeout)
	return t.conn.Read(buf)
}

func (t *TransportDriverSerialImpl) Write(buf []byte, flush bool) (int, error) {
	if t.conn == nil {
		return 0, errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if flush {
		t.conn.ResetOutputBuffer()
	}
	return t.conn.Write(buf)
}

func (t *TransportDriverSerialImpl) Filter(signature []byte, timeout time.Duration) (bool, error) {
	if t.conn == nil {
		return false, errors.New("connection is not opened")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.conn.SetReadTimeout(timeout)

	t.conn.ResetInputBuffer()
	t.conn.ResetOutputBuffer()

	header := make([]byte, len(signature))
	_, err := t.conn.Read(header)
	if err != nil {
		return false, err
	}

	if bytes.Equal(header, signature) {
		return true, nil
	}

	return false, nil
}
