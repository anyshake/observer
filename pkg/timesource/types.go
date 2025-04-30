package timesource

import (
	"sync"
	"time"
)

type Source struct {
	rwMutex       sync.RWMutex
	ntpAddress    string
	queryRetries  int
	queryTimeout  time.Duration
	LocalBaseTime time.Time
	ReferenceTime time.Time
}
