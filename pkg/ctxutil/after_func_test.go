package ctxutil

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"
)

func TestAfterFuncStopPreventsCallback(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		var calls atomic.Int32
		stop := AfterFunc(ctx, func() { calls.Add(1) })
		stop()
		stop()
		cancel()
		synctest.Wait()
		if got := calls.Load(); got != 0 {
			t.Fatalf("callback ran %d times after stop", got)
		}
	})
}

func TestAfterFuncStopWaitsForRunningCallback(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		started := make(chan struct{})
		release := make(chan struct{})
		var finished atomic.Bool
		stop := AfterFunc(ctx, func() {
			close(started)
			<-release
			finished.Store(true)
		})
		cancel()
		<-started
		stopped := make(chan struct{})
		go func() {
			stop()
			close(stopped)
		}()
		synctest.Wait()
		select {
		case <-stopped:
			t.Fatal("stop returned before the callback finished")
		default:
		}
		close(release)
		<-stopped
		if !finished.Load() {
			t.Fatal("callback did not finish before stop returned")
		}
		stop()
	})
}
