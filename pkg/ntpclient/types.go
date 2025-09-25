package ntpclient

import (
	"time"

	"github.com/bclswl0827/ntp"
)

const (
	QUERY_ATTEMPT      = 10
	CONCURRENT_QUERIES = 5
)

type Client struct {
	timeFunc    ntp.TimeFunc
	pool        []string
	retries     int
	readTimeout time.Duration
}

type probeResult struct {
	resp   *ntp.Response
	server string
	err    error
}
