package dnsquery

import (
	"errors"
	"time"

	"github.com/miekg/dns"
)

type DoT struct {
	server string
	client *dns.Client
}

func (d *DoT) Open() error {
	d.client = &dns.Client{Net: "tcp-tls"}
	return nil
}

func (t *DoT) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if t.client == nil {
		return nil, errors.New("DoT client is not initialized")
	}

	t.client.Timeout = timeout
	res, _, err := t.client.Exchange(msg, t.server)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *DoT) Close() error {
	t.client = nil
	return nil
}
