package helicorder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestFilesAndPlotData(t *testing.T) {
	logger.Init()
	database, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(t.Context(), testsupport.NewHardware(), handler, timesource.New(func() time.Time { return fixed }))
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "file_path", dir); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "cache_storage", CACHE_STORAGE_DISABLED); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "cache_path", dir); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if svc.GetName() == "" || svc.GetDescription() == "" || svc.GetStatus() == nil {
		t.Fatal("service metadata is incomplete")
	}
	if _, err := svc.GetAssetList(); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.GetAssetData(filepath.Join(dir, "missing.png")); err == nil {
		t.Fatal("asset data succeeded while stopped")
	}

	svc.dataProvider.setChannelCode("EHZ")
	if svc.dataProvider.GetPlotName() == "" || svc.dataProvider.GetStation() == "" || svc.dataProvider.GetChannel() != "EHZ" {
		t.Fatal("plot provider metadata is incomplete")
	}
	if svc.dataProvider.GetNetwork() == "" || svc.dataProvider.GetLocation() == "" {
		t.Fatal("plot provider codes are empty")
	}
	if svc.getHelicorderFileName(fixed, "EHZ") == "" {
		t.Fatal("empty helicorder file name")
	}
	if wait := svc.getDurationToNextTime(fixed); wait <= 0 || wait > time.Hour {
		t.Fatalf("duration = %s", wait)
	}
	if _, _, err := normalizeDateRange(fixed, fixed.Add(-24*time.Hour)); err == nil {
		t.Fatal("reversed date range was accepted")
	}

	day := time.Date(2026, time.October, 8, 0, 0, 0, 0, time.UTC)
	other := day.AddDate(0, 0, 1)
	if err := os.MkdirAll(filepath.Join(dir, day.Format("2006-01-02")), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, other.Format("2006-01-02")), 0o700); err != nil {
		t.Fatal(err)
	}
	plot := filepath.Join(dir, day.Format("2006-01-02"), "plot.png")
	if err := os.WriteFile(plot, []byte("png"), 0o600); err != nil {
		t.Fatal(err)
	}
	svc.imageFormat = "png"
	svc.status.SetIsRunning(true)
	assets, err := svc.GetAssetList()
	if err != nil || len(assets) != 1 {
		t.Fatalf("assets = %#v, %v", assets, err)
	}
	data, err := svc.GetAssetData(plot)
	if err != nil || string(data.Data) != "png" {
		t.Fatalf("asset = %#v, %v", data, err)
	}
	if _, err := svc.GetAssetData(filepath.Join(dir, "note.txt")); err == nil {
		t.Fatal("non-image asset was accepted")
	}
	if _, err := svc.GetAssetData(filepath.Join(t.TempDir(), "plot.png")); err == nil {
		t.Fatal("outside asset was accepted")
	}
	if err := svc.PurgeHelicorderFilesByDate(t.Context(), day, day); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, day.Format("2006-01-02"))); !os.IsNotExist(err) {
		t.Fatalf("purged directory still exists: %v", err)
	}
	if err := svc.cleanupHelicorderFiles(t.Context(), other); err != nil {
		t.Fatal(err)
	}
	if err := svc.PurgeHelicorderFiles(t.Context()); err != nil {
		t.Fatal(err)
	}

	cache, err := newQueryCache(CACHE_STORAGE_MEMORY, "")
	if err != nil {
		t.Fatal(err)
	}
	svc.dataProvider.queryCache = cache
	svc.dataProvider.ctx = context.Background()
	testsupport.EnsureSeisTable(t, database, fixed)
	record := model.SeisRecord{}
	if err := record.Encode(fixed, 2, []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{3, 4}}}); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(t.Context(), record); err != nil {
		t.Fatal(err)
	}
	plotted, err := svc.dataProvider.GetPlotData(fixed, fixed.Add(time.Second))
	if err != nil || len(plotted) == 0 {
		t.Fatalf("plot data = %#v, %v", plotted, err)
	}
	cached, err := svc.dataProvider.GetPlotData(fixed, fixed.Add(time.Second))
	if err != nil || len(cached) != len(plotted) {
		t.Fatalf("cached plot data = %#v, %v", cached, err)
	}
	canceled, stop := context.WithCancel(context.Background())
	stop()
	svc.dataProvider.ctx = canceled
	svc.dataProvider.queryCache = nil
	if _, err := svc.dataProvider.GetPlotData(fixed, fixed.Add(time.Second)); err == nil {
		t.Fatal("canceled plot query succeeded")
	}
}

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(t.Context(), testsupport.NewHardware(), handler, timesource.New(func() time.Time { return fixed }))
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "file_path", dir); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "cache_storage", CACHE_STORAGE_DISABLED); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "cache_path", dir); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled(t.Context()) && svc.GetStatus().GetIsRunning() {
		t.Fatal("unexpected running state")
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

func setConstraint(t *testing.T, handler *action.Handler, svc *HelicorderServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() == key {
			return constraint.Set(t.Context(), handler, value)
		}
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}

func waitRunning(t *testing.T, svc *HelicorderServiceImpl, running bool) {
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
