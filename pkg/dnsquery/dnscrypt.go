package dnsquery

import (
	"errors"
	"fmt"
	"time"

	"github.com/ameshkov/dnscrypt/v2"
	"github.com/jedisct1/go-dnsstamps"
	"github.com/miekg/dns"
)

type DNSCrypt struct {
	server   string
	client   *dnscrypt.Client
	resolver *dnscrypt.ResolverInfo
}

func (d *DNSCrypt) Open() error {
	d.client = &dnscrypt.Client{Net: "tcp"}

	stamp, err := dnsstamps.NewServerStampFromString(d.server)
	if err != nil {
		return fmt.Errorf("failed to create stamp from provider information: %w", err)
	}

	ro, err := d.client.Dial(stamp.String())
	if err != nil {
		return fmt.Errorf("failed to dial DNSCrypt server: %w", err)
	}

	d.resolver = ro
	return nil
}

func (d *DNSCrypt) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if d.client == nil {
		return nil, errors.New("DNSCrypt client is not initialized")
	}

	d.client.Timeout = timeout
	return d.client.Exchange(msg, d.resolver)
}

func (d *DNSCrypt) Close() error {
	d.resolver = nil
	d.client = nil
	return nil
}
