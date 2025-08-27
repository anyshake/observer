package ntpclient

import (
	"fmt"
	"net/url"
	"time"
)

func New(pool []string, retries int, timeout int) (Client, error) {
	if len(pool) == 0 {
		return Client{}, fmt.Errorf("NTP pool is empty")
	}

	hosts := make([]string, len(pool))
	for _, host := range pool {
		urlObj, err := url.Parse(host)
		if err != nil {
			return Client{}, fmt.Errorf("failed to parse NTP endpoint: %w", err)
		}
		hosts = append(hosts, urlObj.Host)
	}

	return Client{
		pool:        hosts,
		retries:     retries,
		readTimeout: time.Duration(timeout) * time.Second,
	}, nil
}
