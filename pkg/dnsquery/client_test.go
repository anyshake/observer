package dnsquery

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/miekg/dns"
)

type stubServer struct {
	query func(context.Context, *dns.Msg) (*dns.Msg, error)
}

func (s *stubServer) Open() error                                            { return nil }
func (s *stubServer) OpenContext(ctx context.Context, _ time.Duration) error { return ctx.Err() }
func (s *stubServer) Close() error                                           { return nil }
func (s *stubServer) Query(msg *dns.Msg, _ time.Duration) (*dns.Msg, error) {
	return s.query(context.Background(), msg)
}
func (s *stubServer) QueryContext(ctx context.Context, msg *dns.Msg, _ time.Duration) (*dns.Msg, error) {
	return s.query(ctx, msg)
}

func stubClient(t *testing.T, count, retries int, query func(int, context.Context, *dns.Msg) (*dns.Msg, error)) *Client {
	t.Helper()
	pool := make(Resolvers, count)
	for i := range pool {
		pool[i].Server = fmt.Sprintf("tls://resolver-%d.example.test", i)
	}
	client, err := NewClient(pool, retries, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	for i, resolver := range client.pool {
		client.servers[resolver.Server].server = &stubServer{query: func(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
			return query(i, ctx, msg)
		}}
	}
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func TestClientCancelsLosersAndReusesPreferredResolver(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		started := make(chan struct{}, concurrentQueries-1)
		var calls, canceled atomic.Int32
		client := stubClient(t, 7, 1, func(i int, ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
			if i == 0 {
				if calls.Add(1) == 1 {
					for range concurrentQueries - 1 {
						<-started
					}
				}
				return exampleReply(msg), nil
			}
			if i >= concurrentQueries {
				t.Errorf("queried endpoint %d after a successful batch", i)
			}
			started <- struct{}{}
			<-ctx.Done()
			canceled.Add(1)
			return nil, ctx.Err()
		})
		for range 2 {
			question := exampleQuestion()
			response, err := client.Query(question)
			if err != nil || response.Id != question.Id {
				t.Fatalf("Query = %v, %v", response, err)
			}
		}
		if calls.Load() != 2 || canceled.Load() != concurrentQueries-1 {
			t.Fatalf("preferred requests=%d joined losers=%d", calls.Load(), canceled.Load())
		}
		for _, state := range client.servers {
			if state.backoff != 0 {
				t.Fatal("canceling a losing probe penalized its resolver")
			}
		}
	})
}

func TestClientBoundsConcurrencyAndDoesNotRepeatEndpoints(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var active, peak atomic.Int32
		calls := make([]atomic.Int32, 13)
		client := stubClient(t, len(calls), 1, func(i int, _ context.Context, _ *dns.Msg) (*dns.Msg, error) {
			calls[i].Add(1)
			n := active.Add(1)
			for previous := peak.Load(); n > previous; previous = peak.Load() {
				if peak.CompareAndSwap(previous, n) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			active.Add(-1)
			return nil, errors.New("unavailable")
		})
		if _, err := client.Query(exampleQuestion()); err == nil {
			t.Fatal("all failed probes returned success")
		}
		total := int32(0)
		for i := range calls {
			if count := calls[i].Load(); count > 1 {
				t.Fatalf("endpoint %d queried %d times", i, count)
			}
			total += calls[i].Load()
		}
		if total != 2*concurrentQueries || peak.Load() > concurrentQueries || active.Load() != 0 {
			t.Fatalf("requests=%d peak=%d active=%d", total, peak.Load(), active.Load())
		}
	})
}

func TestClientPeriodicDiscoveryIncludesUntriedResolvers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		started := make(chan int, concurrentQueries)
		release := make(chan struct{})
		client := stubClient(t, 6, 0, func(i int, _ context.Context, msg *dns.Msg) (*dns.Msg, error) {
			started <- i
			<-release
			return exampleReply(msg), nil
		})
		for round := range 2 {
			result := make(chan error, 1)
			go func() {
				_, err := client.Query(exampleQuestion())
				result <- err
			}()
			includedLast := false
			for range concurrentQueries {
				includedLast = (<-started == 5) || includedLast
			}
			close(release)
			if err := <-result; err != nil {
				t.Fatal(err)
			}
			if includedLast != (round == 1) {
				t.Fatalf("discovery round %d included untried resolver = %t", round, includedLast)
			}
			release = make(chan struct{})
			time.Sleep(poolRefreshInterval)
		}
	})
}

func TestClientTimedOutDiscoveryTriesUnprobedResolvers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client := stubClient(t, 2*concurrentQueries+1, 2, func(i int, ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
			if i == 2*concurrentQueries {
				return exampleReply(msg), nil
			}
			<-ctx.Done()
			return nil, ctx.Err()
		})
		for range 2 {
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			_, err := client.QueryContext(ctx, exampleQuestion())
			cancel()
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("blocked discovery = %v", err)
			}
		}
		for _, state := range client.servers {
			if state.backoff != 0 || !state.nextTry.IsZero() {
				t.Fatal("caller deadline penalized a resolver")
			}
		}
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		if _, err := client.QueryContext(ctx, exampleQuestion()); err != nil {
			t.Fatalf("unprobed reachable resolver was not tried: %v", err)
		}
	})
}

