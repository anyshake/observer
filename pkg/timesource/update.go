package timesource

import (
	"time"
)

func (s *Source) Update(localTime, refTime time.Time, driftPPM float64, timeFunc TimeFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if timeFunc != nil {
		s.timeFunc = timeFunc
	}

	s.driftPPM = driftPPM

	s.localTime = localTime
	s.refTime = refTime.UTC() // strip monotonic clock data
}
