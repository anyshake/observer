package ntpclient

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/beevik/ntp"
)

func testResponse(offset, rtt time.Duration) *ntp.Response {
	now := time.Unix(1700000000, 0)
	return &ntp.Response{
		Stratum:       2,
		Time:          now,
		ReferenceTime: now.Add(-time.Minute),
		ClockOffset:   offset,
		RTT:           rtt,
		Precision:     time.Microsecond,
	}
}

type queryRecorder struct {
	mu      sync.Mutex
	calls   map[string]int
	respond func(string, int) (*ntp.Response, error)
}

func (r *queryRecorder) query(server string, _ ntp.QueryOptions) (*ntp.Response, error) {
	r.mu.Lock()
	r.calls[server]++
	attempt := r.calls[server]
	r.mu.Unlock()
	return r.respond(server, attempt)
}

func (r *queryRecorder) count(server string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls[server]
}

func newTestClient(t *testing.T, count, retries int, options ...Option) (*Client, *queryRecorder) {
	t.Helper()
	pool := make([]string, count)
	ranks := make(map[string]int, count)
	for i := range count {
		host := fmt.Sprintf("server%d.example.test", i)
		pool[i] = "ntp://" + host
		ranks[host] = i
	}
	client, err := New(pool, retries, 1, nil, options...)
	if err != nil {
		t.Fatal(err)
	}
	recorder := &queryRecorder{
		calls: make(map[string]int),
		respond: func(server string, _ int) (*ntp.Response, error) {
			return testResponse(1234567*time.Nanosecond, time.Duration(ranks[server]+1)*2*time.Millisecond), nil
		},
	}
	client.query = recorder.query
	t.Cleanup(func() { _ = client.Close() })
	return client, recorder
}

func TestQueryDiscoveryAndPreferredSources(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 8, 5)
		offset, err := client.QueryAverage(5)
		if err != nil || offset != 1234567*time.Nanosecond {
			t.Fatalf("QueryAverage() = %v, %v", offset, err)
		}
		for _, server := range client.pool {
			if got := recorder.count(server); got != 1 {
				t.Errorf("discovery requests to %s = %d, want 1", server, got)
			}
		}

		time.Sleep(2 * time.Minute)
		_, server, err := client.Query()
		if err != nil || server != client.pool[0] {
			t.Fatalf("Query() server = %q, error = %v", server, err)
		}
		for i, server := range client.pool {
			want := 1
			if i < QUERY_ATTEMPT {
				want++
			}
			if got := recorder.count(server); got != want {
				t.Errorf("routine requests to %s = %d, want %d", server, got, want)
			}
		}

		time.Sleep(poolRefreshInterval - 2*time.Minute)
		if _, _, err := client.Query(); err != nil {
			t.Fatal(err)
		}
		for i, server := range client.pool {
			want := 2
			if i < QUERY_ATTEMPT {
				want++
			}
			if got := recorder.count(server); got != want {
				t.Errorf("rediscovery requests to %s = %d, want %d", server, got, want)
			}
		}
	})
}

func TestQueryAverageUsesDistinctSources(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 6, 5)
		client.lastDiscovery = time.Now()
		if _, err := client.QueryAverage(2); err != nil {
			t.Fatal(err)
		}
		for i, server := range client.pool {
			want := 0
			if i < 2 {
				want = 1
			}
			if got := recorder.count(server); got != want {
				t.Errorf("requests to %s = %d, want %d", server, got, want)
			}
		}
	})
}

func TestQueryAverageUsesAvailableFreshSamplesWithoutRounding(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 3, 5)
		base := 30*365*24*time.Hour + 123*time.Nanosecond
		recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
			if server != client.pool[0] {
				return nil, errors.New("timeout")
			}
			return testResponse(base+time.Duration(attempt)*time.Nanosecond, time.Millisecond), nil
		}
		for attempt := 1; attempt <= 2; attempt++ {
			if offset, err := client.QueryAverage(5); err != nil || offset != base+time.Duration(attempt)*time.Nanosecond {
				t.Fatalf("query %d = %v, %v: stale or rounded offset", attempt, offset, err)
			}
			time.Sleep(2 * time.Minute)
		}
	})
}

