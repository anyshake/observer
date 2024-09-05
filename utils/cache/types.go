package cache

import (
	"sync"
	"time"
)

type BytesCache struct {
	cache     []byte
	ttl       time.Duration
	createdAt time.Time
	mutex     sync.RWMutex
}
