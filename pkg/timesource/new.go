package timesource

import (
	"time"
)

func New(timeFunc TimeFunc) *Source {
	if timeFunc == nil {
		timeFunc = time.Now
	}

	currentTime := timeFunc()
	return &Source{
		timeFunc:  timeFunc,
		localTime: currentTime,
		refTime:   currentTime,
	}
}
