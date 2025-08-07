package ntpclient

import (
	"fmt"
	"sync"

	"github.com/beevik/ntp"
)

func (c *Client) Query() (ntp.Response, error) {
	const requestCount = 10

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		bestResp *ntp.Response
	)

	wg.Add(requestCount)
	for i := 0; i < requestCount; i++ {
		go func() {
			defer wg.Done()

			resp, err := ntp.QueryWithOptions(c.ntpAddr, ntp.QueryOptions{Timeout: c.readTimeout})
			if err != nil {
				return
			}

			if resp.Stratum == 0 || resp.RootDistance <= 0 {
				return
			}

			mu.Lock()
			defer mu.Unlock()

			if bestResp == nil || resp.RootDistance < bestResp.RootDistance || (resp.RootDistance == bestResp.RootDistance && resp.RTT < bestResp.RTT) {
				bestResp = resp
			}
		}()
	}

	wg.Wait()

	if bestResp == nil {
		return ntp.Response{}, fmt.Errorf("all NTP queries failed or returned invalid data")
	}

	return *bestResp, nil
}
