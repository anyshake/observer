package timesource

import (
	"time"
)

func (s *Source) Fix(currentTime, prevTime time.Time, span time.Duration) time.Time {
	expectedTime := prevTime.Add(span)
	discrepancy := expectedTime.Sub(currentTime)

	return currentTime.Add(discrepancy)
}
