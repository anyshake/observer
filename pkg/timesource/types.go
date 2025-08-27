package timesource

import (
	"sync"
	"time"
)

type Source struct {
	mu        sync.RWMutex
	refTime   time.Time
	localTime time.Time
}
