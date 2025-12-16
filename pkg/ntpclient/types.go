package ntpclient

import (
	"time"

	"github.com/beevik/ntp"
)

const (
	QUERY_ATTEMPT      = 5
	CONCURRENT_QUERIES = 5
)

type TimeFunc func() time.Time

type ProbeResult struct {
	Resp   *ntp.Response
	Server string
	Err    error
}

type Client struct {
	timeFunc    TimeFunc
	pool        []string
	retries     int
	readTimeout time.Duration
}
