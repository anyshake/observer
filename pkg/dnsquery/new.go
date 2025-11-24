package dnsquery

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v2"
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

func NewResolvers() Resolvers {
	fallback := Resolvers{
		{Name: "Google", Server: "udp://8.8.8.8"},
		{Name: "Cloudflare", Server: "udp://1.0.0.1"},
		{Name: "Quad9", Server: "udp://9.9.9.9"},
	}

	data, err := _resolvers.ReadFile("resolvers.yaml")
	if err != nil {
		return fallback
	}

	var resolvers []Resolver
	if err := yaml.Unmarshal(data, &resolvers); err != nil {
		return fallback
	}

	return resolvers
}
