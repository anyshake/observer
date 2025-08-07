package timesource

import (
	"time"
)

func New() *Source {
	currentTime := time.Now()
	return &Source{
		localTime: currentTime,
		refTime:   currentTime,
	}
}
