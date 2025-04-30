package dnsquery

import (
	"time"

	"github.com/miekg/dns"
)

type IServer interface {
	Open() error
	Close() error
	Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error)
}
