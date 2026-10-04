package ntpclient

import (
	"context"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/beevik/ntp"
)

func TestCloseCancelsAndJoinsProbes(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, _ := newTestClient(t, 12, 5)
		var started, finished atomic.Int32
		release := make(chan struct{})
		client.query = func(string, ntp.QueryOptions) (*ntp.Response, error) {
			started.Add(1)
			<-client.ctx.Done()
			<-release
			finished.Add(1)
			return nil, context.Canceled
		}
		queryDone := make(chan error, 1)
		go func() {
			_, _, err := client.Query()
			queryDone <- err
		}()
		synctest.Wait()
		closeDone := make(chan error, 1)
		go func() { closeDone <- client.Close() }()
		synctest.Wait()
		returned := false
		select {
		case err := <-closeDone:
			returned = true
			t.Errorf("Close returned before workers finished: %v", err)
		default:
		}
		close(release)
		if !returned {
			if err := <-closeDone; err != nil {
				t.Fatal(err)
			}
		}
		if got := finished.Load(); got != CONCURRENT_QUERIES || got != started.Load() {
			t.Fatalf("Close left active or queued probes: started=%d, finished=%d", started.Load(), got)
		}
		if err := <-queryDone; !errors.Is(err, context.Canceled) {
			t.Fatalf("interrupted query error = %v", err)
		}
		if _, _, err := client.Query(); !errors.Is(err, ErrClosed) {
			t.Fatalf("query after Close error = %v", err)
		}
		if _, err := client.QueryAverage(5); !errors.Is(err, ErrClosed) {
			t.Fatalf("average after Close error = %v", err)
		}
		if err := client.Close(); err != nil {
			t.Fatalf("repeated Close failed: %v", err)
		}
	})
}

func TestQueryContextCancelsWhileWaiting(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		client, recorder := newTestClient(t, 1, 0)
		activeCtx, cancelActive := context.WithCancel(context.Background())
		defer cancelActive()
		client.query = func(string, ntp.QueryOptions) (*ntp.Response, error) {
			<-activeCtx.Done()
			return nil, activeCtx.Err()
		}
		activeDone := make(chan error, 1)
		go func() {
			_, _, err := client.QueryContext(activeCtx)
			activeDone <- err
		}()
		synctest.Wait()

		queuedCtx, cancelQueued := context.WithCancel(context.Background())
		queuedDone := make(chan error, 1)
		go func() {
			_, err := client.QueryAverageContext(queuedCtx, 5)
			queuedDone <- err
		}()
		synctest.Wait()
		cancelQueued()
		if err := <-queuedDone; !errors.Is(err, context.Canceled) {
			t.Fatalf("queued query error = %v", err)
		}
		select {
		case err := <-activeDone:
			t.Fatalf("canceling waiter interrupted active query: %v", err)
		default:
		}
		cancelActive()
		if err := <-activeDone; !errors.Is(err, context.Canceled) {
			t.Fatalf("active query error = %v", err)
		}
		client.query = recorder.query
		if _, _, err := client.Query(); err == nil || recorder.count(client.pool[0]) != 0 {
			t.Fatal("cancellation bypassed the endpoint cooldown")
		}
		time.Sleep(minPollInterval)
		if _, _, err := client.Query(); err != nil {
			t.Fatalf("canceling one operation made the client unusable: %v", err)
		}
	})
}

func TestQueryContextAlreadyCanceled(t *testing.T) {
	client, recorder := newTestClient(t, 2, 5)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.QueryAverageContext(ctx, 5); !errors.Is(err, context.Canceled) {
		t.Fatalf("query error = %v", err)
	}
	for _, server := range client.pool {
		if recorder.count(server) != 0 {
			t.Errorf("already canceled query contacted %s", server)
		}
	}
}

func TestQueryCancellationClosesConnection(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		name := "cancel"
		if deadline {
			name = "deadline"
		}
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				conn, peer := net.Pipe()
				defer conn.Close()
				defer peer.Close()
				requestSeen := make(chan struct{})
				peerDone := make(chan struct{})
				go func() {
					defer close(peerDone)
					if _, err := io.ReadFull(peer, make([]byte, 48)); err != nil {
						t.Errorf("reading NTP request: %v", err)
					}
					close(requestSeen)
					_, _ = peer.Read(make([]byte, 1))
				}()
				done := make(chan error, 1)
				go func() {
					_, err := queryWithContext(ctx, "time.example.test", ntp.QueryOptions{
						Timeout: time.Hour,
						Dialer:  func(string, string) (net.Conn, error) { return conn, nil },
					}, ntp.QueryWithOptions)
					done <- err
				}()
				<-requestSeen
				wantErr := context.DeadlineExceeded
				if !deadline {
					wantErr = context.Canceled
					cancel()
				}
				if err := <-done; !errors.Is(err, wantErr) {
					t.Fatalf("query error = %v, want %v", err, wantErr)
				}
				<-peerDone
			})
		})
	}
}

func TestSmallPools(t *testing.T) {
	for _, tt := range []struct {
		name       string
		count      int
		failed     int
		second     time.Duration
		wantOffset time.Duration
		wantErr    bool
	}{
		{name: "one source", count: 1, wantOffset: 10 * time.Millisecond},
		{name: "one unavailable source", count: 1, failed: 1, wantErr: true},
		{name: "two agreeing sources", count: 2, second: 12 * time.Millisecond, wantOffset: 11 * time.Millisecond},
		{name: "second source unavailable", count: 2, failed: 2, wantOffset: 10 * time.Millisecond},
		{name: "first source unavailable", count: 2, failed: 1, second: 12 * time.Millisecond, wantOffset: 12 * time.Millisecond},
		{name: "two disagreeing sources", count: 2, second: time.Second, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				client, recorder := newTestClient(t, tt.count, 5)
				recorder.respond = func(server string, _ int) (*ntp.Response, error) {
					if tt.failed > 0 && server == client.pool[tt.failed-1] {
						return nil, errors.New("timeout")
					}
					offset := 10 * time.Millisecond
					if server != client.pool[0] {
						offset = tt.second
					}
					return testResponse(offset, 4*time.Millisecond), nil
				}
				for attempt := 1; attempt <= 2; attempt++ {
					offset, err := client.QueryAverage(5)
					if (err != nil) != tt.wantErr || !tt.wantErr && offset != tt.wantOffset {
						t.Fatalf("query %d = %v, %v", attempt, offset, err)
					}
					for _, server := range client.pool {
						if got := recorder.count(server); got != attempt {
							t.Errorf("requests to %s = %d, want %d", server, got, attempt)
						}
					}
					time.Sleep(2 * time.Minute)
				}
			})
		})
	}
}
