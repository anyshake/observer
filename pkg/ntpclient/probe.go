package ntpclient

import (
	"context"
	"sync"
	"time"

	"github.com/beevik/ntp"
)

// Probe queries every supplied endpoint once without client polling state.
// Use Client.Query for rate-limited synchronization.
func Probe(servers []string, timeout time.Duration, timeFn func() time.Time) []ProbeResult {
	return probe(servers, timeout, timeFn, ntp.QueryWithOptions)
}

func probe(servers []string, timeout time.Duration, timeFn func() time.Time, query queryFunc) []ProbeResult {
	return probeContext(context.Background(), servers, timeout, timeFn, query, nil)
}

func probeContext(ctx context.Context, servers []string, timeout time.Duration, timeFn func() time.Time, query queryFunc, onResult func(ProbeResult)) []ProbeResult {
	results := make([]ProbeResult, len(servers))
	sem := make(chan struct{}, CONCURRENT_QUERIES)
	var wg sync.WaitGroup
	var progressMu sync.Mutex

	for i, server := range servers {
		if ctx.Err() != nil {
			break
		}
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Wait()
			return results
		}
		wg.Add(1)
		go func(i int, server string) {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := queryWithContext(ctx, server, ntp.QueryOptions{
				Timeout:       timeout,
				GetSystemTime: timeFn,
			}, query)
			results[i] = ProbeResult{Server: server, Resp: resp, Err: err}
			if onResult != nil {
				progressMu.Lock()
				onResult(results[i])
				progressMu.Unlock()
			}
		}(i, server)
	}

	wg.Wait()
	return results
}
