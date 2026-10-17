package dnsquery

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ameshkov/dnscrypt/v2"
	"github.com/miekg/dns"
)

type DNSCrypt struct {
	server   string
	client   *dnscrypt.Client
	resolver *dnscrypt.ResolverInfo
}

func (d *DNSCrypt) Open() error {
	return d.OpenContext(context.Background(), 5*time.Second)
}

func (d *DNSCrypt) OpenContext(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := queryContext(ctx, timeout)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	if d.resolver != nil && d.resolver.ResolverCert.VerifyDate() {
		return nil
	}
	resolver, err := dialDNSCrypt(ctx, d.server, timeout)
	if err != nil {
		return fmt.Errorf("failed to dial DNSCrypt server: %w", err)
	}
	d.client = &dnscrypt.Client{Net: "tcp"}
	d.resolver = resolver
	return nil
}

func (d *DNSCrypt) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	return d.QueryContext(context.Background(), msg, timeout)
}

func (d *DNSCrypt) QueryContext(ctx context.Context, msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if d.client == nil || d.resolver == nil {
		return nil, errors.New("DNSCrypt client is not initialized")
	}

	ctx, cancel := queryContext(ctx, timeout)
	defer cancel()
	if err := d.OpenContext(ctx, timeout); err != nil {
		return nil, err
	}
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", d.resolver.ServerAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	stop := closeOnCancel(ctx, func() { _ = conn.Close() })
	defer stop()
	client := *d.client
	client.Timeout = timeout
	response, err := client.ExchangeConn(conn, msg, d.resolver)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return response, err
}

func (d *DNSCrypt) Close() error {
	d.resolver = nil
	d.client = nil
	return nil
}
