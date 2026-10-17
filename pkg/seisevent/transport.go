package seisevent

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/miekg/dns"
)

// Magic function that bypasses the Great Firewall of China using encrypted DNS and SNI-spoofing
func createCustomTransport(client *dnsquery.Client, frontendSni string) *http.Transport {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			hostname, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse address: %w", err)
			}

			if net.ParseIP(hostname) != nil {
				return (&net.Dialer{}).DialContext(ctx, network, addr)
			}
			if client == nil {
				return nil, errors.New("DNS client is not initialized")
			}
			queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			res, err := client.QueryContext(queryCtx, new(dns.Msg).SetQuestion(dns.Fqdn(hostname), dns.TypeA))
			if err != nil {
				return nil, fmt.Errorf("failed to query DNS: %w", err)
			}
			for _, answer := range res.Answer {
				if record, ok := answer.(*dns.A); ok && record.A.To4() != nil {
					return (&net.Dialer{}).DialContext(ctx, network, net.JoinHostPort(record.A.String(), port))
				}
			}
			return nil, errors.New("DNS response has no valid A record")
		},
	}

	if frontendSni != "" {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         frontendSni,
			MinVersion:         tls.VersionTLS12,
		}
	}

	return transport
}
