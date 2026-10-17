package dnsquery

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

type IServer interface {
	Open() error
	OpenContext(ctx context.Context, timeout time.Duration) error
	Close() error
	Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error)
	QueryContext(ctx context.Context, msg *dns.Msg, timeout time.Duration) (*dns.Msg, error)
}

const (
	concurrentQueries   = 5
	poolRefreshInterval = 30 * time.Minute
	minFailureBackoff   = time.Second
	maxFailureBackoff   = time.Minute
)

type Client struct {
	ctx           context.Context
	cancel        context.CancelFunc
	gate          chan struct{}
	pool          Resolvers
	servers       map[string]*resolverState
	retries       int
	timeout       time.Duration
	lastDiscovery time.Time
	logger        Logger
}

type resolverState struct {
	server    IServer
	valid     bool
	rtt       time.Duration
	lastProbe time.Time
	nextTry   time.Time
	backoff   time.Duration
}

type probeResult struct {
	server   string
	response *dns.Msg
	err      error
	rtt      time.Duration
	canceled bool
}
