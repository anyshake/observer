package forwarder

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestLifecycle(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	port := freePort(t)
	if err := handler.SettingsSet(ID, "listen_host", action.String, 0, "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if err := handler.SettingsSet(ID, "listen_port", action.Int, 0, int64(port)); err != nil {
		t.Fatal(err)
	}

	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }))

	if svc.GetName() != "TCP Forwarder Service" {
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
		t.Fatal("forwarder is disabled after Init")
	}
	if svc.GetListenPort() != port {
		t.Fatalf("listen port = %d, want %d", svc.GetListenPort(), port)
	}
	if err := handler.SettingsSet(ID, "enabled", action.Bool, 0, false); err != nil {
		t.Fatal(err)
	}
	if svc.IsEnabled() {
		t.Fatal("forwarder stayed enabled")
	}
	if err := handler.SettingsSet(ID, "enabled", action.Bool, 0, true); err != nil {
		t.Fatal(err)
	}

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() { _ = svc.Stop() })
	hw.WaitSubscribed(t)
	hw.WaitSubscribedRealtime(t)
	waitRunning(t, svc, true)
	readForwarded(t, hw, port, fixed)

	if err := svc.Restart(); err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	hw.WaitSubscribed(t)
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d, want 1", svc.GetStatus().GetRestarts())
	}
	readForwarded(t, hw, port, fixed.Add(time.Second))

	if err := svc.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	waitRunning(t, svc, false)
}

func TestStartRejectsInvalidListenAddress(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	svc := New(testsupport.NewHardware(), handler, timesource.New(nil))
	if err := svc.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	svc.listenHost = "127.0.0.1"
	svc.listenPort = -1
	if err := svc.Start(); err == nil {
		t.Fatal("Start() accepted an invalid listen port")
	}
}

func TestStartSubscribeFailureStopsCleanly(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	port := freePort(t)
	if err := handler.SettingsSet(ID, "listen_host", action.String, 0, "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if err := handler.SettingsSet(ID, "listen_port", action.Int, 0, int64(port)); err != nil {
		t.Fatal(err)
	}
	hw := testsupport.NewHardware()
	hw.SubscribeRealtimeErr = errors.New("realtime unavailable")
	svc := New(hw, handler, timesource.New(nil))
	if err := svc.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := svc.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	waitRunning(t, svc, false)
	if err := svc.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
}

func readForwarded(t *testing.T, hw *testsupport.Hardware, port int, ts time.Time) {
	t.Helper()
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial forwarder: %v", err)
	}
	defer conn.Close()

	event := explorer.Event{
		Timestamp:  ts,
		SampleRate: 50,
		ChannelData: []explorer.ChannelData{{
			ChannelCode: "EHZ",
			ChannelId:   3,
			Data:        []int32{9, 8, 7},
		}},
	}
	deadline := time.Now().Add(2 * time.Second)
	buf := make([]byte, 512)
	for time.Now().Before(deadline) {
		hw.Emit(event)
		_ = conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		n, readErr := conn.Read(buf)
		if n > 0 {
			if !strings.Contains(string(buf[:n]), "EHZ") {
				t.Fatalf("forwarded payload = %q", buf[:n])
			}
			return
		}
		if readErr != nil && !isTimeout(readErr) {
			continue
		}
	}
	t.Fatal("timed out waiting for forwarded data")
}

func isTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	return port
}

func waitRunning(t *testing.T, svc *ForwarderServiceImpl, running bool) {
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
