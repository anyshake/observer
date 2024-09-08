package timesource

import (
	"errors"
	"time"

	"github.com/beevik/ntp"
)

func New(ntpHost string, ntpPort, retries int, timeout time.Duration) (Source, error) {
	for i := 0; i <= retries; i++ {
		res, err := ntp.QueryWithOptions(ntpHost, ntp.QueryOptions{
			Port: ntpPort, Timeout: timeout,
		})
		if err != nil {
			continue
		}
		return Source{
			ntpHost:       ntpHost,
			ntpPort:       ntpPort,
			queryRetries:  retries,
			queryTimeout:  timeout,
			LocalBaseTime: time.Now().UTC(),
			ReferenceTime: res.Time,
		}, nil
	}

	return Source{}, errors.New("failed to read time from NTP server")
}
