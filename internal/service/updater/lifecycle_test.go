package updater

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/internal/upgrade"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"
)

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	dir := t.TempDir()
	exe := filepath.Join(dir, "observer")
	helper := upgrade.NewHelper(exe, semver.New("1", "0", "0", ""), unibuild.New("test", "stable", "abc", "1"))
	svc := New(t.Context(), handler, timesource.New(func() time.Time { return fixed }), helper, exe, make(chan struct{}, 1))
	if svc.GetName() != "Updater Service" || svc.GetDescription() == "" || svc.GetStatus() == nil {
		t.Fatal("service metadata is incomplete")
	}
	if svc.IsEnabled(t.Context()) {
		t.Fatal("enabled before initialization")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "enabled", true); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled(t.Context()) {
		t.Fatal("service stayed disabled")
	}
	if err := svc.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testsupport.StopService(t, svc) })
	waitRunning(t, svc, true)
	if err := svc.Restart(t.Context()); err != nil {
		t.Fatal(err)
	}
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d", svc.GetStatus().GetRestarts())
	}
	if err := svc.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	waitRunning(t, svc, false)
}

func TestRepeatedRestartRequestsDoNotBlock(t *testing.T) {
	restartChan := make(chan struct{}, 1)
	svc := &UpdaterServiceImpl{restartChan: restartChan}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 3 {
			svc.requestRestart()
		}
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("repeated restart requests blocked")
	}
	if len(restartChan) != 1 {
		t.Fatalf("pending restart requests = %d, want 1", len(restartChan))
	}
}

func TestStartRejectsReadOnlyDirectory(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	locked := t.TempDir()
	t.Cleanup(func() { _ = os.Chmod(locked, 0o600) })
	if err := os.Chmod(locked, 0o555); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(locked, "observer")
	helper := upgrade.NewHelper(exe, semver.New("1", "0", "0", ""), unibuild.New("test", "stable", "abc", "1"))
	svc := New(t.Context(), handler, timesource.New(nil), helper, exe, make(chan struct{}, 1))
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := svc.Start(t.Context()); err == nil {
		t.Fatal("read-only directory was accepted")
	}
}

func setConstraint(t *testing.T, handler *action.Handler, svc *UpdaterServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() == key {
			return constraint.Set(t.Context(), handler, value)
		}
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}

func waitRunning(t *testing.T, svc *UpdaterServiceImpl, running bool) {
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
