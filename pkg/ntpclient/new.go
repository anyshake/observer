package ntpclient

import (
	"fmt"
	"net/url"
	"time"
)

func New(endpoint string, retries int, timeout int) (Client, error) {
	if len(endpoint) == 0 {
		return Client{}, fmt.Errorf("NTP endpoint is empty")
	}

	urlObj, err := url.Parse(endpoint)
	if err != nil {
		return Client{}, fmt.Errorf("failed to parse NTP endpoint: %w", err)
	}

	return Client{
		ntpAddr:     urlObj.Host,
		retries:     retries,
		readTimeout: time.Duration(timeout) * time.Second,
	}, nil
}
