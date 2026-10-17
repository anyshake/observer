package dnsquery

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func New(server string) (IServer, error) {
	urlObj, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL: %w", err)
	}
	if urlObj.Scheme != "sdns" {
		if urlObj.Hostname() == "" || urlObj.User != nil || urlObj.Fragment != "" {
			return nil, fmt.Errorf("invalid DNS endpoint: %s", server)
		}
		if port := urlObj.Port(); port != "" {
			number, err := strconv.ParseUint(port, 10, 16)
			if err != nil || number == 0 {
				return nil, fmt.Errorf("invalid DNS port in endpoint: %s", server)
			}
		}
		if urlObj.Scheme != "https" && (urlObj.Path != "" || urlObj.RawQuery != "") {
			return nil, fmt.Errorf("DNS endpoint does not support a path or query: %s", server)
		}
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

// NewClient keeps resolver quality across queries. Construction performs no I/O.
// retries limits additional batches; an endpoint is queried at most once per call.
func NewClient(pool Resolvers, retries int, timeout time.Duration, options ...Option) (*Client, error) {
	if len(pool) == 0 {
		return nil, fmt.Errorf("DNS pool is empty")
	}
	if retries < 0 || timeout <= 0 {
		return nil, fmt.Errorf("DNS retries must not be negative and timeout must be positive")
	}
	servers := make(map[string]*resolverState, len(pool))
	resolvers := make(Resolvers, 0, len(pool))
	for _, resolver := range pool {
		server, err := New(resolver.Server)
		if err != nil {
			return nil, err
		}
		endpoint, _ := url.Parse(resolver.Server)
		if endpoint.Scheme != "sdns" {
			host := strings.TrimSuffix(strings.ToLower(endpoint.Hostname()), ".")
			if ip, err := netip.ParseAddr(host); err == nil {
				host = ip.Unmap().String()
			}
			port := endpoint.Port()
			if port == "" {
				switch endpoint.Scheme {
				case "https":
					port = "443"
				case "tls":
					port = "853"
				case "udp":
					port = "53"
				}
			}
			number, _ := strconv.ParseUint(port, 10, 16)
			endpoint.Host = net.JoinHostPort(host, strconv.FormatUint(number, 10))
			if endpoint.Scheme == "https" && endpoint.Path == "" {
				endpoint.Path = "/"
			}
			resolver.Server = endpoint.String()
		}
		if _, seen := servers[resolver.Server]; seen {
			continue
		}
		servers[resolver.Server] = &resolverState{server: server}
		resolvers = append(resolvers, resolver)
	}
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ctx: ctx, cancel: cancel, gate: make(chan struct{}, 1),
		pool: resolvers, servers: servers, retries: retries, timeout: timeout,
	}
	for _, option := range options {
		option(client)
	}
	return client, nil
}
