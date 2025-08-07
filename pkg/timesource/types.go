package timesource

import (
	"sync"
	"time"
)

type Source struct {
	mu        sync.Mutex
	refTime   time.Time
	localTime time.Time
}
