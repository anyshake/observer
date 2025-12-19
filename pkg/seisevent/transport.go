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
func createCustomTransport(dnsList dnsquery.Resolvers, frontendSni string) *http.Transport {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			hostname, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse address: %w", err)
			}

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			type result struct {
				ip  net.IP
				err error
			}
			resultCh := make(chan result, len(dnsList))

			for _, dnsServer := range dnsList {
				go func() {
					dnsResolver, err := dnsquery.New(dnsServer.Server)
					if err != nil {
						resultCh <- result{nil, fmt.Errorf("failed to create DNS resolver: %w", err)}
						return
					}
					if err := dnsResolver.Open(); err != nil {
						resultCh <- result{nil, fmt.Errorf("failed to open DNS resolver: %w", err)}
						return
					}
					defer dnsResolver.Close()

					res, err := dnsResolver.Query((&dns.Msg{}).SetQuestion(fmt.Sprintf("%s.", hostname), dns.TypeA), 5*time.Second)
					if err != nil {
						resultCh <- result{nil, fmt.Errorf("failed to query DNS: %w", err)}
						return
					}
					if len(res.Answer) == 0 {
						resultCh <- result{nil, errors.New("no answer from DNS")}
						return
					}

					for _, ans := range res.Answer {
						if aRecord, ok := ans.(*dns.A); ok && len(aRecord.A) > 0 {
							select {
							case resultCh <- result{aRecord.A, nil}:
							case <-ctx.Done():
							}
							return
						}
					}
					resultCh <- result{nil, errors.New("no valid A record")}
				}()
			}

			for range dnsList {
				select {
				case r := <-resultCh:
					if r.err == nil {
						return (&net.Dialer{}).DialContext(ctx, network, fmt.Sprintf("%s:%s", r.ip.String(), port))
					}
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}

			return nil, errors.New("all DNS queries failed")
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
