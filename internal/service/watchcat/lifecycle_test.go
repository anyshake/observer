package watchcat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestLifecycle(t *testing.T) {
	hw := testsupport.NewHardware()
	svc := New(t.Context(), hw, timesource.New(func() time.Time {
		return time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	}))

	if svc.GetName() != "WatchCat Service" {
		t.Fatalf("name = %q", svc.GetName())
	}
	if svc.GetDescription() == "" {
		t.Fatal("description is empty")
	}
	if svc.GetStatus() == nil {
		t.Fatal("status is nil")
	}
	if !svc.IsEnabled(t.Context()) {
		t.Fatal("watchcat is always enabled")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("service is running before Start")
	}

	if err := svc.Start(t.Context()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() { testsupport.StopService(t, svc) })
	waitRunning(t, svc, true)

	if err := svc.Restart(t.Context()); err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d, want 1", svc.GetStatus().GetRestarts())
	}

	if err := svc.Stop(t.Context()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	waitRunning(t, svc, false)
	if svc.GetStatus().GetStoppedAt().IsZero() {
		t.Fatal("stopped time was not recorded")
	}
}

func TestStopBeforeStart(t *testing.T) {
	svc := New(t.Context(), testsupport.NewHardware(), timesource.New(nil))
	if err := svc.Stop(t.Context()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("service is running after Stop")
	}
}

func TestStartFollowsApplicationContext(t *testing.T) {
	appCtx, cancelApp := context.WithCancel(t.Context())
	defer cancelApp()
	svc := New(appCtx, testsupport.NewHardware(), timesource.New(nil))
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	if err := svc.Start(canceled); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start with canceled caller context = %v", err)
	}
	if err := svc.Start(t.Context()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() { testsupport.StopService(t, svc) })
	waitRunning(t, svc, true)

	cancelApp()
	exited := make(chan struct{})
	go func() {
		svc.wg.Wait()
		close(exited)
	}()
	select {
	case <-exited:
	case <-time.After(2 * time.Second):
		t.Fatal("worker kept running after the application context was canceled")
	}
	if err := svc.Start(t.Context()); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start after application shutdown = %v", err)
	}
}

func waitRunning(t *testing.T, svc *WatchCatServiceImpl, running bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if svc.GetStatus().GetIsRunning() == running {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("running = %v, want %v", svc.GetStatus().GetIsRunning(), running)
}
