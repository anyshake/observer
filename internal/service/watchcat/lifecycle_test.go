package watchcat

import (
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestLifecycle(t *testing.T) {
	hw := testsupport.NewHardware()
	svc := New(hw, timesource.New(func() time.Time {
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
	if !svc.IsEnabled() {
		t.Fatal("watchcat is always enabled")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("service is running before Start")
	}

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() {
		if err := svc.Stop(); err != nil {
			t.Errorf("cleanup Stop() error = %v", err)
		}
	})
	waitRunning(t, svc, true)

	if err := svc.Restart(); err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d, want 1", svc.GetStatus().GetRestarts())
	}

	if err := svc.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	waitRunning(t, svc, false)
	if svc.GetStatus().GetStoppedAt().IsZero() {
		t.Fatal("stopped time was not recorded")
	}
}

func TestStopBeforeStart(t *testing.T) {
	svc := New(testsupport.NewHardware(), timesource.New(nil))
	if err := svc.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("service is running after Stop")
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
