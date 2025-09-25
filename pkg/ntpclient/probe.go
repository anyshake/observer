package ntpclient

import (
	"sync"

	"github.com/bclswl0827/ntp"
)

func (c *Client) parallelProbe(servers []string) []probeResult {
	var wg sync.WaitGroup
	results := make([]probeResult, len(servers))

	for i, server := range servers {
		wg.Add(1)
		go func(i int, server string) {
			defer wg.Done()
			resp, err := ntp.QueryWithOptions(server, ntp.QueryOptions{Timeout: c.readTimeout})
			results[i] = probeResult{server: server, resp: resp, err: err}
		}(i, server)
	}
	wg.Wait()

	return results
}
