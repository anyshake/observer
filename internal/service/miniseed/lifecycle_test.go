package miniseed

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	hw := testsupport.NewHardware()
	hw.Config.SetGnssAvailability(true)
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }))

	if svc.GetName() != "MiniSEED Service" || svc.GetDescription() == "" || svc.GetStatus() == nil {
		t.Fatal("service metadata is incomplete")
	}
	if svc.IsEnabled() {
		t.Fatal("enabled before settings exist")
	}
	if _, err := svc.GetAssetData("missing.mseed"); err == nil {
		t.Fatal("asset data succeeded before start")
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "file_path", dir); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "enabled", true); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled() {
		t.Fatal("service stayed disabled")
	}

	if err := svc.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.Stop() })
	hw.WaitSubscribed(t)
	waitRunning(t, svc, true)

	for i := range MINISEED_APPEND_INTERVAL {
		hw.Emit(explorer.Event{
			Timestamp:     fixed.Add(time.Duration(i) * time.Second),
			SampleRate:    100,
			GNSSAvailable: true,
			ChannelData: []explorer.ChannelData{{
				ChannelCode: "EHZ",
				ChannelId:   1,
				Data:        []int32{int32(i), int32(i + 1)},
			}},
		})
	}
	assets, err := svc.GetAssetList()
	if err != nil || len(assets) != 1 {
		t.Fatalf("assets = %#v, %v", assets, err)
	}
	data, err := svc.GetAssetData(assets[0].FilePath)
	if err != nil || len(data.Data) == 0 || data.FileName != assets[0].FileName {
		t.Fatalf("asset data = %#v, %v", data, err)
	}
	if _, err := svc.GetAssetData(filepath.Join(dir, "note.txt")); err == nil {
		t.Fatal("non-miniseed asset was accepted")
	}
	if _, err := svc.GetAssetData(filepath.Join(t.TempDir(), "wave.mseed")); err == nil {
		t.Fatal("outside asset was accepted")
	}

	if err := os.WriteFile(svc.dataSequence.filePath, []byte("not-gob"), 0o600); err != nil {
		t.Fatal(err)
	}
	later := fixed.Add(time.Hour)
	for i := range MINISEED_APPEND_INTERVAL {
		hw.Emit(explorer.Event{
			Timestamp:     later.Add(time.Duration(i) * time.Second),
			SampleRate:    100,
			GNSSAvailable: true,
			ChannelData: []explorer.ChannelData{{
				ChannelCode: "EHZ",
				ChannelId:   1,
				Data:        []int32{9, int32(i)},
			}},
		})
	}

	if err := svc.PurgeMiniSeedFilesByDate(fixed.Add(48*time.Hour), fixed); err == nil {
		t.Fatal("reversed purge range was accepted")
	}
	if err := svc.PurgeMiniSeedFilesByDate(fixed, fixed); err != nil {
		t.Fatal(err)
	}
	assets, err = svc.GetAssetList()
	if err != nil || len(assets) != 0 {
		t.Fatalf("assets after purge = %#v, %v", assets, err)
	}
	if err := svc.PurgeMiniSeedFiles(); err != nil {
		t.Fatal(err)
	}

	if err := svc.Restart(); err != nil {
		t.Fatal(err)
	}
	hw.WaitSubscribed(t)
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d", svc.GetStatus().GetRestarts())
	}
	if err := svc.Stop(); err != nil {
		t.Fatal(err)
	}
	waitRunning(t, svc, false)
}

func TestInitRequiresStationSettings(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	svc := New(testsupport.NewHardware(), handler, timesource.New(time.Now))
	if err := svc.Init(); err == nil {
		t.Fatal("Init succeeded without station settings")
	}
}

func setConstraint(t *testing.T, handler *action.Handler, svc *MiniSeedServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() != key {
			continue
		}
		return constraint.Set(handler, value)
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}

func waitRunning(t *testing.T, svc *MiniSeedServiceImpl, running bool) {
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
