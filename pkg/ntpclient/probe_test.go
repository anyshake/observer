package ntpclient

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/beevik/ntp"
)

func TestProbeConcurrencyAndOptions(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		servers := make([]string, 2*CONCURRENT_QUERIES+1)
		for i := range servers {
			servers[i] = fmt.Sprintf("server%d.example.test", i)
		}
		clock := time.Unix(1000, 123)
		timeout := 3 * time.Second
		queryErr := errors.New("query failed")
		started := make(chan string, len(servers))
		release := make(chan struct{})
		done := make(chan []ProbeResult, 1)
		var active atomic.Int32
		go func() {
			done <- probe(servers, timeout, func() time.Time { return clock }, func(server string, options ntp.QueryOptions) (*ntp.Response, error) {
				if got := active.Add(1); got > CONCURRENT_QUERIES {
					t.Errorf("active queries = %d, exceeds %d", got, CONCURRENT_QUERIES)
				}
				defer active.Add(-1)
				if options.Timeout != timeout || options.GetSystemTime == nil || !options.GetSystemTime().Equal(clock) {
					t.Errorf("query options were not forwarded: %+v", options)
				}
				started <- server
				<-release
				if server == servers[1] {
					return nil, queryErr
				}
				return testResponse(0, time.Millisecond), nil
			})
		}()
		synctest.Wait()
		initial := len(started)
		close(release)
		results := <-done
		if initial != CONCURRENT_QUERIES {
			t.Errorf("initial concurrent requests = %d, want %d", initial, CONCURRENT_QUERIES)
		}
		if len(results) != len(servers) || len(started) != len(servers) {
			t.Fatalf("results=%d, requests=%d, want %d each", len(results), len(started), len(servers))
		}
		for i, result := range results {
			if result.Server != servers[i] {
				t.Errorf("result %d server = %q, want %q", i, result.Server, servers[i])
			}
			if i == 1 {
				if !errors.Is(result.Err, queryErr) || result.Resp != nil {
					t.Errorf("failed query result = %+v", result)
				}
			} else if result.Err != nil || result.Resp == nil {
				t.Errorf("successful query result = %+v", result)
			}
		}
	})
}

func TestProbeEmptyPool(t *testing.T) {
	if got := Probe(nil, time.Second, time.Now); len(got) != 0 {
		t.Errorf("Probe(nil) = %v, want no results", got)
	}
}
