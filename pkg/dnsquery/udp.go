package dnsquery

import (
	"errors"
	"time"

	"github.com/miekg/dns"
)

type UDP struct {
	server string
	client *dns.Client
}

func (d *UDP) Open() error {
	d.client = &dns.Client{}
	return nil
}

func (t *UDP) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if t.client == nil {
		return nil, errors.New("UDP client is not initialized")
	}

	t.client.Timeout = timeout
	res, _, err := t.client.Exchange(msg, t.server)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *UDP) Close() error {
	t.client = nil
	return nil
}
