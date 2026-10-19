package seisevent

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"golang.org/x/sync/singleflight"
)

type fetchOutcome struct {
	events []Event
	err    error
}

func TestFetchEventsRetriesAfterLeaderCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var group singleflight.Group
		var calls atomic.Int32
		leaderCtx, cancelLeader := context.WithCancel(t.Context())
		defer cancelLeader()

		leader := make(chan fetchOutcome, 1)
		go func() {
			events, err := fetchEvents(leaderCtx, &group, "key", func(ctx context.Context) ([]Event, error) {
				calls.Add(1)
				<-ctx.Done()
				return nil, ctx.Err()
			})
			leader <- fetchOutcome{events, err}
		}()
		synctest.Wait()

		waiter := make(chan fetchOutcome, 1)
		go func() {
			events, err := fetchEvents(t.Context(), &group, "key", func(context.Context) ([]Event, error) {
				calls.Add(1)
				return []Event{{Event: "fresh"}}, nil
			})
			waiter <- fetchOutcome{events, err}
		}()
		synctest.Wait()

		cancelLeader()
		if got := <-leader; !errors.Is(got.err, context.Canceled) {
			t.Fatalf("leader error = %v, want context.Canceled", got.err)
		}
		got := <-waiter
		if got.err != nil || len(got.events) != 1 || got.events[0].Event != "fresh" {
			t.Fatalf("waiter events = %#v, error = %v", got.events, got.err)
		}
		if calls.Load() != 2 {
			t.Fatalf("fetch calls = %d, want 2", calls.Load())
		}
	})
}

func TestFetchEventsCanceledWaiterReturnsPromptly(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var group singleflight.Group
		var calls atomic.Int32
		release := make(chan struct{})

		leader := make(chan fetchOutcome, 1)
		go func() {
			events, err := fetchEvents(t.Context(), &group, "key", func(context.Context) ([]Event, error) {
				calls.Add(1)
				<-release
				return []Event{{Event: "shared"}}, nil
			})
			leader <- fetchOutcome{events, err}
		}()
		synctest.Wait()

		waiterCtx, cancelWaiter := context.WithCancel(t.Context())
		waiter := make(chan fetchOutcome, 1)
		go func() {
			events, err := fetchEvents(waiterCtx, &group, "key", func(context.Context) ([]Event, error) {
				calls.Add(1)
				return nil, errors.New("waiter must join the existing flight")
			})
			waiter <- fetchOutcome{events, err}
		}()
		synctest.Wait()

		cancelWaiter()
		if got := <-waiter; !errors.Is(got.err, context.Canceled) {
			t.Fatalf("waiter error = %v, want context.Canceled", got.err)
		}
		close(release)
		got := <-leader
		if got.err != nil || len(got.events) != 1 || got.events[0].Event != "shared" {
			t.Fatalf("leader events = %#v, error = %v", got.events, got.err)
		}
		if calls.Load() != 1 {
			t.Fatalf("fetch calls = %d, want 1", calls.Load())
		}
	})
}

func TestFetchEventsReturnsSourceErrorsAndPanics(t *testing.T) {
	var group singleflight.Group
	sourceErr := errors.New("source failed")
	if _, err := fetchEvents(t.Context(), &group, "error", func(context.Context) ([]Event, error) {
		return nil, sourceErr
	}); !errors.Is(err, sourceErr) {
		t.Fatalf("error = %v, want %v", err, sourceErr)
	}
	if _, err := fetchEvents(t.Context(), &group, "panic", func(context.Context) ([]Event, error) {
		panic("malformed response")
	}); err == nil {
		t.Fatal("panic was not converted to an error")
	}
}
