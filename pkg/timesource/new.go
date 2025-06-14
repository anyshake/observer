package timesource

import (
	"fmt"
	"net/url"
	"time"
)

func NewNtpClient(endpoint string, retries int, timeout time.Duration) (*Source, error) {
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NTP endpoint: %w", err)
	}

	ntpAddress := urlObj.Host

	avgOffset, err := measureStableOffset(ntpAddress, retries, timeout, 10)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Source{
		ntpAddress:    ntpAddress,
		queryRetries:  retries,
		queryTimeout:  timeout,
		LocalBaseTime: now,
		ReferenceTime: now.Add(avgOffset),
	}, nil
}
