package action

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"gorm.io/gorm"
)

func TestSettingsAndUsers(t *testing.T) {
	_, handler := openActionDB(t)

	if _, _, _, err := NewHandler(nil).SettingsGet("ns", "key"); err == nil {
		t.Fatal("nil dao settings get succeeded")
	}
	if err := NewHandler(nil).SettingsSet("ns", "key", String, 0, "x"); err == nil {
		t.Fatal("nil dao settings set succeeded")
	}
	if _, err := NewHandler(nil).SettingsInit("ns", "key", String, 0, "x"); err == nil {
		t.Fatal("nil dao settings init succeeded")
	}

	created, err := handler.SettingsInit("station", "name", String, 1, "AnyShake")
	if err != nil || created {
		t.Fatalf("init = %v, %v", created, err)
	}
	again, err := handler.SettingsInit("station", "name", String, 2, "ignored")
	if err != nil || again {
		t.Fatalf("second init = %v, %v", again, err)
	}
	value, valueType, version, err := handler.SettingsGet("station", "name")
	if err != nil || value != "AnyShake" || valueType != String || version != 1 {
		t.Fatalf("settings = %#v %s %d %v", value, valueType, version, err)
	}

	if err := handler.SettingsSet("station", "kind", String, 1, "old"); err != nil {
		t.Fatal(err)
	}
	needsWrite, err := handler.SettingsInit("station", "kind", Int, 2, int64(1))
	if err != nil || !needsWrite {
		t.Fatalf("type/version mismatch init = %v, %v", needsWrite, err)
	}
	if err := handler.SettingsSet("station", "blank", String, 0, ""); err != nil {
		t.Fatal(err)
	}
	if got, _, _, err := handler.SettingsGet("station", "blank"); err != nil || got != "" {
		t.Fatalf("blank string = %#v, %v", got, err)
	}
	if err := handler.SettingsSet("station", "hidden", String, 0, "a\u200Bb\x01"); err != nil {
		t.Fatal(err)
	}
	if got, _, _, err := handler.SettingsGet("station", "hidden"); err != nil || got != "ab" {
		t.Fatalf("cleaned string = %#v, %v", got, err)
	}

	for _, tc := range []struct {
		key   string
		kind  SettingType
		value any
	}{
		{"enabled", Bool, false},
		{"count", Int, int64(7)},
		{"ratio", Float, 1.5},
		{"codes", StringArray, []string{"EHZ"}},
		{"days", IntArray, []int64{1, 2}},
		{"gains", FloatArray, []float64{0.5}},
	} {
		if err := handler.SettingsSet("station", tc.key, tc.kind, 0, tc.value); err != nil {
			t.Fatalf("set %s: %v", tc.key, err)
		}
		got, kind, _, err := handler.SettingsGet("station", tc.key)
		if err != nil || kind != tc.kind {
			t.Fatalf("get %s = %#v %s %v", tc.key, got, kind, err)
		}
	}
	if err := handler.SettingsSet("station", "bad", Bool, 0, "nope"); err == nil {
		t.Fatal("bool accepted a string")
	}
	if err := handler.SettingsSet("station", "bad", String, 0, 1); err == nil {
		t.Fatal("string accepted an int")
	}
	if err := handler.SettingsSet("station", "bad", SettingType("nope"), 0, 1); err == nil {
		t.Fatal("unknown setting type accepted")
	}
	if err := handler.daoObj.Database.Create(&model.UserSettings{
		Namespace: "station", ConfigKey: "mystery", ConfigType: "mystery", ConfigValue: []byte("x"),
	}).Error; err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := handler.SettingsGet("station", "mystery"); err == nil {
		t.Fatal("unknown stored type accepted")
	}
	if _, _, _, err := handler.SettingsGet("station", "missing"); err == nil {
		t.Fatal("missing setting accepted")
	}

	if _, err := handler.SysUserCreate("ab", "Anyshake@12#$", true); err == nil {
		t.Fatal("short username accepted")
	}
	userID, err := handler.SysUserCreate("anyshake_admin", "Anyshake@12#$", true)
	if err != nil || userID == "" {
		t.Fatal(err)
	}
	if _, err := handler.SysUserCreate("anyshake_admin", "Anyshake@12#$", false); err == nil {
		t.Fatal("duplicate user accepted")
	}
	hasAdmin, err := handler.SysUserHasAdmin()
	if err != nil || !hasAdmin {
		t.Fatalf("has admin = %v, %v", hasAdmin, err)
	}
	users, err := handler.SysUserList()
	if err != nil || len(users) != 1 {
		t.Fatalf("users = %d, %v", len(users), err)
	}
	if _, err := handler.SysUserLogin("anyshake_admin", "wrong-password", "", ""); err == nil {
		t.Fatal("bad password accepted")
	}
	loggedIn, err := handler.SysUserLogin("anyshake_admin", "Anyshake@12#$", "", "")
	if err != nil || loggedIn != userID {
		t.Fatalf("login = %s, %v", loggedIn, err)
	}
	if _, err := handler.SysUserLogin("anyshake_admin", "Anyshake@12#$", "agent", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	stored, err := handler.SysUserGetByUserId(userID)
	if err != nil || stored.UserIp != "127.0.0.1" {
		t.Fatalf("stored user = %+v, %v", stored, err)
	}
	stored.Username = "renamed_admin"
	if err := handler.SysUserUpdate(userID, stored); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.SysUserGetByUsername("missing"); err == nil {
		t.Fatal("missing username accepted")
	}
	if err := handler.SysUserRemove(userID); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.SysUserGetByUserId(userID); err == nil {
		t.Fatal("removed user still exists")
	}
	if err := handler.SysUserUpdate(userID, stored); err == nil {
		t.Fatal("update of missing user succeeded")
	}

	for _, password := range []string{"short", "alllowercase1!", "NoSpecial123"} {
		if err := handler.SysUserCheckPassword(password); err == nil {
			t.Fatalf("password %q accepted", password)
		}
	}
	if err := handler.SysUserCheckPassword("GoodPass1!"); err != nil {
		t.Fatal(err)
	}

	nilHandler := NewHandler(nil)
	if _, err := nilHandler.SysUserHasAdmin(); err == nil {
		t.Fatal("nil dao has-admin succeeded")
	}
	if _, err := nilHandler.SysUserList(); err == nil {
		t.Fatal("nil dao list succeeded")
	}
	if _, err := nilHandler.SysUserCreate("valid_name", "GoodPass1!", true); err == nil {
		t.Fatal("nil dao create succeeded")
	}
	if _, err := nilHandler.SysUserLogin("valid_name", "GoodPass1!", "", ""); err == nil {
		t.Fatal("nil dao login succeeded")
	}
	if err := nilHandler.SysUserUpdate("id", model.SysUser{}); err == nil {
		t.Fatal("nil dao update succeeded")
	}
	if err := nilHandler.SysUserRemove("id"); err == nil {
		t.Fatal("nil dao remove succeeded")
	}
}

func TestSchemaVersion(t *testing.T) {
	_, handler := openActionDB(t)
	if err := handler.SchemaVersionInit(); err != nil {
		t.Fatal(err)
	}
	if err := handler.SchemaVersionInit(); err != nil {
		t.Fatal(err)
	}
	version, err := handler.SchemaVersionGetCurrent()
	if err != nil || version != 1 {
		t.Fatalf("version = %d, %v", version, err)
	}
	if err := handler.SchemaVersionUpdate(1, 2, "next"); err != nil {
		t.Fatal(err)
	}
	if err := handler.SchemaVersionUpdate(1, 3, "stale"); !errors.Is(err, gorm.ErrInvalidData) {
		t.Fatalf("stale update error = %v", err)
	}
}

func TestSeismicRecords(t *testing.T) {
	obj, handler := openActionDB(t)
	if err := obj.AutoMigrate(&model.SeisRecord{}); err != nil {
		t.Fatal(err)
	}

	if err := NewHandler(nil).SeisRecordsCreate(model.SeisRecord{}); err == nil {
		t.Fatal("nil dao create succeeded")
	}
	if err := NewHandler(nil).SeisRecordsQueryEachContext(context.Background(), time.Now(), time.Now(), func(model.SeisRecord) error { return nil }); err == nil {
		t.Fatal("nil dao query succeeded")
	}
	if err := handler.SeisRecordsQueryEachContext(nil, time.Now(), time.Now(), func(model.SeisRecord) error { return nil }); err == nil {
		t.Fatal("nil context accepted")
	}
	if err := handler.SeisRecordsQueryEachContext(context.Background(), time.Now(), time.Now(), nil); err == nil {
		t.Fatal("nil callback accepted")
	}
	end := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	start := end.Add(-2 * time.Hour)
	if _, err := handler.SeisRecordsQuery(start, end); err == nil {
		t.Fatal("oversized window accepted")
	}
	if _, err := handler.SeisRecordsQuery(end, start); err == nil {
		t.Fatal("reversed window accepted")
	}
	if err := handler.SeisRecordsPurge(end, start); err == nil {
		t.Fatal("reversed purge accepted")
	}

	first := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	second := first.Add(time.Minute)
	nextDay := first.Add(24 * time.Hour)
	records := []model.SeisRecord{{RecordTime: first.UnixMilli(), SampleRate: 100, ChannelData: []byte{1}}, {RecordTime: second.UnixMilli(), SampleRate: 100, ChannelData: []byte{2}}}
	if err := handler.SeisRecordsCreate(records...); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(model.SeisRecord{RecordTime: nextDay.UnixMilli(), SampleRate: 50, ChannelData: []byte{3}}); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(records[0]); err == nil {
		t.Fatal("duplicate record accepted")
	}

	got, err := handler.SeisRecordsQuery(first.Add(-time.Second), second.Add(time.Second))
	if err != nil || len(got) != 2 {
		t.Fatalf("query = %d, %v", len(got), err)
	}
	if err := handler.SeisRecordsQueryEach(first, second, func(model.SeisRecord) error {
		return errors.New("stop")
	}); err == nil || err.Error() != "stop" {
		t.Fatalf("callback error = %v", err)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := handler.SeisRecordsQueryEachContext(canceled, first, second, func(model.SeisRecord) error { return nil }); err == nil {
		t.Fatal("canceled query succeeded")
	}

	sqlDB, err := obj.Database.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(4)
	if err := handler.SeisRecordsPurge(first, second); err != nil {
		t.Fatal(err)
	}
	remaining, err := handler.SeisRecordsQuery(nextDay.Add(-time.Second), nextDay.Add(time.Second))
	if err != nil || len(remaining) != 1 {
		t.Fatalf("remaining = %d, %v", len(remaining), err)
	}
	if err := handler.SeisRecordsPurgeAll(); err != nil {
		t.Fatal(err)
	}
	if got, err := handler.SeisRecordsQuery(nextDay.Add(-time.Second), nextDay.Add(time.Second)); err != nil || len(got) != 0 {
		t.Fatalf("after purge all = %d, %v", len(got), err)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- handler.CleanupExclusive(func() error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started
	if err := handler.CleanupExclusive(func() error { return nil }); !errors.Is(err, ErrCleanupRunning) {
		t.Fatalf("overlapping cleanup = %v", err)
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}

	if _, err := handler.quoteSQLIdentifier(""); err == nil {
		t.Fatal("empty identifier accepted")
	}
	if _, err := handler.quoteSQLIdentifier("bad-name"); err == nil {
		t.Fatal("hyphenated identifier accepted")
	}
	quoted, err := handler.quoteSQLIdentifier("seis_records_1")
	if err != nil || quoted == "" {
		t.Fatalf("quoted = %q, %v", quoted, err)
	}
}

func openActionDB(t *testing.T) (*dao.DAO, *Handler) {
	t.Helper()
	obj, err := dao.New("sqlite3://localhost", "", "", "cov_", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err := obj.Open(filepath.Join(t.TempDir(), "action.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = obj.Close() })
	if err := obj.AutoMigrate(&model.UserSettings{}, &model.SysUser{}, &model.SchemaVersion{}); err != nil {
		t.Fatal(err)
	}
	return obj, NewHandler(obj)
}

func TestRecordCodec(t *testing.T) {
	t.Parallel()

	record := model.SeisRecord{}
	when := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	channels := []explorer.ChannelData{{ChannelCode: "EHZ", ChannelId: 1, ByteSize: 4, DataType: "int32", Data: []int32{1, -2, 3}}}
	if err := record.Encode(when, 100, channels); err != nil {
		t.Fatal(err)
	}
	gotWhen, rate, gotChannels, err := record.Decode()
	if err != nil || !gotWhen.Equal(when) || rate != 100 || len(gotChannels) != 1 || gotChannels[0].Data[1] != -2 {
		t.Fatalf("decoded = %v %d %#v %v", gotWhen, rate, gotChannels, err)
	}
	record.ChannelData = []byte{0}
	if _, _, _, err := record.Decode(); err == nil {
		t.Fatal("truncated channel data accepted")
	}
	if bytes.Equal(record.ChannelData, nil) {
		t.Fatal("channel data unexpectedly cleared")
	}
}
