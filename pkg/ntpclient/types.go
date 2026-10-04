package ntpclient

import (
	"context"
	"time"

	"github.com/beevik/ntp"
)

const (
	QUERY_ATTEMPT      = 5
	CONCURRENT_QUERIES = 5

	poolRefreshInterval = 30 * time.Minute
	minPollInterval     = 64 * time.Second
	maxFailureBackoff   = time.Hour
)

type TimeFunc func() time.Time

type ProbeResult struct {
	Resp   *ntp.Response
	Server string
	Err    error
}

type Client struct {
	ctx           context.Context
	cancel        context.CancelFunc
	gate          chan struct{}
	timeFunc      TimeFunc
	pool          []string
	retries       int
	readTimeout   time.Duration
	query         queryFunc
	logger        Logger
	servers       map[string]*serverState
	lastDiscovery time.Time
}

type queryFunc func(string, ntp.QueryOptions) (*ntp.Response, error)

type serverState struct {
	nextPoll time.Time
	interval time.Duration
	backoff  time.Duration
	distance time.Duration
	valid    bool
	disabled bool
}

type clockSample struct {
	server   string
	offset   time.Duration
	distance time.Duration
}
