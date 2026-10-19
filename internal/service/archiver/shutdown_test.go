package archiver

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestSubscriptionFailureCanStop(t *testing.T) {
	hw := testsupport.NewHardware()
	hw.SubscribeErr = errors.New("subscription failed")
	svc := New(t.Context(), hw, nil, timesource.New(nil))
	if err := svc.Start(t.Context()); err == nil {
		t.Fatal("Start succeeded despite a subscription failure")
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("failed subscription left the service running")
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := svc.Stop(ctx); err != nil {
		t.Fatalf("Stop after subscription failure: %v", err)
	}
}

func TestCallbackPanicStopsServiceAndAllowsRestart(t *testing.T) {
	hw := testsupport.NewHardware()
	svc := New(t.Context(), hw, nil, timesource.New(nil))
	if err := svc.Start(t.Context()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() { testsupport.StopService(t, svc) })
	hw.WaitSubscribed(t)
	if err := svc.Start(t.Context()); err == nil {
		t.Fatal("duplicate Start succeeded")
	}

	event := explorer.Event{
		Timestamp:   time.Now(),
		SampleRate:  100,
		ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{1, 2, 3}}},
	}
	for range RECORDS_INSERT_INTERVAL {
		hw.Emit(event)
	}
	waitRunning(t, svc, false)

	if err := svc.Start(t.Context()); err != nil {
		t.Fatalf("Start after callback panic: %v", err)
	}
	hw.WaitSubscribed(t)
	waitRunning(t, svc, true)
}

func TestStopWaitsForWorkersWithinCallerDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		svc := New(t.Context(), testsupport.NewHardware(), nil, timesource.New(nil))
		svc.wg.Add(1)
		go func() {
			defer svc.wg.Done()
			time.Sleep(12 * time.Second)
		}()
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		defer cancel()
		started := time.Now()
		if err := svc.Stop(ctx); err != nil {
			t.Fatal(err)
		}
		if elapsed := time.Since(started); elapsed != 12*time.Second {
			t.Fatalf("Stop returned after %s before the worker finished", elapsed)
		}
	})
}

func TestInitHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	svc := New(t.Context(), testsupport.NewHardware(), nil, timesource.New(nil))
	if err := svc.Init(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Init with canceled context = %v", err)
	}
}

func TestStopTimeoutKeepsTrackingWorker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		svc := New(t.Context(), testsupport.NewHardware(), nil, timesource.New(nil))
		svc.wg.Add(1)
		finished := make(chan struct{})
		go func() {
			defer svc.wg.Done()
			time.Sleep(12 * time.Second)
			close(finished)
		}()
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		if err := svc.Stop(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Stop before worker completion = %v", err)
		}
		ctx, cancel = context.WithTimeout(t.Context(), 30*time.Second)
		defer cancel()
		if err := svc.Stop(ctx); err != nil {
			t.Fatalf("waiting again for shutdown: %v", err)
		}
		select {
		case <-finished:
		default:
			t.Fatal("Stop succeeded while the worker was still running")
		}
	})
}
