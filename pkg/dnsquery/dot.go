package dnsquery

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/miekg/dns"
)

type DoT struct {
	server string
	client *dns.Client
}

func (d *DoT) Open() error {
	return d.OpenContext(context.Background(), 5*time.Second)
}

func (d *DoT) OpenContext(ctx context.Context, _ time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d.client == nil {
		d.client = &dns.Client{Net: "tcp-tls", TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	}
	return nil
}

func (t *DoT) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	return t.QueryContext(context.Background(), msg, timeout)
}

func (t *DoT) QueryContext(ctx context.Context, msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if t.client == nil {
		return nil, errors.New("DoT client is not initialized")
	}

	ctx, cancel := queryContext(ctx, timeout)
	defer cancel()
	client := *t.client
	client.Timeout = timeout
	return exchangeContext(ctx, &client, msg, t.server)
}

func (t *DoT) Close() error {
	t.client = nil
	return nil
}
