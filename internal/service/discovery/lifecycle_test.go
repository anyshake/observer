package mdns_discovery

import (
	"net"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	service_forwarder "github.com/anyshake/observer/internal/service/forwarder"
	service_seedlink "github.com/anyshake/observer/internal/service/seedlink"
	service_winston "github.com/anyshake/observer/internal/service/winston"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestRegistrationBranches(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	clock := timesource.New(func() time.Time { return fixed })
	hw := testsupport.NewHardware()
	svc := New("not-a-host", map[string]service.IService{}, handler, clock)
	if svc.GetName() != "mDNS Discovery Service" || svc.GetDescription() == "" || svc.GetStatus() == nil {
		t.Fatal("service metadata is incomplete")
	}
	if svc.IsEnabled() {
		t.Fatal("enabled before initialization")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}

	svc.syncRegistrations()
	svc.webServerAddr = "127.0.0.1:not-a-port"
	svc.syncRegistrations()

	seedlinkSvc := service_seedlink.New(hw, handler, clock)
	if err := seedlinkSvc.Init(); err != nil {
		t.Fatal(err)
	}
	winstonSvc := service_winston.New(hw, handler, clock)
	if err := winstonSvc.Init(); err != nil {
		t.Fatal(err)
	}
	forwarderSvc := service_forwarder.New(hw, handler, clock)
	if err := forwarderSvc.Init(); err != nil {
		t.Fatal(err)
	}
	svc.serviceMap = map[string]service.IService{
		service_seedlink.ID:  svc,
		service_winston.ID:   winstonSvc,
		service_forwarder.ID: forwarderSvc,
	}
	svc.registerSeedLink = true
	svc.registerWinston = true
	svc.registerForwarder = true
	svc.syncRegistrations()

	svc.serviceMap[service_seedlink.ID] = seedlinkSvc
	svc.seedlinkReg = &mdnsRegistration{port: 1}
	svc.winstonReg = &mdnsRegistration{port: 1}
	svc.forwarderReg = &mdnsRegistration{port: 1}
	svc.registerSeedLink = false
	svc.registerWinston = false
	svc.registerForwarder = false
	svc.syncRegistrations()
	if svc.seedlinkReg != nil || svc.winstonReg != nil || svc.forwarderReg != nil {
		t.Fatal("disabled registrations were kept")
	}

	svc.registerSeedLink = true
	svc.serviceMap = map[string]service.IService{}
	svc.syncRegistrations()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()
	for _, constraint := range seedlinkSvc.GetConfigConstraint() {
		switch constraint.GetKey() {
		case "listen_host":
			if err := constraint.Set(handler, "127.0.0.1"); err != nil {
				t.Fatal(err)
			}
		case "listen_port":
			if err := constraint.Set(handler, int64(port)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := seedlinkSvc.Init(); err != nil {
		t.Fatal(err)
	}
	if err := seedlinkSvc.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = seedlinkSvc.Stop() })
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && !seedlinkSvc.GetStatus().GetIsRunning() {
		time.Sleep(5 * time.Millisecond)
	}
	svc.serviceMap = map[string]service.IService{service_seedlink.ID: seedlinkSvc}
	svc.syncRegistrations()
}

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New("not-a-host", nil, handler, timesource.New(func() time.Time { return fixed }))
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "enabled", true); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled() {
		t.Fatal("service stayed disabled")
	}
	if err := svc.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.Stop() })
	waitRunning(t, svc, true)
	if err := svc.Restart(); err != nil {
		t.Fatal(err)
	}
	waitRunning(t, svc, true)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d", svc.GetStatus().GetRestarts())
	}
	if err := svc.Stop(); err != nil {
		t.Fatal(err)
	}
	waitRunning(t, svc, false)
}

func setConstraint(t *testing.T, handler *action.Handler, svc *DiscoveryServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() == key {
			return constraint.Set(handler, value)
		}
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}

func waitRunning(t *testing.T, svc *DiscoveryServiceImpl, running bool) {
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
