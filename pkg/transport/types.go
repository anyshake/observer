package transport

import (
	"context"
	"time"
)

type ITransport interface {
	Open() error
	Close() error
	Flush() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	GetLatency(packetSize int) time.Duration
	SetTimeout(timeout time.Duration) error
	ReadUntil(ctx context.Context, delim []byte, maxBytes int, timeout time.Duration) (dataBytes []byte, isTimeout bool, elapsed time.Duration, err error)
}