func TestConcurrentQueriesRespectCooldown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 1, 5)
		var wg sync.WaitGroup
		results := make(chan error, 10)
		for range cap(results) {
			wg.Go(func() {
				_, _, err := client.Query()
				results <- err
			})
		}
		wg.Wait()
		close(results)
		successes := 0
		for err := range results {
			if err == nil {
				successes++
			}
		}
		if successes != 1 || recorder.count(client.pool[0]) != 1 {
			t.Fatalf("concurrent queries: successes=%d, requests=%d, want 1 each", successes, recorder.count(client.pool[0]))
		}
		time.Sleep(minPollInterval - time.Nanosecond)
		if _, _, err := client.Query(); err == nil {
			t.Fatal("query before minimum poll interval succeeded")
		}
		if got := recorder.count(client.pool[0]); got != 1 {
			t.Fatalf("cooldown generated %d requests, want 1", got)
		}
		time.Sleep(time.Nanosecond)
		if _, _, err := client.Query(); err != nil {
			t.Fatalf("query at cooldown boundary failed: %v", err)
		}
		if got := recorder.count(client.pool[0]); got != 2 {
			t.Fatalf("requests after cooldown = %d, want 2", got)
		}
	})
}

func TestQueryFailureUsesOtherServers(t *testing.T) {
	for _, retries := range []int{0, 1} {
		t.Run(fmt.Sprintf("retries=%d", retries), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				client, recorder := newTestClient(t, 8, retries)
				client.lastDiscovery = time.Now()
				respond := recorder.respond
				recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
					if server == client.pool[0] || server == client.pool[1] {
						return nil, errors.New("timeout")
					}
					return respond(server, attempt)
				}
				if _, _, err := client.Query(); err != nil {
					t.Fatal(err)
				}
				for i, server := range client.pool {
					want := 0
					if i < 5 || retries > 0 && i < 7 {
						want = 1
					}
					if got := recorder.count(server); got != want {
						t.Errorf("requests to %s = %d, want %d", server, got, want)
					}
				}
			})
		})
	}
}

func TestQueryRejectsUnusableResponses(t *testing.T) {
	for _, kind := range []string{"network error", "nil response", "invalid stratum", "unsynchronized", "stale", "excess dispersion"} {
		t.Run(kind, func(t *testing.T) {
			client, recorder := newTestClient(t, 3, 5)
			recorder.respond = func(string, int) (*ntp.Response, error) {
				response := testResponse(0, time.Millisecond)
				switch kind {
				case "network error":
					return nil, errors.New("timeout")
				case "nil response":
					return nil, nil
				case "invalid stratum":
					response.Stratum = 16
				case "unsynchronized":
					response.Leap = ntp.LeapNotInSync
				case "stale":
					response.ReferenceTime = response.Time.Add(-48 * time.Hour)
				case "excess dispersion":
					response.RootDispersion = 17 * time.Second
				}
				return response, nil
			}
			if _, _, err := client.Query(); err == nil {
				t.Fatal("query accepted unusable responses")
			}
			for _, server := range client.pool {
				if got := recorder.count(server); got != 1 {
					t.Errorf("requests to failing server %s = %d, want 1", server, got)
				}
			}
		})
	}
}

func TestQueryFailureBackoff(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 1, 5)
		respond := recorder.respond
		recorder.respond = func(string, int) (*ntp.Response, error) { return nil, errors.New("timeout") }
		for i, delay := range []time.Duration{64, 128, 256, 512, 1024, 2048, 3600, 3600} {
			if _, _, err := client.Query(); err == nil {
				t.Fatal("failed source unexpectedly succeeded")
			}
			if got := recorder.count(client.pool[0]); got != i+1 {
				t.Fatalf("failure round %d sent %d requests", i, got)
			}
			time.Sleep(delay*time.Second - time.Nanosecond)
			if _, _, err := client.Query(); err == nil {
				t.Fatal("query during backoff succeeded")
			}
			if got := recorder.count(client.pool[0]); got != i+1 {
				t.Fatalf("backoff round %d sent an early request", i)
			}
			time.Sleep(time.Nanosecond)
		}
		recorder.respond = respond
		if _, _, err := client.Query(); err != nil {
			t.Fatalf("recovered source failed: %v", err)
		}
		time.Sleep(minPollInterval)
		if _, _, err := client.Query(); err != nil {
			t.Fatalf("successful query did not reset failure backoff: %v", err)
		}
	})
}

