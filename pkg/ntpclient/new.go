package ntpclient

import (
	"fmt"
	"net/url"
	"time"

	"github.com/bclswl0827/ntp"
)

func New(pool []string, retries int, timeout int, timeFunc ntp.TimeFunc) (Client, error) {
	if len(pool) == 0 {
		return Client{}, fmt.Errorf("NTP pool is empty")
	}

	if timeFunc == nil {
		timeFunc = time.Now
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
		timeFunc:    timeFunc,
		pool:        hosts,
		retries:     retries,
		readTimeout: time.Duration(timeout) * time.Second,
	}, nil
}
