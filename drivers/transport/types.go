package transport

import "time"

type TransportDependency struct {
	DSN    string
	Engine string
}

type TransportDriver interface {
	Open(deps *TransportDependency) error
	Close() error
	Read(buf []byte, timeout time.Duration, flush bool) (int, error)
	Write(buf []byte, flush bool) (int, error)
	Filter(signature []byte, timeout time.Duration) (bool, error)
}
