package graph_resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/server/middleware/auth_jwt"
	graph_model "github.com/anyshake/observer/internal/server/router/graph/model"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/internal/service/helicorder"
	"github.com/anyshake/observer/internal/service/miniseed"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/jobtracker"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/seisevent"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"
	"github.com/gin-gonic/gin"
)

const testPassword = "GoodPass1!"

func TestResolverHelpers(t *testing.T) {
	env := newEnv(t)
	if env.resolver.getCurrentUserId(context.Background()) != "" {
		t.Fatal("missing user status returned an id")
	}
	if env.resolver.getCurrentUserId(withStatus("nope")) != "" {
		t.Fatal("wrong user status type returned an id")
	}
	if env.resolver.getCurrentUserId(withStatus(map[string]any{})) != "" {
		t.Fatal("missing user id returned an id")
	}
	if got := env.resolver.getCurrentUserId(adminContext("user-1")); got != "user-1" {
		t.Fatalf("user id = %s", got)
	}
	if env.resolver.checkIsAdmin(context.Background()) || env.resolver.checkIsAdmin(withStatus(42)) || env.resolver.checkIsAdmin(withStatus(map[string]any{})) {
		t.Fatal("non-admin context was treated as admin")
	}
	if !env.resolver.checkIsAdmin(adminContext("user-1")) || env.resolver.checkIsAdmin(userContext("user-1", false)) {
		t.Fatal("admin flag was parsed incorrectly")
	}

	if _, err := env.resolver.getMiniSeedService(); err == nil {
		t.Fatal("missing miniseed service accepted")
	}
	env.resolver.ServiceMap["other"] = &fakeService{}
	env.resolver.ServiceMap[miniseed.ID] = &fakeService{}
	if _, err := env.resolver.getMiniSeedService(); err == nil {
		t.Fatal("unexpected miniseed type accepted")
	}
	env.resolver.ServiceMap[helicorder.ID] = &fakeService{}
	if _, err := env.resolver.getHelicorderService(); err == nil {
		t.Fatal("unexpected helicorder type accepted")
	}
	mini := env.prepareMiniSeed(t)
	heli := env.prepareHelicorder(t)
	env.resolver.ServiceMap[miniseed.ID] = mini
	env.resolver.ServiceMap[helicorder.ID] = heli
	if got, err := env.resolver.getMiniSeedService(); err != nil || got != mini {
		t.Fatalf("miniseed service = %v, %v", got, err)
	}
	if got, err := env.resolver.getHelicorderService(); err != nil || got != heli {
		t.Fatalf("helicorder service = %v, %v", got, err)
	}

	now := time.Date(2026, 10, 8, 1, 2, 3, 0, time.UTC)
	finished := now.Add(time.Second)
	jobErr := errors.New("boom")
	for _, tc := range []struct {
		status jobtracker.JobStatus
		want   graph_model.JobStatus
	}{
		{jobtracker.JobStatusIdle, graph_model.JobStatusIdle},
		{jobtracker.JobStatusRunning, graph_model.JobStatusRunning},
		{jobtracker.JobStatusSucceeded, graph_model.JobStatusSucceeded},
		{jobtracker.JobStatusFailed, graph_model.JobStatusFailed},
		{jobtracker.JobStatus("NOPE"), graph_model.JobStatusIdle},
	} {
		got := env.resolver.toPurgeDataJobResponse(&jobtracker.Job{ID: "9", Kind: "kind", Status: tc.status})
		if got.Status != tc.want || got.StartedAt != nil || got.Error != nil {
			t.Fatalf("status %s mapped to %+v", tc.status, got)
		}
	}
	full := env.resolver.toPurgeDataJobResponse(&jobtracker.Job{
		ID: "9", Kind: "kind", Status: jobtracker.JobStatusFailed,
		StartedAt: &now, FinishedAt: &finished, Error: jobErr,
	})
	if full.StartedAt == nil || *full.StartedAt != now.UnixMilli() || full.FinishedAt == nil || full.Error == nil || *full.Error != "boom" {
		t.Fatalf("full job = %+v", full)
	}

	if LoadOrCreatePurgeDataJob(nil) != nil {
		t.Fatal("nil resolver created a tracker")
	}
	first := LoadOrCreatePurgeDataJob(env.resolver)
	second := LoadOrCreatePurgeDataJob(env.resolver)
	if first == nil || first != second || first.Get().Kind != "data_purge_job" {
		t.Fatal("purge tracker was not reused")
	}
	if env.resolver.Mutation() == nil || env.resolver.Query() == nil {
		t.Fatal("resolver factories returned nil")
	}
}

func TestPermissionChecks(t *testing.T) {
	env := newEnv(t)
	mut := env.mutation()
	query := env.query()
	password := testPassword
	serviceID := "svc"
	for _, ctx := range []context.Context{context.Background(), userContext("user-1", false)} {
		if _, err := mut.CreateSysUser(ctx, "abc", testPassword, false); err == nil {
			t.Fatal("create user allowed")
		}
		if _, err := mut.RemoveSysUser(ctx, "id"); err == nil {
			t.Fatal("remove user allowed")
		}
		if _, err := mut.UpdateSysUser(ctx, "id", "abc", &password, true); err == nil {
			t.Fatal("update user allowed")
		}
		if _, err := mut.PurgeSeisRecords(ctx); err == nil {
			t.Fatal("purge records allowed")
		}
		if _, err := mut.PurgeSeisRecordsByDate(ctx, 1, 2); err == nil {
			t.Fatal("purge records by date allowed")
		}
		if _, err := mut.PurgeMiniSeedFiles(ctx); err == nil {
			t.Fatal("purge miniseed allowed")
		}
		if _, err := mut.PurgeMiniSeedFilesByDate(ctx, 1, 2); err == nil {
			t.Fatal("purge miniseed by date allowed")
		}
		if _, err := mut.PurgeHelicorderFiles(ctx); err == nil {
			t.Fatal("purge helicorder allowed")
		}
		if _, err := mut.PurgeHelicorderFilesByDate(ctx, 1, 2); err == nil {
			t.Fatal("purge helicorder by date allowed")
		}
		if _, err := mut.UpdateStationConfig(ctx, "station_code", "AB"); err == nil {
			t.Fatal("station update allowed")
		}
		if _, err := mut.RestoreStationConfig(ctx); err == nil {
			t.Fatal("station restore allowed")
		}
		if _, err := mut.ImportGlobalConfig(ctx, "{}"); err == nil {
			t.Fatal("import allowed")
		}
		if _, err := mut.RestartApplication(ctx); err == nil {
			t.Fatal("restart allowed")
		}
		if _, err := mut.StopService(ctx, serviceID); err == nil {
			t.Fatal("stop allowed")
		}
		if _, err := mut.StartService(ctx, serviceID); err == nil {
			t.Fatal("start allowed")
		}
		if _, err := mut.RestartService(ctx, serviceID); err == nil {
			t.Fatal("service restart allowed")
		}
		if _, err := mut.RestoreServiceConfig(ctx, &serviceID); err == nil {
			t.Fatal("service restore allowed")
		}
		if _, err := mut.UpdateServiceConfig(ctx, serviceID, "mode", "on"); err == nil {
			t.Fatal("service update allowed")
		}
		if _, err := query.GetServiceConfigConstraint(ctx); err == nil {
			t.Fatal("service constraints allowed")
		}
		if _, err := query.GetStationConfigConstraint(ctx); err == nil {
			t.Fatal("station constraints allowed")
		}
		if _, err := query.GetSysUsers(ctx); err == nil {
			t.Fatal("user list allowed")
		}
		if _, err := query.IsGenuineProduct(ctx); err == nil {
			t.Fatal("genuine check allowed")
		}
		if _, err := query.GetApplicationLogs(ctx); err == nil {
			t.Fatal("logs allowed")
		}
		if _, err := query.ExportGlobalConfig(ctx); err == nil {
			t.Fatal("export allowed")
		}
		if _, err := query.GetUpgradeStatus(ctx); err == nil {
			t.Fatal("upgrade status allowed")
		}
		if _, err := query.GetCleanupStatus(ctx); err == nil {
			t.Fatal("cleanup status allowed")
		}
	}
}

