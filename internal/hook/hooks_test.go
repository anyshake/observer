package hook_test

import (
	"context"
	"errors"
	"testing"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/hook/cleaner/close_database"
	"github.com/anyshake/observer/internal/hook/cleaner/close_explorer"
	"github.com/anyshake/observer/internal/hook/startup/migrate_database"
	"github.com/anyshake/observer/internal/hook/startup/setup_admin"
	"github.com/anyshake/observer/internal/hook/startup/setup_station"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/metadata"
)

func TestStartupAndCleanerHooks(t *testing.T) {
	logger.Init()
	obj, handler := testsupport.OpenDAO(t)

	admin := &setup_admin.SetupAdminStartupImpl{ActionHandler: handler}
	if admin.GetName() == "" {
		t.Fatal("empty admin hook name")
	}
	if err := admin.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := admin.Execute(); err != nil {
		t.Fatal(err)
	}
	user, err := handler.SysUserGetByUsername(model.DEFAULT_USERNAME)
	if err != nil {
		t.Fatal(err)
	}
	user.HashedPassword, err = user.GetHashedPassword("Changed@123")
	if err != nil {
		t.Fatal(err)
	}
	if err := handler.SysUserUpdate(user.UserId, user); err != nil {
		t.Fatal(err)
	}
	if err := admin.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := (&setup_admin.SetupAdminStartupImpl{ActionHandler: action.NewHandler(nil)}).Execute(); err == nil {
		t.Fatal("admin hook without a database succeeded")
	}

	station := &setup_station.SetupStationStartupImpl{
		ActionHandler:            handler,
		StationConfigConstraints: config.NewStationConstraints(),
	}
	if station.GetName() == "" {
		t.Fatal("empty station hook name")
	}
	if err := station.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := (&setup_station.SetupStationStartupImpl{
		ActionHandler:            action.NewHandler(nil),
		StationConfigConstraints: config.NewStationConstraints(),
	}).Execute(); err == nil {
		t.Fatal("station hook without a database succeeded")
	}

	migration := &migrate_database.MigrateDatabaseStartupImpl{ActionHandler: handler}
	if migration.GetName() == "" {
		t.Fatal("empty migration hook name")
	}
	if err := migration.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := migration.Execute(); err != nil {
		t.Fatal(err)
	}

	databaseCleaner := &close_database.CloseDatabaseCleanerImpl{}
	if databaseCleaner.GetName() == "" {
		t.Fatal("empty database cleaner name")
	}
	if err := databaseCleaner.Execute(); err != nil {
		t.Fatal(err)
	}
	databaseCleaner.DAO = obj
	if err := databaseCleaner.Execute(); err != nil {
		t.Fatal(err)
	}

	explorerCleaner := &close_explorer.CloseExplorerCleanerImpl{}
	if explorerCleaner.GetName() == "" {
		t.Fatal("empty explorer cleaner name")
	}
	if err := explorerCleaner.Execute(); err != nil {
		t.Fatal(err)
	}
	explorerCleaner.HardwareDev = fakeHardware{}
	if err := explorerCleaner.Execute(); err != nil {
		t.Fatal(err)
	}
	explorerCleaner.HardwareDev = fakeHardware{closeErr: errors.New("busy")}
	if err := explorerCleaner.Execute(); err == nil {
		t.Fatal("hardware close error was ignored")
	}
}

type fakeHardware struct {
	closeErr error
}

func (f fakeHardware) Open(context.Context) (context.Context, context.CancelFunc, error) {
	return context.Background(), func() {}, nil
}
func (f fakeHardware) Close() error { return f.closeErr }
func (fakeHardware) Flush() error   { return nil }
func (fakeHardware) Subscribe(string, explorer.EventHandler) error {
	return nil
}
func (fakeHardware) Unsubscribe(string) error { return nil }
func (fakeHardware) SubscribeRealtime(string, explorer.EventHandler) error {
	return nil
}
func (fakeHardware) UnsubscribeRealtime(string) error { return nil }
func (fakeHardware) GetConfig() explorer.DeviceConfig { return explorer.DeviceConfig{} }
func (fakeHardware) GetStatus() explorer.DeviceStatus { return explorer.DeviceStatus{} }
func (fakeHardware) GetCoordinates(bool) (float64, float64, float64, error) {
	return 0, 0, 0, nil
}
func (fakeHardware) GetTemperature() (float64, error) { return 0, nil }
func (fakeHardware) GetDeviceId() string              { return "device" }
func (fakeHardware) GetMetadata(string, string, string, string, string, string, string, bool) (*metadata.Render, error) {
	return &metadata.Render{}, nil
}
