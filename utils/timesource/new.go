package timesource

import (
	"errors"
	"time"

	"github.com/beevik/ntp"
)

func New(ntpHost string, ntpPort, attempts int, timeout time.Duration) (Source, error) {
	for i := 0; i < attempts; i++ {
		res, err := ntp.QueryWithOptions(ntpHost, ntp.QueryOptions{
			Port: ntpPort, Timeout: timeout,
		})
		if err != nil {
			continue
		}
		return Source{
			ntpHost:       ntpHost,
			ntpPort:       ntpPort,
			queryAttempts: attempts,
			queryTimeout:  timeout,
			LocalBaseTime: time.Now().UTC(),
			ReferenceTime: res.Time,
		}, nil
	}

	return Source{}, errors.New("failed to read time from NTP server")
}
