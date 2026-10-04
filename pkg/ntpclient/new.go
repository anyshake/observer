package ntpclient

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
)

func New(pool []string, retries int, timeout int, timeFunc TimeFunc, options ...Option) (*Client, error) {
	if len(pool) == 0 {
		return nil, fmt.Errorf("NTP pool is empty")
	}
	if retries < 0 || timeout < 0 {
		return nil, fmt.Errorf("NTP retries and timeout must not be negative")
	}

	if timeFunc == nil {
		timeFunc = time.Now
	}

	hosts := make([]string, 0, len(pool))
	seen := make(map[string]bool, len(pool))
	servers := make(map[string]*serverState, len(pool))
	for _, host := range pool {
		urlObj, err := url.Parse(host)
		if err != nil {
			return nil, fmt.Errorf("failed to parse NTP endpoint: %w", err)
		}
		if urlObj.Scheme != "ntp" || urlObj.Hostname() == "" {
			return nil, fmt.Errorf("invalid NTP endpoint: %s", host)
		}
		hostname := urlObj.Hostname()
		if address, err := netip.ParseAddr(hostname); err == nil {
			hostname = address.Unmap().String()
		} else {
			hostname = strings.TrimSuffix(strings.ToLower(hostname), ".")
		}
		port := urlObj.Port()
		if port == "" {
			port = "123"
		}
		portNumber, err := strconv.ParseUint(port, 10, 16)
		if err != nil || portNumber == 0 {
			return nil, fmt.Errorf("invalid NTP port in endpoint: %s", host)
		}
		key := net.JoinHostPort(hostname, strconv.FormatUint(portNumber, 10))
		if seen[key] {
			continue
		}
		seen[key] = true
		hosts = append(hosts, urlObj.Host)
		servers[urlObj.Host] = &serverState{interval: minPollInterval}
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ctx:         ctx,
		cancel:      cancel,
		gate:        make(chan struct{}, 1),
		timeFunc:    timeFunc,
		pool:        hosts,
		retries:     retries,
		readTimeout: time.Duration(timeout) * time.Second,
		query:       ntp.QueryWithOptions,
		servers:     servers,
	}
	for _, option := range options {
		option(client)
	}
	return client, nil
}
