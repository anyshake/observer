package ntpclient

import (
	"sync"
	"time"

	"github.com/beevik/ntp"
)

func Probe(servers []string, timeout time.Duration, timeFn func() time.Time) []ProbeResult {
	results := make([]ProbeResult, len(servers))
	sem := make(chan struct{}, CONCURRENT_QUERIES)
	var wg sync.WaitGroup

	for i, server := range servers {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, server string) {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := ntp.QueryWithOptions(server, ntp.QueryOptions{
				Timeout:       timeout,
				GetSystemTime: timeFn,
			})
			results[i] = ProbeResult{Server: server, Resp: resp, Err: err}
		}(i, server)
	}

	wg.Wait()
	return results
}
