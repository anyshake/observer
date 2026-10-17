package dnsquery

import (
	"context"
	"errors"
	"time"

	"github.com/miekg/dns"
)

type UDP struct {
	server string
	client *dns.Client
}

func (d *UDP) Open() error {
	return d.OpenContext(context.Background(), 5*time.Second)
}

func (d *UDP) OpenContext(ctx context.Context, _ time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d.client == nil {
		d.client = &dns.Client{}
	}
	return nil
}

func (t *UDP) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	return t.QueryContext(context.Background(), msg, timeout)
}

func (t *UDP) QueryContext(ctx context.Context, msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if t.client == nil {
		return nil, errors.New("UDP client is not initialized")
	}

	ctx, cancel := queryContext(ctx, timeout)
	defer cancel()
	client := *t.client
	client.Timeout = timeout
	return exchangeContext(ctx, &client, msg, t.server)
}

func (t *UDP) Close() error {
	t.client = nil
	return nil
}
