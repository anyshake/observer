package timesource

import (
	"time"

	"github.com/beevik/ntp"
)

func (s *Source) Update() error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	res, err := ntp.QueryWithOptions(s.ntpHost, ntp.QueryOptions{
		Port: s.ntpPort, Timeout: s.queryTimeout,
	})
	if err != nil {
		return err
	}

	s.LocalBaseTime = time.Now().UTC()
	s.ReferenceTime = res.Time
	return nil
}
