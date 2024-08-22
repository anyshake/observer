package timesource

import (
	"time"

	"github.com/beevik/ntp"
)

func New(ntpHost string, ntpPort int, queryTimeout time.Duration) (Source, error) {
	res, err := ntp.QueryWithOptions(ntpHost, ntp.QueryOptions{
		Port: ntpPort, Timeout: queryTimeout,
	})
	if err != nil {
		return Source{}, err
	}

	return Source{
		ntpHost:       ntpHost,
		ntpPort:       ntpPort,
		queryTimeout:  queryTimeout,
		LocalBaseTime: time.Now().UTC(),
		ReferenceTime: res.Time,
	}, nil
}
