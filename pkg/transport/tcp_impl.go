package transport

import (
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

func (t *TcpTransportImpl) SetTimeout(timeout time.Duration) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.timeout = timeout
	return nil
}
