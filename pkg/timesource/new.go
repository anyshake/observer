package timesource

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/beevik/ntp"
)

func NewNtpClient(endpoint string, retries int, timeout time.Duration) (*Source, error) {
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NTP endpoint: %w", err)
	}

	ntpAddress := urlObj.Host
	for i := 0; i <= retries; i++ {
		res, err := ntp.QueryWithOptions(ntpAddress, ntp.QueryOptions{Timeout: timeout})
		if err != nil {
			continue
		}
		return &Source{
			ntpAddress:    ntpAddress,
			queryRetries:  retries,
			queryTimeout:  timeout,
			LocalBaseTime: time.Now().UTC(),
			ReferenceTime: res.Time,
		}, nil
	}

	return nil, errors.New("failed to read time from NTP server")
}
