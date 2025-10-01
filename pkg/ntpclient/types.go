package ntpclient

import (
	"time"

	"github.com/beevik/ntp"
)

const (
	QUERY_ATTEMPT      = 10
	CONCURRENT_QUERIES = 5
)

type TimeFunc func() time.Time

type Client struct {
	timeFunc    TimeFunc
	pool        []string
	retries     int
	readTimeout time.Duration
}

type probeResult struct {
	resp   *ntp.Response
	server string
	err    error
}