func TestClientFailureBackoffAndRecovery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		client := stubClient(t, 1, 0, func(_ int, _ context.Context, msg *dns.Msg) (*dns.Msg, error) {
			calls++
			if calls == 1 {
				return nil, errors.New("unavailable")
			}
			return exampleReply(msg), nil
		})
		for range 2 {
			if _, err := client.Query(exampleQuestion()); err == nil {
				t.Fatal("failed or cooling resolver returned success")
			}
		}
		if calls != 1 {
			t.Fatalf("queries while cooling down = %d", calls)
		}
		time.Sleep(minFailureBackoff)
		if _, err := client.Query(exampleQuestion()); err != nil {
			t.Fatal(err)
		}
		if calls != 2 || client.servers[client.pool[0].Server].backoff != 0 {
			t.Fatal("successful reply did not reset failure backoff")
		}
	})
}

func TestClientNegativeAnswerDoesNotBlockOtherNames(t *testing.T) {
	client := stubClient(t, 1, 0, func(_ int, _ context.Context, msg *dns.Msg) (*dns.Msg, error) {
		response := exampleReply(msg)
		if msg.Question[0].Name == "missing.example.test." {
			response.Rcode, response.Answer = dns.RcodeNameError, nil
		}
		return response, nil
	})
	if _, err := client.Query(new(dns.Msg).SetQuestion("missing.example.test.", dns.TypeA)); !errors.Is(err, ErrNoAnswer) {
		t.Fatalf("negative query = %v", err)
	}
	if _, err := client.Query(exampleQuestion()); err != nil {
		t.Fatalf("negative reply blocked a different name: %v", err)
	}
}

func TestClientCancellationAndCloseJoinActiveQuery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		started, exited := make(chan struct{}), make(chan struct{})
		client := stubClient(t, 1, 0, func(_ int, ctx context.Context, _ *dns.Msg) (*dns.Msg, error) {
			close(started)
			<-ctx.Done()
			close(exited)
			return nil, ctx.Err()
		})
		result := make(chan error, 1)
		go func() {
			_, err := client.Query(exampleQuestion())
			result <- err
		}()
		<-started
		ctx, cancel := context.WithCancel(context.Background())
		queued := make(chan error, 1)
		go func() {
			_, err := client.QueryContext(ctx, exampleQuestion())
			queued <- err
		}()
		synctest.Wait()
		cancel()
		if err := <-queued; !errors.Is(err, context.Canceled) {
			t.Fatalf("queued cancellation = %v", err)
		}
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
		select {
		case <-exited:
		default:
			t.Fatal("Close returned before its worker exited")
		}
		if err := <-result; !errors.Is(err, ErrClosed) {
			t.Fatalf("active query after Close = %v", err)
		}
		if _, err := client.Query(exampleQuestion()); !errors.Is(err, ErrClosed) {
			t.Fatalf("closed query = %v", err)
		}
	})
}

func TestResolverCatalogFailsClosed(t *testing.T) {
	for _, data := range []string{
		"", "[]", "[invalid", "- name: plain\n  server: udp://8.8.8.8",
		"- name: encrypted\n  server: tls://dns.example.test\n- name: plain\n  server: udp://1.1.1.1",
	} {
		if pool, err := parseResolvers([]byte(data)); err == nil || len(pool) != 0 {
			t.Errorf("parseResolvers(%q) = %v, %v", data, pool, err)
		}
	}
}

func TestClientValidatesAndDeduplicatesPool(t *testing.T) {
	for _, pool := range []Resolvers{nil, {{Server: "tls://"}}, {{Server: "http://dns.example.test"}}} {
		if client, err := NewClient(pool, 0, time.Second); err == nil || client != nil {
			t.Errorf("NewClient(%v) = %v, %v", pool, client, err)
		}
	}
	pool := Resolvers{{Server: "tls://DNS.Example.test."}, {Server: "tls://dns.example.test:853"}}
	client, err := NewClient(pool, 0, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	if len(client.pool) != 1 {
		t.Fatalf("deduplicated pool = %v", client.pool)
	}
	if _, err := NewClient(pool, -1, time.Second); err == nil {
		t.Fatal("negative retries accepted")
	}
	if _, err := NewClient(pool, 0, 0); err == nil {
		t.Fatal("zero timeout accepted")
	}
}

func TestResponseValidation(t *testing.T) {
	request := exampleQuestion()
	for name, mutate := range map[string]func(*dns.Msg){
		"request packet":     func(msg *dns.Msg) { msg.Response = false },
		"wrong ID":           func(msg *dns.Msg) { msg.Id++ },
		"unrelated answer":   func(msg *dns.Msg) { msg.Answer[0].Header().Name = "other.test." },
		"wrong answer class": func(msg *dns.Msg) { msg.Answer[0].Header().Class = dns.ClassCHAOS },
		"wrong question": func(msg *dns.Msg) {
			msg.Question = []dns.Question{{Name: "other.test.", Qtype: dns.TypeA, Qclass: dns.ClassINET}}
		},
		"wrong record type": func(msg *dns.Msg) {
			msg.Answer = []dns.RR{&dns.TXT{Hdr: dns.RR_Header{Rrtype: dns.TypeTXT}, Txt: []string{"unexpected"}}}
		},
		"truncated":      func(msg *dns.Msg) { msg.Truncated = true },
		"server failure": func(msg *dns.Msg) { msg.Rcode = dns.RcodeServerFailure },
	} {
		t.Run(name, func(t *testing.T) {
			response := exampleReply(request)
			mutate(response)
			if err := validateResponse(request, response); err == nil {
				t.Fatal("invalid response accepted")
			}
		})
	}
}
