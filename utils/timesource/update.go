package timesource

import (
	"errors"
	"time"

	"github.com/beevik/ntp"
)

func (s *Source) Update() error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	for i := 0; i < s.queryAttempts; i++ {
		res, err := ntp.QueryWithOptions(s.ntpHost, ntp.QueryOptions{
			Port: s.ntpPort, Timeout: s.queryTimeout,
		})
		if err != nil {
			continue
		}
		s.LocalBaseTime = time.Now().UTC()
		s.ReferenceTime = res.Time
		return nil
	}

	return errors.New("failed to update time from NTP server")
}
