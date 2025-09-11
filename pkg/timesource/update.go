package timesource

import (
	"time"
)

func (s *Source) Update(localTime, refTime time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.localTime = localTime
	s.refTime = refTime.UTC() // strip monotonic clock data
}
