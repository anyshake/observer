package message_test

import (
	"errors"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/message"
)

func TestBusSubscriptionLifecycle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := message.NewBus[int]("samples")
		defer bus.Close()
		first, second := make(chan int, 2), make(chan int, 2)
		if err := bus.Subscribe("first", message.SubscriptionOptions{}, func(v int) { first <- v }); err != nil {
			t.Fatal(err)
		}
		if err := bus.Subscribe("second", message.SubscriptionOptions{}, func(v int) { second <- v }); err != nil {
			t.Fatal(err)
		}
		if err := bus.Subscribe("nil", message.SubscriptionOptions{}, nil); !errors.Is(err, message.ErrInvalidHandler) {
			t.Fatalf("nil handler error = %v, want ErrInvalidHandler", err)
		}
		if err := bus.Subscribe("first", message.SubscriptionOptions{}, func(int) {}); err == nil {
			t.Fatal("duplicate subscriber accepted")
		}
		if got := bus.Publish(10); got != (message.PublishResult{Delivered: 2}) {
			t.Fatalf("Publish() = %+v, want two deliveries", got)
		}
		synctest.Wait()
		if got := receive(t, first); got != 10 {
			t.Fatalf("first subscriber received %d, want 10", got)
		}
		if got := receive(t, second); got != 10 {
			t.Fatalf("second subscriber received %d, want 10", got)
		}
		if err := bus.Unsubscribe("first"); err != nil {
			t.Fatal(err)
		}
		if err := bus.Unsubscribe("first"); err == nil {
			t.Fatal("unsubscribing an unknown client succeeded")
		}
		if got := bus.Publish(20); got != (message.PublishResult{Delivered: 1}) {
			t.Fatalf("Publish() after unsubscribe = %+v, want one delivery", got)
		}
		synctest.Wait()
		if got := receive(t, second); got != 20 {
			t.Fatalf("remaining subscriber received %d, want 20", got)
		}
		if len(first) != 0 {
			t.Fatal("unsubscribed client received another event")
		}
		bus.Close()
		bus.Close()
		if err := bus.Subscribe("late", message.SubscriptionOptions{}, func(int) {}); !errors.Is(err, message.ErrBusClosed) {
			t.Fatalf("Subscribe() after Close() error = %v, want ErrBusClosed", err)
		}
		if got := bus.Publish(30); got != (message.PublishResult{}) {
			t.Fatalf("Publish() after Close() = %+v, want no deliveries", got)
		}
		want := message.Stats{Published: 3, Delivered: 3}
		if got := bus.GetStats(); got != want {
			t.Fatalf("stats = %+v, want %+v", got, want)
		}
	})
}

func TestBusOverflowPolicies(t *testing.T) {
	tests := []struct {
		name       string
		policy     message.OverflowPolicy
		wantResult message.PublishResult
		wantNext   int
		wantErr    error
		wantStats  message.Stats
	}{
		{
			name: "drop newest", policy: message.OverflowDropNewest,
			wantResult: message.PublishResult{Dropped: 1}, wantNext: 2,
			wantStats: message.Stats{Published: 3, Delivered: 2, Dropped: 1, Subscribers: 1},
		},
		{
			name: "drop oldest", policy: message.OverflowDropOldest,
			wantResult: message.PublishResult{Delivered: 1, Dropped: 1}, wantNext: 3,
			wantStats: message.Stats{Published: 3, Delivered: 3, Dropped: 1, Subscribers: 1},
		},
		{
			name: "block times out", policy: message.OverflowBlock,
			wantResult: message.PublishResult{Dropped: 1}, wantNext: 2, wantErr: message.ErrPublishTimeout,
			wantStats: message.Stats{Published: 3, Delivered: 2, Dropped: 1, Subscribers: 1},
		},
		{
			name: "disconnect", policy: message.OverflowDisconnect,
			wantResult: message.PublishResult{Disconnected: 1}, wantErr: message.ErrSubscriberTooSlow,
			wantStats: message.Stats{Published: 3, Delivered: 2, Disconnected: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				bus := message.NewBus[int]("samples")
				defer bus.Close()
				release := make(chan struct{})
				unblock := sync.OnceFunc(func() { close(release) })
				defer unblock()
				received := make(chan int, 3)
				reported := make(chan error, 3)
				err := bus.Subscribe("slow", message.SubscriptionOptions{
					BufferSize: 1, Overflow: tt.policy, EnqueueTimeout: time.Second,
					OnError: func(err error) { reported <- err },
				}, func(v int) {
					received <- v
					if v == 1 {
						<-release
					}
				})
				if err != nil {
					t.Fatal(err)
				}
				if got := bus.Publish(1); got != (message.PublishResult{Delivered: 1}) {
					t.Fatalf("first Publish() = %+v", got)
				}
				synctest.Wait()
				if got := receive(t, received); got != 1 {
					t.Fatalf("first event = %d, want 1", got)
				}
				// The handler is blocked, so the second event fills the queue.
				if got := bus.Publish(2); got != (message.PublishResult{Delivered: 1}) {
					t.Fatalf("second Publish() = %+v", got)
				}
				if got := bus.Publish(3); got != tt.wantResult {
					t.Fatalf("overflow result = %+v, want %+v", got, tt.wantResult)
				}
				synctest.Wait()
				if tt.wantErr != nil {
					if err := receive(t, reported); !errors.Is(err, tt.wantErr) {
						t.Fatalf("reported error = %v, want %v", err, tt.wantErr)
					}
				} else if len(reported) != 0 {
					t.Fatalf("unexpected error: %v", receive(t, reported))
				}
				if got := bus.GetStats(); got != tt.wantStats {
					t.Fatalf("stats = %+v, want %+v", got, tt.wantStats)
				}
				unblock()
				synctest.Wait()
				if tt.wantNext != 0 {
					if got := receive(t, received); got != tt.wantNext {
						t.Fatalf("queued event = %d, want %d", got, tt.wantNext)
					}
				}
				if len(received) != 0 {
					t.Fatal("dropped or disconnected event reached the handler")
				}
			})
		})
	}
}