func TestQueryKissOfDeath(t *testing.T) {
	for _, code := range []string{"RATE", "DENY", "RSTR"} {
		t.Run(code, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				client, recorder := newTestClient(t, 1, 5)
				respond := recorder.respond
				recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
					if attempt == 1 {
						return &ntp.Response{KissCode: code, Poll: 10 * time.Minute}, ntp.ErrKissOfDeath
					}
					return respond(server, attempt)
				}
				if _, _, err := client.Query(); !errors.Is(err, ntp.ErrKissOfDeath) {
					t.Fatalf("KoD query error = %v", err)
				}
				time.Sleep(10*time.Minute - time.Nanosecond)
				if _, _, err := client.Query(); err == nil || recorder.count(client.pool[0]) != 1 {
					t.Fatalf("KoD cooldown ignored: requests=%d, error=%v", recorder.count(client.pool[0]), err)
				}
				time.Sleep(time.Nanosecond)
				if code == "RATE" {
					if _, _, err := client.Query(); err != nil {
						t.Fatalf("query after RATE cooldown failed: %v", err)
					}
					time.Sleep(2 * time.Minute)
					if _, _, err := client.Query(); err == nil || recorder.count(client.pool[0]) != 2 {
						t.Fatal("successful reply reset the RATE polling interval")
					}
				} else {
					time.Sleep(24 * time.Hour)
					if _, _, err := client.Query(); err == nil || recorder.count(client.pool[0]) != 1 {
						t.Fatalf("disabled source was queried again: error=%v", err)
					}
				}
			})
		})
	}
}

func TestRateResponseIncreasesIntervalWithoutPollHint(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 1, 0)
		recorder.respond = func(string, int) (*ntp.Response, error) {
			return &ntp.Response{KissCode: "RATE"}, nil
		}
		for _, delay := range []time.Duration{128 * time.Second, 256 * time.Second} {
			if _, _, err := client.Query(); !errors.Is(err, ntp.ErrKissOfDeath) {
				t.Fatalf("KoD query error = %v", err)
			}
			if got := client.servers[client.pool[0]].interval; got != delay {
				t.Fatalf("RATE interval = %v, want %v", got, delay)
			}
			time.Sleep(delay)
		}
	})
}

func TestQueryRanksBySynchronizationDistance(t *testing.T) {
	client, recorder := newTestClient(t, 3, 0)
	respond := recorder.respond
	recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
		response, err := respond(server, attempt)
		if server == client.pool[0] {
			response.RootDispersion = time.Second
		}
		return response, err
	}
	if _, server, err := client.Query(); err != nil || server != client.pool[1] {
		t.Fatalf("best source = %q, error = %v; root dispersion was ignored", server, err)
	}
}

func TestQueryDemotesDisagreeingSource(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 6, 1)
		respond := recorder.respond
		recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
			if server == client.pool[0] {
				return testResponse(time.Hour, time.Microsecond), nil
			}
			return respond(server, attempt)
		}
		if offset, server, err := client.Query(); err != nil || offset != 1234567*time.Nanosecond || server != client.pool[1] {
			t.Fatalf("outlier affected synchronization: offset=%v, server=%q, error=%v", offset, server, err)
		}
		time.Sleep(2 * time.Minute)
		if _, _, err := client.Query(); err != nil {
			t.Fatal(err)
		}
		if got := recorder.count(client.pool[0]); got != 1 {
			t.Fatalf("disagreeing source remained preferred: %d requests", got)
		}
	})
}

func TestQueryDisagreementUsesOtherSources(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 8, 1)
		client.lastDiscovery = time.Now()
		recorder.respond = func(server string, _ int) (*ntp.Response, error) {
			offset := time.Duration(0)
			if server == client.pool[0] || server == client.pool[1] {
				offset = time.Second
			}
			if server == client.pool[2] {
				offset = -time.Second
			}
			return testResponse(offset, time.Millisecond), nil
		}
		if offset, _, err := client.Query(); err != nil || offset != 0 {
			t.Fatalf("fallback consensus = %v, %v", offset, err)
		}
		for _, server := range client.pool {
			if got := recorder.count(server); got != 1 {
				t.Errorf("requests to %s = %d, want 1", server, got)
			}
		}
	})
}

func TestQueryRejectsDisagreement(t *testing.T) {
	client, recorder := newTestClient(t, 2, 5)
	recorder.respond = func(server string, _ int) (*ntp.Response, error) {
		offset := time.Duration(0)
		if server == client.pool[0] {
			offset = time.Second
		}
		return testResponse(offset, time.Millisecond), nil
	}
	if _, _, err := client.Query(); err == nil {
		t.Fatal("disagreeing sources produced a clock correction")
	}
}