func TestSysUsers(t *testing.T) {
	env := newEnv(t)
	mut := env.mutation()
	query := env.query()
	admin := adminContext("bootstrap")
	if _, err := mut.CreateSysUser(admin, "ab", testPassword, true); err == nil {
		t.Fatal("short username accepted")
	}
	if _, err := mut.CreateSysUser(admin, "valid_user", "short", true); err == nil {
		t.Fatal("weak password accepted")
	}
	userID, err := mut.CreateSysUser(admin, "valid_user", testPassword, true)
	if err != nil || userID == "" {
		t.Fatal(err)
	}
	if _, err := mut.CreateSysUser(admin, "valid_user", testPassword, false); err == nil {
		t.Fatal("duplicate user accepted")
	}

	current, err := query.GetCurrentUser(adminContext(userID))
	if err != nil || current.Username != "valid_user" || !current.Admin || current.UserID != userID {
		t.Fatalf("current = %+v, %v", current, err)
	}
	if _, err := query.GetCurrentUser(adminContext("missing")); err == nil {
		t.Fatal("missing current user accepted")
	}
	users, err := query.GetSysUsers(admin)
	if err != nil || len(users) != 1 || users[0].UserID != userID {
		t.Fatalf("users = %+v, %v", users, err)
	}

	if _, err := mut.UpdateSysUser(admin, userID, "xy", nil, true); err == nil {
		t.Fatal("short renamed username accepted")
	}
	bad := "nocapitals!"
	if _, err := mut.UpdateSysUser(admin, userID, "valid_user", &bad, true); err == nil {
		t.Fatal("weak updated password accepted")
	}
	if _, err := mut.UpdateSysUser(adminContext(userID), userID, "valid_user", nil, false); err == nil {
		t.Fatal("admin downgraded itself")
	}
	next := "BetterPass1!"
	if _, err := mut.UpdateSysUser(admin, userID, "renamed_user", &next, true); err != nil {
		t.Fatal(err)
	}
	if _, err := env.handler.SysUserLogin("renamed_user", next, "agent", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.UpdateSysUser(admin, userID, "renamed_user", nil, true); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.RemoveSysUser(adminContext(userID), userID); err == nil {
		t.Fatal("current user was removed")
	}
	otherID, err := mut.CreateSysUser(admin, "other_user", testPassword, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mut.RemoveSysUser(admin, otherID); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.RemoveSysUser(admin, otherID); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.UpdateSysUser(admin, "missing", "valid_name", nil, false); err == nil {
		t.Fatal("missing user was updated")
	}
	env.resolver.ActionHandler = action.NewHandler(nil)
	if _, err := mut.RemoveSysUser(admin, userID); err == nil {
		t.Fatal("nil database removed a user")
	}
	if _, err := query.GetSysUsers(admin); err == nil {
		t.Fatal("nil database listed users")
	}
}

func TestStationAndServiceConfig(t *testing.T) {
	env := newEnv(t)
	mut := env.mutation()
	query := env.query()
	admin := adminContext("admin")

	if _, err := mut.UpdateStationConfig(admin, "missing", "x"); err == nil {
		t.Fatal("missing station key accepted")
	}
	if _, err := mut.UpdateStationConfig(admin, "station_code", ""); err == nil {
		t.Fatal("empty station code accepted")
	}
	if _, err := mut.UpdateStationConfig(admin, "station_code", "ab"); err != nil {
		t.Fatal(err)
	}
	cfg, err := query.GetStationConfig(admin)
	if err != nil || cfg["station_code"] != "AB" {
		t.Fatalf("station config = %#v, %v", cfg, err)
	}
	if _, err := mut.RestoreStationConfig(admin); err != nil {
		t.Fatal(err)
	}
	cfg, err = query.GetStationConfig(context.Background())
	if err != nil || cfg["station_code"] != "SHAKE" {
		t.Fatalf("restored config = %#v, %v", cfg, err)
	}

	constraints, err := query.GetStationConfigConstraint(admin)
	if err != nil || len(constraints) != len(env.resolver.StationConfigConstraints) {
		t.Fatalf("constraints = %d, %v", len(constraints), err)
	}
	broken := &fakeConstraint{namespace: "global_station", key: "broken", name: "Broken", description: "bad", typ: action.String, getErr: errors.New("unreadable")}
	env.resolver.StationConfigConstraints = append(env.resolver.StationConfigConstraints, broken)
	if _, err := query.GetStationConfig(admin); err == nil {
		t.Fatal("broken station config accepted")
	}
	if _, err := query.GetStationConfigConstraint(admin); err == nil {
		t.Fatal("broken station constraint accepted")
	}
	env.resolver.StationConfigConstraints = []config.IConstraint{&fakeConstraint{namespace: "global_station", key: "broken", name: "Broken", description: "bad", typ: action.String, restoreErr: errors.New("restore failed")}}
	if _, err := mut.RestoreStationConfig(admin); err == nil {
		t.Fatal("broken restore accepted")
	}

	mode := &fakeConstraint{namespace: "custom", key: "mode", name: "Mode", description: "mode", typ: action.String, required: true, version: 0, value: "off", options: map[string]any{"On": "on"}}
	env.resolver.ServiceMap["custom"] = &fakeService{name: "Custom", description: "custom service", constraints: []config.IConstraint{mode}}
	if _, err := mut.UpdateServiceConfig(admin, "missing", "mode", "on"); err == nil {
		t.Fatal("missing service accepted")
	}
	if _, err := mut.UpdateServiceConfig(admin, "custom", "missing", "on"); err == nil {
		t.Fatal("missing service key accepted")
	}
	mode.setErr = errors.New("set failed")
	if _, err := mut.UpdateServiceConfig(admin, "custom", "mode", "on"); err == nil {
		t.Fatal("service set error ignored")
	}
	mode.setErr = nil
	if _, err := mut.UpdateServiceConfig(admin, "custom", "mode", "on"); err != nil || mode.value != "on" {
		t.Fatalf("mode = %#v, %v", mode.value, err)
	}
	serviceID := "custom"
	if _, err := mut.RestoreServiceConfig(admin, &serviceID); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.RestoreServiceConfig(admin, nil); err != nil {
		t.Fatal(err)
	}
	missingID := "missing"
	if _, err := mut.RestoreServiceConfig(admin, &missingID); err == nil {
		t.Fatal("missing service restore accepted")
	}
	mode.restoreErr = errors.New("restore failed")
	if _, err := mut.RestoreServiceConfig(admin, &serviceID); err == nil {
		t.Fatal("service restore error ignored")
	}
	if _, err := mut.RestoreServiceConfig(admin, nil); err == nil {
		t.Fatal("restore of every service ignored an error")
	}
	mode.restoreErr = nil

	listed, err := query.GetServiceConfigConstraint(admin)
	if err != nil || len(listed) != 1 || listed[0].ServiceID != "custom" || len(listed[0].Constraints) != 1 || listed[0].Constraints[0].CurrentValue != "on" {
		t.Fatalf("service constraints = %+v, %v", listed, err)
	}
	mode.getErr = errors.New("unreadable")
	if _, err := query.GetServiceConfigConstraint(admin); err == nil {
		t.Fatal("broken service constraint accepted")
	}
	mode.getErr = nil
}

func TestImportAndExportConfig(t *testing.T) {
	env := newEnv(t)
	mut := env.mutation()
	query := env.query()
	admin := adminContext("admin")
	mode := &fakeConstraint{namespace: "custom", key: "mode", name: "Mode", description: "mode", typ: action.String, value: "off"}
	env.resolver.ServiceMap["custom"] = &fakeService{name: "Custom", constraints: []config.IConstraint{mode}}

	if _, err := mut.ImportGlobalConfig(admin, "{"); err == nil {
		t.Fatal("invalid json accepted")
	}
	if _, err := mut.ImportGlobalConfig(admin, "{}"); err == nil {
		t.Fatal("empty config accepted")
	}
	if _, err := mut.ImportGlobalConfig(admin, `{"unused":{}}`); err != nil {
		t.Fatal(err)
	}
	payload := `{
		"global_station": {
			"station_code": {"type": "string", "version": 0, "value": "ab"},
			"skipped": {"type": "int", "version": 0, "value": 1},
			"network_code": {"type": "string", "value": "tw"}
		},
		"custom": {
			"mode": {"type": "string", "version": 0, "value": "on"},
			"absent": {"type": "string", "version": 0, "value": "x"}
		}
	}`
	if _, err := mut.ImportGlobalConfig(admin, payload); err != nil {
		t.Fatal(err)
	}
	if mode.value != "on" {
		t.Fatalf("imported mode = %#v", mode.value)
	}
	cfg, err := query.GetStationConfig(admin)
	if err != nil || cfg["station_code"] != "AB" || cfg["network_code"] != "TW" {
		t.Fatalf("imported station = %#v, %v", cfg, err)
	}

	skipped := &fakeConstraint{namespace: "custom", key: "typed", name: "Typed", description: "typed", typ: action.String, version: 3}
	env.resolver.ServiceMap["custom"] = &fakeService{constraints: []config.IConstraint{mode, skipped}}
	if _, err := mut.ImportGlobalConfig(admin, `{"custom":{"typed":{"type":"int","version":"bad","value":"x"},"mode":{"type":"string","version":0}}}`); err != nil {
		t.Fatal(err)
	}
	if skipped.value != nil {
		t.Fatalf("skipped constraint stored %#v", skipped.value)
	}

	mode.setErr = errors.New("set failed")
	if _, err := mut.ImportGlobalConfig(admin, `{"custom":{"mode":{"type":"string","version":0,"value":"bad"}}}`); err == nil {
		t.Fatal("service import error ignored")
	}
	mode.setErr = nil
	env.resolver.StationConfigConstraints = []config.IConstraint{&fakeConstraint{namespace: config.STATION_NAMESPACE, key: "station_code", typ: action.String, setErr: errors.New("set failed")}}
	if _, err := mut.ImportGlobalConfig(admin, `{"global_station":{"station_code":{"type":"string","version":0,"value":"ZZ"}}}`); err == nil {
		t.Fatal("station import error ignored")
	}
	env.resolver.StationConfigConstraints = []config.IConstraint{
		&fakeConstraint{namespace: config.STATION_NAMESPACE, key: "station_code", typ: action.String},
		&fakeConstraint{namespace: config.STATION_NAMESPACE, key: "network_code", typ: action.String, version: 4},
		&fakeConstraint{namespace: config.STATION_NAMESPACE, key: "location_code", typ: action.String},
	}
	env.resolver.ServiceMap = map[string]service.IService{"custom": &fakeService{constraints: []config.IConstraint{
		&fakeConstraint{namespace: "custom", key: "missing", typ: action.String},
		&fakeConstraint{namespace: "custom", key: "typed", typ: action.String},
		&fakeConstraint{namespace: "custom", key: "ver", typ: action.String, version: 5},
	}}}
	if _, err := mut.ImportGlobalConfig(admin, `{
		"global_station": {
			"station_code": {"type": "int", "version": 0, "value": "AB"},
			"network_code": {"type": "string", "version": "bad", "value": "TW"},
			"location_code": {"type": "string", "version": 0}
		},
		"custom": {
			"typed": {"type": "int", "version": 0, "value": "x"},
			"ver": {"type": "string", "version": "bad", "value": "x"}
		}
	}`); err != nil {
		t.Fatal(err)
	}

	env = newEnv(t)
	env.resolver.ServiceMap["custom"] = &fakeService{constraints: []config.IConstraint{&fakeConstraint{namespace: "custom", key: "mode", name: "Mode", description: "mode", typ: action.String, value: "off"}}}
	exported, err := env.query().ExportGlobalConfig(adminContext("admin"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]map[string]map[string]any
	if err := json.Unmarshal([]byte(exported), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["global_station"]["station_code"]["value"] != "SHAKE" || decoded["custom"]["mode"]["value"] != "off" {
		t.Fatalf("exported = %s", exported)
	}
	env.resolver.StationConfigConstraints = []config.IConstraint{&fakeConstraint{namespace: "global_station", key: "broken", typ: action.String, getErr: errors.New("unreadable")}}
	if _, err := env.query().ExportGlobalConfig(adminContext("admin")); err == nil {
		t.Fatal("broken export accepted")
	}
	env.resolver.StationConfigConstraints = nil
	env.resolver.ServiceMap = map[string]service.IService{"custom": &fakeService{constraints: []config.IConstraint{
		&fakeConstraint{namespace: "custom", key: "mode", typ: action.String, getErr: errors.New("unreadable")},
	}}}
	if _, err := env.query().ExportGlobalConfig(adminContext("admin")); err == nil {
		t.Fatal("broken service export accepted")
	}
	env.resolver.ServiceMap = nil
	env.resolver.StationConfigConstraints = []config.IConstraint{&fakeConstraint{namespace: "global_station", key: "broken", typ: action.String, value: make(chan int)}}
	if _, err := env.query().ExportGlobalConfig(adminContext("admin")); err == nil {
		t.Fatal("unmarshalable export accepted")
	}
}

func TestServiceLifecycleAndAssets(t *testing.T) {
	env := newEnv(t)
	mut := env.mutation()
	query := env.query()
	admin := adminContext("admin")
	started := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	svc := &fakeService{name: "Custom", description: "custom service", started: started, updated: started, stopped: started}
	env.resolver.ServiceMap["custom"] = svc

	if _, err := mut.StopService(admin, "missing"); err == nil {
		t.Fatal("missing stop accepted")
	}
	if _, err := mut.StopService(admin, "custom"); err == nil {
		t.Fatal("stopped service was stopped again")
	}
	svc.running = true
	svc.stopErr = errors.New("busy")
	if _, err := mut.StopService(admin, "custom"); err == nil {
		t.Fatal("stop error ignored")
	}
	svc.stopErr = nil
	if _, err := mut.StopService(admin, "custom"); err != nil || svc.running {
		t.Fatalf("running = %v, %v", svc.running, err)
	}

	svc.running = true
	if _, err := mut.StartService(admin, "custom"); err == nil {
		t.Fatal("running service was started")
	}
	svc.running = false
	if _, err := mut.StartService(admin, "missing"); err == nil {
		t.Fatal("missing start accepted")
	}
	svc.initErr = errors.New("init")
	if _, err := mut.StartService(admin, "custom"); err == nil {
		t.Fatal("init error ignored")
	}
	svc.initErr = nil
	svc.startErr = errors.New("start")
	if _, err := mut.StartService(admin, "custom"); err == nil {
		t.Fatal("start error ignored")
	}
	svc.startErr = nil
	if _, err := mut.StartService(admin, "custom"); err != nil || !svc.running {
		t.Fatalf("started = %v, %v", svc.running, err)
	}

	if _, err := mut.RestartService(admin, "missing"); err == nil {
		t.Fatal("missing restart accepted")
	}
	svc.restartErr = errors.New("restart")
	if _, err := mut.RestartService(admin, "custom"); err == nil {
		t.Fatal("restart error ignored")
	}
	svc.restartErr = nil
	if _, err := mut.RestartService(admin, "custom"); err != nil || svc.restarts != 1 {
		t.Fatalf("restarts = %d, %v", svc.restarts, err)
	}

	status, err := query.GetServiceStatus(context.Background())
	if err != nil || len(status) != 1 || status[0].ServiceID != "custom" || !status[0].IsRunning || status[0].Restarts != 1 {
		t.Fatalf("status = %+v, %v", status, err)
	}

	svc.assets = []service.Asset{{FilePath: "/tmp/a.mseed", FileName: "a.mseed", Size: 4, ModifiedAt: 9}}
	env.resolver.ServiceMap[miniseed.ID] = svc
	files, err := query.GetMiniSeedFiles(context.Background())
	if err != nil || len(files) != 1 || files[0].Namespace != miniseed.ID || files[0].FileName != "a.mseed" {
		t.Fatalf("miniseed files = %+v, %v", files, err)
	}
	svc.assetErr = errors.New("unavailable")
	if _, err := query.GetMiniSeedFiles(context.Background()); err == nil {
		t.Fatal("miniseed asset error ignored")
	}
	delete(env.resolver.ServiceMap, miniseed.ID)
	if _, err := query.GetMiniSeedFiles(context.Background()); err == nil {
		t.Fatal("missing miniseed service accepted")
	}

	svc.assetErr = nil
	svc.assets = []service.Asset{{FilePath: "/tmp/a.png", FileName: "a.png", Size: 8, ModifiedAt: 10}}
	env.resolver.ServiceMap[helicorder.ID] = svc
	images, err := query.GetHelicorderFiles(context.Background())
	if err != nil || len(images) != 1 || images[0].Namespace != helicorder.ID || images[0].Size != 8 {
		t.Fatalf("helicorder files = %+v, %v", images, err)
	}
	svc.assetErr = errors.New("unavailable")
	if _, err := query.GetHelicorderFiles(context.Background()); err == nil {
		t.Fatal("helicorder asset error ignored")
	}
	delete(env.resolver.ServiceMap, helicorder.ID)
	if _, err := query.GetHelicorderFiles(context.Background()); err == nil {
		t.Fatal("missing helicorder service accepted")
	}
}

func TestPurgeJobs(t *testing.T) {
	env := newEnv(t)
	mut := env.mutation()
	query := env.query()
	admin := adminContext("admin")

	if _, err := (&mutationResolver{}).PurgeSeisRecords(admin); err == nil {
		t.Fatal("nil resolver created a purge job")
	}
	if _, err := (&mutationResolver{}).PurgeSeisRecordsByDate(admin, 1, 2); err == nil {
		t.Fatal("nil resolver created a dated purge job")
	}
	if _, err := (&queryResolver{}).GetCleanupStatus(admin); err == nil {
		t.Fatal("nil resolver returned cleanup status")
	}

	idle, err := query.GetCleanupStatus(admin)
	if err != nil || idle.Status != graph_model.JobStatusIdle || idle.Kind != "data_purge_job" {
		t.Fatalf("idle = %+v, %v", idle, err)
	}

	full, err := mut.PurgeSeisRecords(admin)
	if err != nil || full.Kind != "seis_records_full" || full.Status != graph_model.JobStatusRunning {
		t.Fatalf("full purge = %+v, %v", full, err)
	}
	again, err := mut.PurgeSeisRecords(admin)
	if err != nil || again.ID != full.ID {
		t.Fatalf("overlapping purge = %+v, %v", again, err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusSucceeded {
		t.Fatalf("full purge finished as %+v", job)
	}

	failed, err := mut.PurgeSeisRecordsByDate(admin, time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC).UnixMilli(), time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli())
	if err != nil || failed.Kind != "seis_records_by_date" {
		t.Fatalf("invalid date purge = %+v, %v", failed, err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusFailed || job.Error == nil {
		t.Fatalf("invalid date purge finished as %+v", job)
	}
	byDate, err := mut.PurgeSeisRecordsByDate(admin, time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli(), time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC).UnixMilli())
	if err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusSucceeded || job.Kind != byDate.Kind {
		t.Fatalf("date purge finished as %+v", job)
	}
	cleanup, err := query.GetCleanupStatus(admin)
	if err != nil || cleanup.Status != graph_model.JobStatusSucceeded {
		t.Fatalf("cleanup = %+v, %v", cleanup, err)
	}

	env.resolver.ServiceMap = map[string]service.IService{}
	if _, err := mut.PurgeMiniSeedFiles(admin); err == nil {
		t.Fatal("missing miniseed purge accepted")
	}
	env.resolver.ServiceMap[miniseed.ID] = &fakeService{}
	if _, err := mut.PurgeMiniSeedFiles(admin); err == nil {
		t.Fatal("fake miniseed purge accepted")
	}
	miniDir := t.TempDir()
	mini := env.prepareMiniSeedIn(t, miniDir)
	env.resolver.ServiceMap[miniseed.ID] = mini
	if err := os.MkdirAll(miniDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(miniDir, "keep.mseed"), []byte("mseed"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.PurgeMiniSeedFiles(admin); err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusSucceeded {
		t.Fatalf("miniseed purge = %+v", job)
	}
	if _, err := os.Stat(miniDir); !os.IsNotExist(err) {
		t.Fatalf("miniseed directory stat = %v", err)
	}

	miniDir = t.TempDir()
	dayDir := filepath.Join(miniDir, "2026-01-02")
	otherDir := filepath.Join(miniDir, "2026-01-03")
	if err := os.MkdirAll(dayDir, 0o755); err != nil || os.MkdirAll(otherDir, 0o755) != nil {
		t.Fatal(err)
	}
	mini = env.prepareMiniSeedIn(t, miniDir)
	env.resolver.ServiceMap[miniseed.ID] = mini
	start := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC).UnixMilli()
	end := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli()
	if _, err := mut.PurgeMiniSeedFilesByDate(admin, start, end); err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusFailed {
		t.Fatalf("invalid miniseed date purge = %+v", job)
	}
	if _, err := mut.PurgeMiniSeedFilesByDate(admin, time.Date(2026, 1, 2, 8, 0, 0, 0, time.UTC).UnixMilli(), time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()); err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusSucceeded {
		t.Fatalf("miniseed date purge = %+v", job)
	}
	if _, err := os.Stat(dayDir); !os.IsNotExist(err) {
		t.Fatalf("purged day stat = %v", err)
	}
	if _, err := os.Stat(otherDir); err != nil {
		t.Fatal(err)
	}

	if _, err := mut.PurgeHelicorderFiles(admin); err == nil {
		t.Fatal("missing helicorder purge accepted")
	}
	env.resolver.ServiceMap[helicorder.ID] = &fakeService{}
	if _, err := mut.PurgeHelicorderFiles(admin); err == nil {
		t.Fatal("fake helicorder purge accepted")
	}
	heliDir := t.TempDir()
	heli := env.prepareHelicorderIn(t, heliDir)
	env.resolver.ServiceMap[helicorder.ID] = heli
	if err := os.WriteFile(filepath.Join(heliDir, "plot.png"), []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := mut.PurgeHelicorderFiles(admin); err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusSucceeded {
		t.Fatalf("helicorder purge = %+v", job)
	}
	if _, err := os.Stat(heliDir); !os.IsNotExist(err) {
		t.Fatalf("helicorder directory stat = %v", err)
	}

	heliDir = t.TempDir()
	heliDay := filepath.Join(heliDir, "2026-02-02")
	if err := os.MkdirAll(heliDay, 0o755); err != nil {
		t.Fatal(err)
	}
	heli = env.prepareHelicorderIn(t, heliDir)
	env.resolver.ServiceMap[helicorder.ID] = heli
	if _, err := mut.PurgeHelicorderFilesByDate(admin, time.Date(2026, 2, 3, 0, 0, 0, 0, time.UTC).UnixMilli(), time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC).UnixMilli()); err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusFailed {
		t.Fatalf("invalid helicorder date purge = %+v", job)
	}
	if _, err := mut.PurgeHelicorderFilesByDate(admin, time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC).UnixMilli(), time.Date(2026, 2, 2, 1, 0, 0, 0, time.UTC).UnixMilli()); err != nil {
		t.Fatal(err)
	}
	if job := waitJob(t, env.resolver); job.Status != jobtracker.JobStatusSucceeded {
		t.Fatalf("helicorder date purge = %+v", job)
	}
	if _, err := os.Stat(heliDay); !os.IsNotExist(err) {
		t.Fatalf("purged helicorder day stat = %v", err)
	}
}

func TestDeviceQueries(t *testing.T) {
	env := newEnv(t)
	query := env.query()
	admin := adminContext("admin")

	version, err := query.GetSoftwareVersion(context.Background())
	if err != nil || version != "v1.2.3-abc123-1700000000" {
		t.Fatalf("version = %s, %v", version, err)
	}
	env.resolver.CurrentBuild = unibuild.New("", "", "", "0")
	env.resolver.CurrentVersion = semver.New("1", "2", "3", "")
	version, err = query.GetSoftwareVersion(context.Background())
	if err != nil || version != "v1.2.3-<out-of-tree>" {
		t.Fatalf("custom version = %s, %v", version, err)
	}
	now, err := query.GetCurrentTime(context.Background())
	if err != nil || now != env.clock.UnixMilli() {
		t.Fatalf("now = %d, %v", now, err)
	}
	systemStatus, err := query.GetSystemStatus(context.Background())
	if err != nil || systemStatus.Uptime < 0 || systemStatus.CPU < 0 || systemStatus.Memory < 0 || systemStatus.Disk < 0 {
		t.Fatalf("system = %+v, %v", systemStatus, err)
	}

	hidden, err := query.GetDeviceID(userContext("user", false))
	if err != nil || hidden != "********" {
		t.Fatalf("hidden id = %s, %v", hidden, err)
	}
	visible, err := query.GetDeviceID(admin)
	if err != nil || visible != "ABCDEF01" {
		t.Fatalf("device id = %s, %v", visible, err)
	}
	cfg, err := query.GetDeviceConfig(context.Background())
	if err != nil || cfg.SampleRate != 100 || cfg.Protocol != "v2" || len(cfg.ChannelCodes) != 2 || !cfg.GnssEnabled || cfg.Model != "E-C111G" {
		t.Fatalf("device config = %+v, %v", cfg, err)
	}
	status, err := query.GetDeviceStatus(context.Background())
	if err != nil || status.Frames != 3 || status.Errors != 1 || status.Messages != 2 || status.StartedAt != env.clock.UnixMilli() {
		t.Fatalf("device status = %+v, %v", status, err)
	}

	info, err := query.GetDeviceInfo(admin)
	if err != nil || info.Temperature == nil || *info.Temperature != 36.5 || info.Latitude == nil || *info.Latitude != 25 || env.hardware.lastFuzzy == nil || *env.hardware.lastFuzzy {
		t.Fatalf("admin device info = %+v fuzzy %v, %v", info, env.hardware.lastFuzzy, err)
	}
	info, err = query.GetDeviceInfo(userContext("user", false))
	if err != nil || env.hardware.lastFuzzy == nil || !*env.hardware.lastFuzzy {
		t.Fatalf("public device info did not request fuzzy coordinates: %+v %v", info, err)
	}
	env.hardware.tempErr = errors.New("temp")
	env.hardware.coordErr = errors.New("coords")
	info, err = query.GetDeviceInfo(admin)
	if err != nil || info.Temperature != nil || info.Latitude != nil {
		t.Fatalf("partial device info = %+v, %v", info, err)
	}

	seiscomp, err := query.GetStationMetadata(admin, "seiscomp_xml")
	if err != nil || !strings.Contains(seiscomp, "SHAKE") {
		t.Fatalf("seiscomp = %s, %v", seiscomp, err)
	}
	stationXML, err := query.GetStationMetadata(userContext("user", false), "station_xml")
	if err != nil || !strings.Contains(stationXML, "SHAKE") {
		t.Fatalf("stationxml = %s, %v", stationXML, err)
	}
	if _, err := query.GetStationMetadata(admin, "csv"); err == nil {
		t.Fatal("unknown metadata format accepted")
	}
	env.hardware.metaErr = errors.New("offline")
	message, err := query.GetStationMetadata(admin, "seiscomp_xml")
	if err != nil || !strings.Contains(message, "offline") {
		t.Fatalf("metadata error = %s, %v", message, err)
	}

	_, bareHandler := testsupport.OpenDAO(t)
	bare := newEnv(t)
	bare.resolver.ActionHandler = bareHandler
	steps := []config.IConstraint{
		&config.StationAffiliationConfigConstraintImpl{},
		&config.StationDescriptionConfigConstraintImpl{},
		&config.StationCountryConfigConstraintImpl{},
		&config.StationPlaceConfigConstraintImpl{},
		&config.StationStationCodeConfigConstraintImpl{},
		&config.StationNetworkCodeConfigConstraintImpl{},
		&config.StationLocationCodeConfigConstraintImpl{},
	}
	for _, constraint := range steps {
		if _, err := bare.query().GetStationMetadata(admin, "seiscomp_xml"); err == nil {
			t.Fatalf("metadata succeeded before %s was initialized", constraint.GetKey())
		}
		if err := constraint.Init(bareHandler); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := bare.query().GetStationMetadata(admin, "station_xml"); err != nil {
		t.Fatal(err)
	}

	env.hardware.protocol = "v1"
	genuine, err := query.IsGenuineProduct(admin)
	if err != nil || genuine {
		t.Fatalf("v1 genuine = %v, %v", genuine, err)
	}
	env.hardware.protocol = "v2"
	env.hardware.id = "FFFFFFFF"
	genuine, err = query.IsGenuineProduct(admin)
	if err != nil || genuine {
		t.Fatalf("blank genuine = %v, %v", genuine, err)
	}
	env.hardware.id = "00000000"
	genuine, err = query.IsGenuineProduct(admin)
	if err != nil || genuine {
		t.Fatalf("zero genuine = %v, %v", genuine, err)
	}
	env.hardware.id = "ABCDEF01"
	genuine, err = query.IsGenuineProduct(admin)
	if err != nil || !genuine {
		t.Fatalf("real genuine = %v, %v", genuine, err)
	}

	logs, err := query.GetApplicationLogs(admin)
	if err != nil || len(logs) != 2 || logs[0] != "alpha" || logs[1] != "beta" {
		t.Fatalf("logs = %#v, %v", logs, err)
	}
	upgrade, err := query.GetUpgradeStatus(admin)
	if err != nil || upgrade != nil {
		t.Fatalf("upgrade = %+v, %v", upgrade, err)
	}
}

func TestEventsAndRecords(t *testing.T) {
	env := newEnv(t)
	query := env.query()
	source := &fakeSource{
		property: seisevent.DataSourceProperty{ID: "local", Country: "TW", Default: "en", Locales: map[string]string{"en": "Local"}},
		events: []seisevent.Event{{
			Verfied: true, Timestamp: 10, Event: "evt", Region: "Taipei", Depth: 8,
			Latitude: 25, Longitude: 121, Distance: 12,
			Magnitude:  []seisevent.Magnitude{{Type: "M", Value: 4.2}},
			Estimation: seisevent.Estimation{P_Wave: 1.5, S_Wave: 2.5},
		}},
	}
	env.resolver.SeisEventSource = map[string]seisevent.IDataSource{"local": source}
	sources, err := query.GetEventSource(context.Background())
	if err != nil || len(sources) != 1 || sources[0].ID != "local" || sources[0].Locales["en"] != "Local" {
		t.Fatalf("sources = %+v, %v", sources, err)
	}
	if _, err := query.GetEventsBySource(context.Background(), "missing"); err == nil {
		t.Fatal("missing event source accepted")
	}
	events, err := query.GetEventsBySource(adminContext("admin"), "local")
	if err != nil || len(events) != 1 || events[0].EventID != "evt" || events[0].Magnitude["M"] != 4.2 || len(events[0].Estimation) != 2 {
		t.Fatalf("events = %+v, %v", events, err)
	}
	env.hardware.coordErr = errors.New("coords")
	if _, err := query.GetEventsBySource(adminContext("admin"), "local"); err == nil {
		t.Fatal("coordinate error ignored")
	}
	env.hardware.coordErr = nil
	source.err = errors.New("upstream")
	if _, err := query.GetEventsBySource(context.Background(), "local"); err == nil {
		t.Fatal("event source error ignored")
	}

	start := time.Date(2026, 10, 8, 1, 2, 3, 0, time.UTC)
	insertSeis(t, env.db, []model.SeisRecord{encodedRecord(t, start, 100, []explorer.ChannelData{{
		ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: []int32{1, -2, 3},
	}})})
	records, err := query.GetSeisRecordsByTime(context.Background(), start.Add(-time.Millisecond).UnixMilli(), start.Add(time.Second).UnixMilli())
	if err != nil || len(records) != 1 || records[0].SampleRate != 100 || records[0].ChannelData[0].ChannelCode != "EHZ" || records[0].ChannelData[0].Data[1] != -2 {
		t.Fatalf("records = %+v, %v", records, err)
	}
	if _, err := query.GetSeisRecordsByTime(context.Background(), start.Add(time.Hour).UnixMilli(), start.UnixMilli()); err == nil {
		t.Fatal("reversed record query accepted")
	}
	empty, err := query.GetSeisRecordsByTime(context.Background(), start.Add(time.Hour).UnixMilli(), start.Add(time.Hour).Add(time.Second).UnixMilli())
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty records = %+v, %v", empty, err)
	}
	insertSeis(t, env.db, []model.SeisRecord{{RecordTime: start.Add(2 * time.Second).UnixMilli(), SampleRate: 100, ChannelData: []byte{1}}})
	if _, err := query.GetSeisRecordsByTime(context.Background(), start.UnixMilli(), start.Add(2*time.Second).UnixMilli()); err == nil {
		t.Fatal("corrupt record accepted")
	}
}

func TestRestartApplication(t *testing.T) {
	env := newEnv(t)
	ok, err := env.mutation().RestartApplication(adminContext("admin"))
	if err != nil || !ok {
		t.Fatal(err)
	}
	select {
	case <-env.resolver.RestartChan:
	case <-time.After(5 * time.Second):
		t.Fatal("restart signal was not sent")
	}
}

type graphEnv struct {
	t        *testing.T
	db       *dao.DAO
	handler  *action.Handler
	resolver *Resolver
	hardware *fakeHardware
	clock    time.Time
}

func newEnv(t *testing.T) *graphEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, handler := testsupport.OpenDAO(t)
	sqlDB, err := db.Database.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	for _, constraint := range config.NewStationConstraints() {
		if err := constraint.Init(handler); err != nil {
			t.Fatal(err)
		}
	}
	clock := time.Date(2026, 10, 8, 1, 2, 3, 0, time.UTC)
	meta, err := metadata.New("E-C111G", metadata.Options{
		StartTime: clock, SampleRate: 100, Latitude: 25, Longitude: 121, Elevation: 10,
		NetworkCode: "AS", StationCode: "SHAKE", LocationCode: "00", ChannelCodes: []string{"EHZ", "EHE"},
		StationPlace: "Lab", StationCountry: "TW", StationAffiliation: "Test", StationDescription: "Unit",
	})
	if err != nil {
		t.Fatal(err)
	}
	logs := ringbuf.New[string](4)
	logs.Push("alpha", "beta")
	hw := &fakeHardware{
		id: "ABCDEF01", sampleRate: 100, interval: 50 * time.Millisecond, channels: []string{"EHZ", "EHE"},
		gnss: true, model: "E-C111G", protocol: "v2", started: clock, updated: clock,
		frames: 3, errors: 1, messages: 2, temp: 36.5, lat: 25, lon: 121, elv: 10, meta: meta,
	}
	resolver := &Resolver{
		RestartChan:              make(chan struct{}, 1),
		CurrentVersion:           semver.New("1", "2", "3", ""),
		CurrentBuild:             unibuild.New("test-toolchain", "stable", "abc123", "1700000000"),
		HardwareDev:              hw,
		TimeSource:               timesource.New(func() time.Time { return clock }),
		ActionHandler:            handler,
		LogBuffer:                logs,
		StationConfigConstraints: config.NewStationConstraints(),
		ServiceMap:               map[string]service.IService{},
		SeisEventSource:          map[string]seisevent.IDataSource{},
	}
	return &graphEnv{t: t, db: db, handler: handler, resolver: resolver, hardware: hw, clock: clock}
}

func (e *graphEnv) mutation() *mutationResolver { return &mutationResolver{e.resolver} }
func (e *graphEnv) query() *queryResolver       { return &queryResolver{e.resolver} }

func (e *graphEnv) prepareMiniSeed(t *testing.T) *miniseed.MiniSeedServiceImpl {
	t.Helper()
	return e.prepareMiniSeedIn(t, t.TempDir())
}

func (e *graphEnv) prepareMiniSeedIn(t *testing.T, dir string) *miniseed.MiniSeedServiceImpl {
	t.Helper()
	if err := e.handler.SettingsSet(miniseed.ID, "file_path", action.String, 0, dir); err != nil {
		t.Fatal(err)
	}
	svc := miniseed.New(e.hardware, e.handler, e.resolver.TimeSource)
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	return svc
}

func (e *graphEnv) prepareHelicorder(t *testing.T) *helicorder.HelicorderServiceImpl {
	t.Helper()
	return e.prepareHelicorderIn(t, t.TempDir())
}

func (e *graphEnv) prepareHelicorderIn(t *testing.T, dir string) *helicorder.HelicorderServiceImpl {
	t.Helper()
	if err := e.handler.SettingsSet(helicorder.ID, "file_path", action.String, 0, dir); err != nil {
		t.Fatal(err)
	}
	if err := e.handler.SettingsSet(helicorder.ID, "cache_storage", action.String, 0, helicorder.CACHE_STORAGE_DISABLED); err != nil {
		t.Fatal(err)
	}
	svc := helicorder.New(e.hardware, e.handler, e.resolver.TimeSource)
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	return svc
}

func adminContext(userID string) context.Context { return userContext(userID, true) }

func userContext(userID string, admin bool) context.Context {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = httptest.NewRequest(http.MethodPost, "/api/graphql", nil)
	ginCtx.Set(auth_jwt.IsAdminKey, admin)
	ginCtx.Set(auth_jwt.UserIdKey, userID)
	return context.WithValue(ginCtx.Request.Context(), ContextKey("user_status"), map[string]any{
		auth_jwt.IsAdminKey: ginCtx.GetBool(auth_jwt.IsAdminKey),
		auth_jwt.UserIdKey:  ginCtx.GetString(auth_jwt.UserIdKey),
	})
}

func withStatus(status any) context.Context {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = httptest.NewRequest(http.MethodPost, "/api/graphql", nil)
	return context.WithValue(ginCtx.Request.Context(), ContextKey("user_status"), status)
}

func waitJob(t *testing.T, resolver *Resolver) *jobtracker.Job {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	var last *jobtracker.Job
	for time.Now().Before(deadline) {
		last = resolver.DataPurgeJob.Get()
		if last.Status == jobtracker.JobStatusSucceeded || last.Status == jobtracker.JobStatusFailed {
			return last
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("job did not finish: %+v", last)
	return nil
}

func insertSeis(t *testing.T, db *dao.DAO, records []model.SeisRecord) {
	t.Helper()
	seen := map[string]struct{}{}
	for _, record := range records {
		day := time.UnixMilli(record.RecordTime).UTC().YearDay()
		table := fmt.Sprintf("%sseis_records_%d", db.GetPrefix(), day%model.SEIS_RECORD_SHARDS)
		if _, ok := seen[table]; ok {
			continue
		}
		if err := db.Database.Table(table).AutoMigrate(&model.SeisRecord{}); err != nil {
			t.Fatal(err)
		}
		seen[table] = struct{}{}
	}
	if err := action.NewHandler(db).SeisRecordsCreate(records...); err != nil {
		t.Fatal(err)
	}
}

func encodedRecord(t *testing.T, when time.Time, rate int, channels []explorer.ChannelData) model.SeisRecord {
	t.Helper()
	var record model.SeisRecord
	if err := record.Encode(when, rate, channels); err != nil {
		t.Fatal(err)
	}
	return record
}

type fakeHardware struct {
	id                         string
	sampleRate                 int
	interval                   time.Duration
	channels                   []string
	gnss                       bool
	model, protocol            string
	started, updated           time.Time
	frames, errors, messages   int64
	temp, lat, lon, elv        float64
	tempErr, coordErr, metaErr error
	lastFuzzy                  *bool
	meta                       *metadata.Render
}

func (f *fakeHardware) Open(context.Context) (context.Context, context.CancelFunc, error) {
	return context.Background(), func() {}, nil
}
func (fakeHardware) Close() error                                  { return nil }
func (fakeHardware) Flush() error                                  { return nil }
func (fakeHardware) Subscribe(string, explorer.EventHandler) error { return nil }
func (fakeHardware) Unsubscribe(string) error                      { return nil }
func (fakeHardware) SubscribeRealtime(string, explorer.EventHandler) error {
	return nil
}
func (fakeHardware) UnsubscribeRealtime(string) error { return nil }
func (f *fakeHardware) GetConfig() explorer.DeviceConfig {
	cfg := explorer.DeviceConfig{}
	cfg.SetPacketInterval(f.interval)
	cfg.SetSampleRate(f.sampleRate)
	cfg.SetChannelCodes(append([]string(nil), f.channels...))
	cfg.SetGnssAvailability(f.gnss)
	cfg.SetModel(f.model)
	cfg.SetProtocol(f.protocol)
	return cfg
}
func (f *fakeHardware) GetStatus() explorer.DeviceStatus {
	status := explorer.DeviceStatus{}
	status.SetStartedAt(f.started)
	status.SetUpdatedAt(f.updated)
	for i := int64(0); i < f.frames; i++ {
		status.IncrementFrames()
	}
	for i := int64(0); i < f.errors; i++ {
		status.IncrementErrors()
	}
	for i := int64(0); i < f.messages; i++ {
		status.IncrementMessages()
	}
	return status
}
func (f *fakeHardware) GetCoordinates(fuzzy bool) (float64, float64, float64, error) {
	f.lastFuzzy = &fuzzy
	if f.coordErr != nil {
		return 0, 0, 0, f.coordErr
	}
	return f.lat, f.lon, f.elv, nil
}
func (f *fakeHardware) GetTemperature() (float64, error) {
	if f.tempErr != nil {
		return 0, f.tempErr
	}
	return f.temp, nil
}
func (f *fakeHardware) GetDeviceId() string { return f.id }
func (f *fakeHardware) GetMetadata(string, string, string, string, string, string, string, bool) (*metadata.Render, error) {
	if f.metaErr != nil {
		return nil, f.metaErr
	}
	return f.meta, nil
}

type fakeService struct {
	name, description                                string
	running                                          bool
	restarts                                         int
	started, updated, stopped                        time.Time
	initErr, startErr, stopErr, restartErr, assetErr error
	assets                                           []service.Asset
	constraints                                      []config.IConstraint
}

func (f *fakeService) GetStatus() *service.Status {
	status := &service.Status{}
	status.SetIsRunning(f.running)
	status.SetRestarts(f.restarts)
	status.SetStartedAt(f.started)
	status.SetUpdatedAt(f.updated)
	status.SetStoppedAt(f.stopped)
	return status
}
func (f *fakeService) GetName() string        { return f.name }
func (f *fakeService) GetDescription() string { return f.description }
func (f *fakeService) Init() error            { return f.initErr }
func (f *fakeService) IsEnabled() bool        { return true }
func (f *fakeService) Start() error {
	if f.startErr != nil {
		return f.startErr
	}
	f.running = true
	return nil
}
func (f *fakeService) Stop() error {
	if f.stopErr != nil {
		return f.stopErr
	}
	f.running = false
	return nil
}
func (f *fakeService) Restart() error {
	if f.restartErr != nil {
		return f.restartErr
	}
	f.restarts++
	return nil
}
func (f *fakeService) GetAssetList() ([]service.Asset, error) {
	if f.assetErr != nil {
		return nil, f.assetErr
	}
	return f.assets, nil
}
func (f *fakeService) GetAssetData(string) (*service.AssetData, error) { return nil, nil }
func (f *fakeService) GetConfigConstraint() []config.IConstraint       { return f.constraints }

type fakeSource struct {
	property seisevent.DataSourceProperty
	events   []seisevent.Event
	err      error
}

func (f *fakeSource) GetProperty() seisevent.DataSourceProperty { return f.property }
func (f *fakeSource) GetEvents(float64, float64) ([]seisevent.Event, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.events, nil
}

type fakeConstraint struct {
	namespace, key, name, description   string
	typ                                 action.SettingType
	version                             int
	required                            bool
	options                             map[string]any
	value                               any
	getErr, setErr, restoreErr, initErr error
}

func (f *fakeConstraint) IsRequired() bool            { return f.required }
func (f *fakeConstraint) GetDefaultValue() any        { return f.value }
func (f *fakeConstraint) GetNamespace() string        { return f.namespace }
func (f *fakeConstraint) GetDescription() string      { return f.description }
func (f *fakeConstraint) GetName() string             { return f.name }
func (f *fakeConstraint) GetType() action.SettingType { return f.typ }
func (f *fakeConstraint) GetKey() string              { return f.key }
func (f *fakeConstraint) GetVersion() int             { return f.version }
func (f *fakeConstraint) GetOptions() map[string]any  { return f.options }
func (f *fakeConstraint) Init(*action.Handler) error  { return f.initErr }
func (f *fakeConstraint) Set(_ *action.Handler, newVal any) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.value = newVal
	return nil
}
func (f *fakeConstraint) Get(*action.Handler) (any, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.value, nil
}
func (f *fakeConstraint) Restore(*action.Handler) error { return f.restoreErr }

var (
	_ hardware.IHardware    = (*fakeHardware)(nil)
	_ service.IService      = (*fakeService)(nil)
	_ seisevent.IDataSource = (*fakeSource)(nil)
	_ config.IConstraint    = (*fakeConstraint)(nil)
)
