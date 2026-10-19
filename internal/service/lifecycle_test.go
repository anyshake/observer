package service

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestCanceledLifecycleWaitDoesNotRunCleanup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var lifecycle Lifecycle
		if err := lifecycle.Lock(t.Context()); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		called := false
		if err := lifecycle.Stop(ctx, func() error {
			called = true
			return nil
		}); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Stop waiting for lifecycle = %v", err)
		}
		lifecycle.Unlock()
		synctest.Wait()
		if called {
			t.Fatal("a canceled operation started cleanup after the lifecycle was released")
		}
	})
}

func TestStopTimeoutKeepsLifecycleLockedUntilCleanupFinishes(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var lifecycle Lifecycle
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		finished := false
		if err := lifecycle.Stop(ctx, func() error {
			time.Sleep(12 * time.Second)
			finished = true
			return nil
		}); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Stop before cleanup finished = %v", err)
		}
		if err := lifecycle.Lock(t.Context()); err != nil {
			t.Fatal(err)
		}
		defer lifecycle.Unlock()
		if !finished {
			t.Fatal("lifecycle was released before cleanup finished")
		}
	})
}
