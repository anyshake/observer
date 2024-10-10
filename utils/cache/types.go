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
