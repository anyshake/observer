package timesource

import (
	"sync"
	"time"
)

type Source struct {
	rwMutex       sync.RWMutex
	ntpHost       string
	ntpPort       int
	queryTimeout  time.Duration
	LocalBaseTime time.Time
	ReferenceTime time.Time
}
