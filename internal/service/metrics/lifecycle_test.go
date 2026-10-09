package metrics

import (
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestLifecycle(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	hw := testsupport.NewHardware()
	hw.Status.SetUpdatedAt(time.Unix(10, 0))
	hw.Status.IncrementFrames()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }), semver.New("1", "2", "3", ""), unibuild.New("linux_arm32_v7a", "stable", "abc123", "1700000000"))

	if svc.GetName() != "Metrics Service" {
		t.Fatalf("name = %q", svc.GetName())
	}
	if svc.GetDescription() == "" {
		t.Fatal("description is empty")
	}
	if svc.GetStatus() == nil {
		t.Fatal("status is nil")
	}
	if svc.IsEnabled() {
		t.Fatal("enabled before settings exist")
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
	if !svc.IsEnabled() {
		t.Fatal("metrics is disabled after Init")
	}
	if err := handler.SettingsSet(ID, "enabled", action.Bool, 0, false); err != nil {
		t.Fatal(err)
	}
	if svc.IsEnabled() {
		t.Fatal("metrics stayed enabled")
	}
	if err := handler.SettingsSet(ID, "enabled", action.Bool, 0, true); err != nil {
		t.Fatal(err)
	}

	withoutToolchain := New(hw, handler, timesource.New(nil), semver.New("0", "0", "1", ""), unibuild.New("", "dev", "commit", "0"))
	if len(withoutToolchain.getBuildInfo()) < 4 {
		t.Fatal("build info is missing the base attributes")
	}
	info := svc.getBuildInfo()
	if !hasAttribute(info, "service.build.commit", "abc123") {
		t.Fatalf("build info = %#v", info)
	}
	if !hasAttribute(info, "service.build.toolchain.name", "linux-arm32-v7a") {
		t.Fatalf("toolchain was not included in %#v", info)
	}
	svc.reportCurrentStatus(noop.NewTracerProvider().Tracer("test"))

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() { _ = svc.Stop() })
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
}

func TestInitRequiresDatabase(t *testing.T) {
	svc := New(testsupport.NewHardware(), action.NewHandler(nil), timesource.New(nil), semver.New("0", "0", "1", ""), unibuild.New("", "", "", ""))
	if err := svc.Init(); err == nil {
		t.Fatal("Init() succeeded without a database")
	}
}

func hasAttribute(values []attribute.KeyValue, key, want string) bool {
	for _, value := range values {
		if string(value.Key) == key && value.Value.AsString() == want {
			return true
		}
	}
	return false
}

func waitRunning(t *testing.T, svc *MetricsServiceImpl, running bool) {
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
