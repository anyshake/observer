package dnsquery

import (
	"fmt"
	"net/url"
)

func New(server string) (IServer, error) {
	urlObj, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL: %w", err)
	}

	switch urlObj.Scheme {
	case "https":
		return &DoH{server: server}, nil
	case "tls":
		return &DoT{server: fixPortNumber(urlObj, "853").Host}, nil
	case "udp":
		return &UDP{server: fixPortNumber(urlObj, "53").Host}, nil
	case "sdns":
		return &DNSCrypt{server: server}, nil
	}

	return nil, fmt.Errorf("unsupported server %s", server)
}
