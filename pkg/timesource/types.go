package timesource

import (
	"sync"
	"time"
)

type TimeFunc func() time.Time

type Source struct {
	mu       sync.RWMutex
	timeFunc TimeFunc

	refTime   time.Time
	localTime time.Time
}