func TestBusSubscriberPanicIsolated(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := message.NewBus[int]("samples")
		defer bus.Close()
		reported := make(chan error, 1)
		received := make(chan int, 2)
		if err := bus.Subscribe("broken", message.SubscriptionOptions{
			OnError: func(err error) { reported <- err },
		}, func(int) { panic("handler failed") }); err != nil {
			t.Fatal(err)
		}
		if err := bus.Subscribe("healthy", message.SubscriptionOptions{}, func(v int) { received <- v }); err != nil {
			t.Fatal(err)
		}
		bus.Publish(1)
		synctest.Wait()
		if err := receive(t, reported); !errors.Is(err, message.ErrSubscriberPanicked) {
			t.Fatalf("panic error = %v, want ErrSubscriberPanicked", err)
		}
		if got := receive(t, received); got != 1 {
			t.Fatalf("healthy subscriber received %d, want 1", got)
		}
		if got := bus.Publish(2); got != (message.PublishResult{Delivered: 1}) {
			t.Fatalf("Publish() after panic = %+v, want one delivery", got)
		}
		synctest.Wait()
		if got := receive(t, received); got != 2 {
			t.Fatalf("healthy subscriber received %d, want 2", got)
		}
		want := message.Stats{Published: 2, Delivered: 3, Disconnected: 1, Subscribers: 1}
		if got := bus.GetStats(); got != want {
			t.Fatalf("stats after panic = %+v, want %+v", got, want)
		}
	})
}

func TestBusCloseUnblocksPublisher(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := message.NewBus[int]("samples")
		defer bus.Close()
		release := make(chan struct{})
		defer close(release)
		reported := make(chan error, 1)
		if err := bus.Subscribe("slow", message.SubscriptionOptions{
			BufferSize: 1, Overflow: message.OverflowBlock, EnqueueTimeout: time.Hour,
			OnError: func(err error) { reported <- err },
		}, func(int) { <-release }); err != nil {
			t.Fatal(err)
		}
		bus.Publish(1)
		synctest.Wait()
		bus.Publish(2)
		result := make(chan message.PublishResult, 1)
		go func() { result <- bus.Publish(3) }()
		synctest.Wait()
		if len(result) != 0 {
			t.Fatal("publisher did not block on the full queue")
		}
		bus.Close()
		synctest.Wait()
		if got := receive(t, result); got != (message.PublishResult{Dropped: 1}) {
			t.Fatalf("blocked Publish() = %+v, want one dropped event", got)
		}
		if err := receive(t, reported); !errors.Is(err, message.ErrBusClosed) {
			t.Fatalf("close error = %v, want ErrBusClosed", err)
		}
	})
}

func receive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	default:
		t.Fatal("expected a queued value after synctest.Wait()")
		var zero T
		return zero
	}
}
