package dnsquery

import (
	"net"
	"net/url"
)

func fixPortNumber(urlObj *url.URL, defaultPort string) *url.URL {
	if urlObj.Port() == "" {
		copyURL := *urlObj
		copyURL.Host = net.JoinHostPort(urlObj.Hostname(), defaultPort)
		urlObj = &copyURL
	}

	return urlObj
}
