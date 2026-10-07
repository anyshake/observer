package archiver

import (
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestLifecycle(t *testing.T) {
	obj, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }))

	if svc.GetName() != "Archiver Service" {
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
		t.Fatal("archiver is disabled after Init")
	}
	if err := handler.SettingsSet(ID, "enabled", action.Bool, 0, false); err != nil {
		t.Fatalf("disable archiver: %v", err)
	}
	if svc.IsEnabled() {
		t.Fatal("archiver stayed enabled")
	}
	if err := handler.SettingsSet(ID, "enabled", action.Bool, 0, true); err != nil {
		t.Fatalf("enable archiver: %v", err)
	}

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() {
		_ = svc.Stop()
	})
	hw.WaitSubscribed(t)
	waitRunning(t, svc, true)

	testsupport.EnsureSeisTable(t, obj, fixed)
	for i := range RECORDS_INSERT_INTERVAL {
		hw.Emit(explorer.Event{
			Timestamp:  fixed.Add(time.Duration(i) * time.Second),
			SampleRate: 50,
			ChannelData: []explorer.ChannelData{{
				ChannelCode: "EHZ",
				ChannelId:   1,
				Data:        []int32{int32(i), int32(i + 1)},
			}},
		})
	}
	records, err := handler.SeisRecordsQuery(fixed.Add(-time.Second), fixed.Add(time.Minute))
	if err != nil {
		t.Fatalf("query records: %v", err)
	}
	if len(records) != RECORDS_INSERT_INTERVAL {
		t.Fatalf("inserted records = %d, want %d", len(records), RECORDS_INSERT_INTERVAL)
	}

	if err := svc.Restart(); err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	hw.WaitSubscribed(t)
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d, want 1", svc.GetStatus().GetRestarts())
	}

	if err := svc.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	waitRunning(t, svc, false)
}

func TestInitRequiresStationSettings(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	svc := New(testsupport.NewHardware(), handler, timesource.New(nil))
	if err := svc.Init(); err == nil {
		t.Fatal("Init() succeeded without station settings")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
}

func waitRunning(t *testing.T, svc *ArchiverServiceImpl, running bool) {
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
