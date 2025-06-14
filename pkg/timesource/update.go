package timesource

import (
	"fmt"
	"time"
)

func (s *Source) Update() error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	avgOffset, err := measureStableOffset(s.ntpAddress, s.queryRetries, s.queryTimeout, 10)
	if err != nil {
		return fmt.Errorf("failed to update time from NTP server: %w", err)
	}

	s.LocalBaseTime = time.Now().UTC()
	s.ReferenceTime = s.LocalBaseTime.Add(avgOffset)
	return nil
}
