package cache

import (
	"sync"
	"time"
)

type AnyCache struct {
	cache     any
	ttl       time.Duration
	createdAt time.Time
	mutex     sync.RWMutex
}

type KvCache[T any] struct {
	cache     map[any]T
	ttl       time.Duration
	createdAt time.Time
	mutex     sync.RWMutex
}
