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

// Magic function that bypasses the Great Firewall of China
func createCustomResolver(dnsList []string, frontendSni string) *http.Transport {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			hostname, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse address: %w", err)
			}

			type result struct {
				ip  net.IP
				err error
			}

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			resultCh := make(chan result, len(dnsList))

			for _, dnsServer := range dnsList {
				dnsServer := dnsServer
				go func() {
					dnsResolver, err := dnsquery.New(dnsServer)
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

			for i := 0; i < len(dnsList); i++ {
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

	if len(frontendSni) > 0 {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         frontendSni,
			MinVersion:         tls.VersionTLS12,
		}
	}

	return transport
}

// Most overseas DoH / DoT providers are blocked in China
// Recommended DNSCrypt: https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/master/v3/public-resolvers.md
func getCustomDnsList() []string {
	return []string{
		// cs-tokyo
		"sdns://AQcAAAAAAAAADDE0Ni43MC4zMS40MyAxM3KtWVYywkFrhy8Jj4Ub3bllKExsvppPGQlkMNupWh4yLmRuc2NyeXB0LWNlcnQuY3J5cHRvc3Rvcm0uaXM",
		// cs-austria
		"sdns://AQcAAAAAAAAADTk0LjE5OC40MS4yMzUgMTNyrVlWMsJBa4cvCY-FG925ZShMbL6aTxkJZDDbqVoeMi5kbnNjcnlwdC1jZXJ0LmNyeXB0b3N0b3JtLmlz",
		// cs-london
		"sdns://AQcAAAAAAAAADTc4LjEyOS4yNDguNjcgMTNyrVlWMsJBa4cvCY-FG925ZShMbL6aTxkJZDDbqVoeMi5kbnNjcnlwdC1jZXJ0LmNyeXB0b3N0b3JtLmlz",
		// cs-la
		"sdns://AQcAAAAAAAAADzE5NS4yMDYuMTA0LjIwMyAxM3KtWVYywkFrhy8Jj4Ub3bllKExsvppPGQlkMNupWh4yLmRuc2NyeXB0LWNlcnQuY3J5cHRvc3Rvcm0uaXM",
		// cs-montreal
		"sdns://AQcAAAAAAAAADTE3Ni4xMTMuNzQuMTkgMTNyrVlWMsJBa4cvCY-FG925ZShMbL6aTxkJZDDbqVoeMi5kbnNjcnlwdC1jZXJ0LmNyeXB0b3N0b3JtLmlz",
		// cs-nyc
		"sdns://AQcAAAAAAAAADTE0Ni43MC4xNTQuNjcgMTNyrVlWMsJBa4cvCY-FG925ZShMbL6aTxkJZDDbqVoeMi5kbnNjcnlwdC1jZXJ0LmNyeXB0b3N0b3JtLmlz",
		// cs-ore
		"sdns://AQcAAAAAAAAADTE3OS42MS4yMjMuNDcgMTNyrVlWMsJBa4cvCY-FG925ZShMbL6aTxkJZDDbqVoeMi5kbnNjcnlwdC1jZXJ0LmNyeXB0b3N0b3JtLmlz",
		// cs-singapore
		"sdns://AQcAAAAAAAAADTM3LjEyMC4xNTEuMTEgMTNyrVlWMsJBa4cvCY-FG925ZShMbL6aTxkJZDDbqVoeMi5kbnNjcnlwdC1jZXJ0LmNyeXB0b3N0b3JtLmlz",
		// cs-sk
		"sdns://AQcAAAAAAAAADjEwOC4xODEuNTAuMjE4IDEzcq1ZVjLCQWuHLwmPhRvduWUoTGy-mk8ZCWQw26laHjIuZG5zY3J5cHQtY2VydC5jcnlwdG9zdG9ybS5pcw",
		// dnscry.pt-tokyo-ipv4
		"sdns://AQcAAAAAAAAADDQ1LjY3Ljg2LjEyMyBDK5aRHZnKfdd6Q9ufEJY83WAQ9X5z7OAQa5CeptBCYBkyLmRuc2NyeXB0LWNlcnQuZG5zY3J5LnB0",
		// dnscry.pt-tokyo02-ipv4
		"sdns://AQcAAAAAAAAADDEwMy4xNzkuNDUuNiDfai5sp1im-BPHwbM1GCnTqn20FIbQfuvvybKsGf0pjhkyLmRuc2NyeXB0LWNlcnQuZG5zY3J5LnB0",
		// jp.tiar.app
		"sdns://AQcAAAAAAAAAEjE3Mi4xMDQuOTMuODA6MTQ0MyAyuHY-8b9lNqHeahPAzW9IoXnjiLaZpTeNbVs8TN9UUxsyLmRuc2NyeXB0LWNlcnQuanAudGlhci5hcHA",
		// dnscry.pt-seoul-ipv4
		"sdns://AQcAAAAAAAAADTkyLjM4LjEzNS4xMjggyHfVGamJyxLfoAWjERmO4pY3KzKkqY-vSa2UnVx_gYAZMi5kbnNjcnlwdC1jZXJ0LmRuc2NyeS5wdA",
		// digitalprivacy.diy-dnscrypt-ipv4
		"sdns://AQcAAAAAAAAAEjM3LjIyMS4xOTQuODQ6NDQzNCCiyGRvm0TcyJmI7lTXstgh-8AoAAiFcTQQp7Od_brTYCIyLmRuc2NyeXB0LWNlcnQuZGlnaXRhbHByaXZhY3kuZGl5",
	}
}
