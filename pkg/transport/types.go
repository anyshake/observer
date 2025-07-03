package transport

import "time"

type ITransport interface {
	Open() error
	Close() error
	Flush() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	SetTimeout(timeout time.Duration) error
	ReadUntil(delim []byte, maxBytes int, timeout time.Duration) (dataBytes []byte, isTimeout bool, err error)
}
