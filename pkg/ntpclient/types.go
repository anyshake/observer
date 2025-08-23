package ntpclient

import (
	"time"

	"github.com/beevik/ntp"
)

const (
	QUERY_ATTEMPT      = 30
	CONCURRENT_QUERIES = 5
)

type Client struct {
	pool        []string
	retries     int
	readTimeout time.Duration
}

type probeResult struct {
	resp   *ntp.Response
	server string
	err    error
}
