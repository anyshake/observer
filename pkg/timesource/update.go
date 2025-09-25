package timesource

import (
	"time"
)

func (s *Source) Update(localTime, refTime time.Time, timeFunc TimeFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if timeFunc != nil {
		s.timeFunc = timeFunc
	}

	s.localTime = localTime
	s.refTime = refTime.UTC() // strip monotonic clock data
}
