package message

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type OverflowPolicy uint8

const (
	OverflowBlock OverflowPolicy = iota
	OverflowDropNewest
	OverflowDropOldest
	OverflowDisconnect
)

var (
	ErrBusClosed          = errors.New("message bus is closed")
	ErrInvalidHandler     = errors.New("message handler is nil")
	ErrPublishTimeout     = errors.New("message delivery timed out")
	ErrSubscriberTooSlow  = errors.New("message subscriber is too slow")
	ErrSubscriberPanicked = errors.New("message subscriber panicked")
)

type SubscriptionOptions struct {
	BufferSize     int
	Overflow       OverflowPolicy
	EnqueueTimeout time.Duration
	OnError        func(error)
}

type PublishResult struct {
	Delivered    int
	Dropped      int
	Disconnected int
}

type Stats struct {
	Published    uint64
	Delivered    uint64
	Dropped      uint64
	Disconnected uint64
	Subscribers  int
}

type Bus[T any] struct {
	topicName string

	mu          sync.Mutex
	closed      atomic.Bool
	subscribers map[string]*subscriber[T]
	snapshot    atomic.Pointer[subscriberSnapshot[T]]

	published    atomic.Uint64
	delivered    atomic.Uint64
	dropped      atomic.Uint64
	disconnected atomic.Uint64
}

type subscriberRef[T any] struct {
	clientID   string
	subscriber *subscriber[T]
}

type subscriberSnapshot[T any] struct {
	subscribers []subscriberRef[T]
}

type subscriber[T any] struct {
	handler    func(T)
	options    SubscriptionOptions
	queue      chan T
	errorQueue chan error
	done       chan struct{}
	stopOnce   sync.Once
}
