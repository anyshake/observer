package testsupport

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/hardware/explorer"
)

func TestHelpers(t *testing.T) {
	database, handler := OpenDAO(t)
	ExerciseConstraint(t, handler, &config.StationNameConfigConstraintImpl{})
	InitStation(t, handler)
	EnsureSeisTable(t, database, time.Date(2026, time.October, 8, 0, 0, 0, 0, time.UTC))

	hw := NewHardware()
	_ = hw.GetConfig()
	_ = hw.GetStatus()
	if hw.GetDeviceId() == "" {
		t.Fatal("hardware defaults are empty")
	}
	lat, lon, alt, err := hw.GetCoordinates(false)
	if err != nil || lat == 0 || lon == 0 || alt == 0 {
		t.Fatalf("coordinates = %v %v %v %v", lat, lon, alt, err)
	}
	if temp, err := hw.GetTemperature(); err != nil || temp == 0 {
		t.Fatalf("temperature = %v %v", temp, err)
	}
	if _, err := hw.GetMetadata("", "", "", "", "", "", "", false); err != nil {
		t.Fatal(err)
	}
	if _, _, err := hw.Open(nil); err != nil {
		t.Fatal(err)
	}
	if err := hw.Flush(); err != nil || hw.Flushes != 1 {
		t.Fatal(err)
	}
	if err := hw.Subscribe("stream", func(explorer.Event) {}); err != nil {
		t.Fatal(err)
	}
	hw.WaitSubscribed(t)
	hw.Emit(explorer.Event{})
	if err := hw.Unsubscribe("stream"); err != nil {
		t.Fatal(err)
	}
	if err := hw.SubscribeRealtime("realtime", func(explorer.Event) {}); err != nil {
		t.Fatal(err)
	}
	hw.WaitSubscribedRealtime(t)
	hw.EmitRealtime(explorer.Event{})
	if err := hw.UnsubscribeRealtime("realtime"); err != nil {
		t.Fatal(err)
	}

	hw.OpenErr = errors.New("open")
	if _, _, err := hw.Open(context.Background()); err == nil {
		t.Fatal("open error was ignored")
	}
	hw.CloseErr = errors.New("close")
	if err := hw.Close(); err == nil {
		t.Fatal("close error was ignored")
	}
	hw.FlushErr = errors.New("flush")
	if err := hw.Flush(); err == nil {
		t.Fatal("flush error was ignored")
	}
	hw.SubscribeErr = errors.New("subscribe")
	if err := hw.Subscribe("stream", nil); err == nil {
		t.Fatal("subscribe error was ignored")
	}
	hw.SubscribeRealtimeErr = errors.New("realtime")
	if err := hw.SubscribeRealtime("realtime", nil); err == nil {
		t.Fatal("realtime subscribe error was ignored")
	}
	hw.CoordErr = errors.New("coord")
	if _, _, _, err := hw.GetCoordinates(true); err == nil {
		t.Fatal("coordinate error was ignored")
	}
	hw.TempErr = errors.New("temp")
	if _, err := hw.GetTemperature(); err == nil {
		t.Fatal("temperature error was ignored")
	}
	hw.MetaErr = errors.New("meta")
	if _, err := hw.GetMetadata("", "", "", "", "", "", "", false); err == nil {
		t.Fatal("metadata error was ignored")
	}
}
