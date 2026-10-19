package frp_client

import (
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(t.Context(), "bad-address", handler, timesource.New(func() time.Time { return fixed }))
	if svc.GetName() != "FRP Client Service" || svc.GetDescription() == "" || svc.GetStatus() == nil {
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
	if _, _, err := svc.parseLocalServerAddr(); err == nil {
		t.Fatal("invalid local address was accepted")
	}
	svc.localServerAddr = "127.0.0.1:nope"
	if _, _, err := svc.parseLocalServerAddr(); err == nil {
		t.Fatal("invalid local port was accepted")
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "enabled", true); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "server_addr", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "server_port", int64(1)); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled(t.Context()) {
		t.Fatal("service stayed disabled")
	}
	if err := svc.Start(t.Context()); err == nil {
		t.Fatal("invalid local address was started")
	}

	svc.localServerAddr = "127.0.0.1:8080"
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

	if err := setConstraint(t, handler, svc, "use_domain_access", false); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "user", "alice"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := svc.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	waitRunning(t, svc, true)
	if err := svc.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestLogWriterFormatsFRPLines(t *testing.T) {
	writer := logWriter{}
	short := []byte("short")
	if n, err := writer.Write(short); err != nil || n != len(short) {
		t.Fatalf("short write = %d, %v", n, err)
	}
	for _, level := range []string{"I", "W", "E", "D"} {
		line := strings.Repeat("x", 25) + level + " message\n"
		if n, err := writer.Write([]byte(line)); err != nil || n != len(line) {
			t.Fatalf("level %s write = %d, %v", level, n, err)
		}
	}
}

func setConstraint(t *testing.T, handler *action.Handler, svc *FrpClientServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() == key {
			return constraint.Set(t.Context(), handler, value)
		}
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}

func waitRunning(t *testing.T, svc *FrpClientServiceImpl, running bool) {
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
