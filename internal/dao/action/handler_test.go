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

	if _, _, _, err := NewHandler(nil).SettingsGet(t.Context(), "ns", "key"); err == nil {
		t.Fatal("nil dao settings get succeeded")
	}
	if err := NewHandler(nil).SettingsSet(t.Context(), "ns", "key", String, 0, "x"); err == nil {
		t.Fatal("nil dao settings set succeeded")
	}
	if _, err := NewHandler(nil).SettingsInit(t.Context(), "ns", "key", String, 0, "x"); err == nil {
		t.Fatal("nil dao settings init succeeded")
	}

	created, err := handler.SettingsInit(t.Context(), "station", "name", String, 1, "AnyShake")
	if err != nil || created {
		t.Fatalf("init = %v, %v", created, err)
	}
	again, err := handler.SettingsInit(t.Context(), "station", "name", String, 2, "ignored")
	if err != nil || again {
		t.Fatalf("second init = %v, %v", again, err)
	}
	value, valueType, version, err := handler.SettingsGet(t.Context(), "station", "name")
	if err != nil || value != "AnyShake" || valueType != String || version != 1 {
		t.Fatalf("settings = %#v %s %d %v", value, valueType, version, err)
	}

	if err := handler.SettingsSet(t.Context(), "station", "kind", String, 1, "old"); err != nil {
		t.Fatal(err)
	}
	needsWrite, err := handler.SettingsInit(t.Context(), "station", "kind", Int, 2, int64(1))
	if err != nil || !needsWrite {
		t.Fatalf("type/version mismatch init = %v, %v", needsWrite, err)
	}
	if err := handler.SettingsSet(t.Context(), "station", "blank", String, 0, ""); err != nil {
		t.Fatal(err)
	}
	if got, _, _, err := handler.SettingsGet(t.Context(), "station", "blank"); err != nil || got != "" {
		t.Fatalf("blank string = %#v, %v", got, err)
	}
	if err := handler.SettingsSet(t.Context(), "station", "hidden", String, 0, "a\u200Bb\x01"); err != nil {
		t.Fatal(err)
	}
	if got, _, _, err := handler.SettingsGet(t.Context(), "station", "hidden"); err != nil || got != "ab" {
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
		if err := handler.SettingsSet(t.Context(), "station", tc.key, tc.kind, 0, tc.value); err != nil {
			t.Fatalf("set %s: %v", tc.key, err)
		}
		got, kind, _, err := handler.SettingsGet(t.Context(), "station", tc.key)
		if err != nil || kind != tc.kind {
			t.Fatalf("get %s = %#v %s %v", tc.key, got, kind, err)
		}
	}
	if err := handler.SettingsSet(t.Context(), "station", "bad", Bool, 0, "nope"); err == nil {
		t.Fatal("bool accepted a string")
	}
	if err := handler.SettingsSet(t.Context(), "station", "bad", String, 0, 1); err == nil {
		t.Fatal("string accepted an int")
	}
	if err := handler.SettingsSet(t.Context(), "station", "bad", SettingType("nope"), 0, 1); err == nil {
		t.Fatal("unknown setting type accepted")
	}
	if err := handler.daoObj.Database.Create(&model.UserSettings{
		Namespace: "station", ConfigKey: "mystery", ConfigType: "mystery", ConfigValue: []byte("x"),
	}).Error; err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := handler.SettingsGet(t.Context(), "station", "mystery"); err == nil {
		t.Fatal("unknown stored type accepted")
	}
	if _, _, _, err := handler.SettingsGet(t.Context(), "station", "missing"); err == nil {
		t.Fatal("missing setting accepted")
	}

	if _, err := handler.SysUserCreate(t.Context(), "ab", "Anyshake@12#$", true); err == nil {
		t.Fatal("short username accepted")
	}
	userID, err := handler.SysUserCreate(t.Context(), "anyshake_admin", "Anyshake@12#$", true)
	if err != nil || userID == "" {
		t.Fatal(err)
	}
	if _, err := handler.SysUserCreate(t.Context(), "anyshake_admin", "Anyshake@12#$", false); err == nil {
		t.Fatal("duplicate user accepted")
	}
	hasAdmin, err := handler.SysUserHasAdmin(t.Context())
	if err != nil || !hasAdmin {
		t.Fatalf("has admin = %v, %v", hasAdmin, err)
	}
	users, err := handler.SysUserList(t.Context())
	if err != nil || len(users) != 1 {
		t.Fatalf("users = %d, %v", len(users), err)
	}
	if _, err := handler.SysUserLogin(t.Context(), "anyshake_admin", "wrong-password", "", ""); err == nil {
		t.Fatal("bad password accepted")
	}
	loggedIn, err := handler.SysUserLogin(t.Context(), "anyshake_admin", "Anyshake@12#$", "", "")
	if err != nil || loggedIn != userID {
		t.Fatalf("login = %s, %v", loggedIn, err)
	}
	if _, err := handler.SysUserLogin(t.Context(), "anyshake_admin", "Anyshake@12#$", "agent", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	stored, err := handler.SysUserGetByUserId(t.Context(), userID)
	if err != nil || stored.UserIp != "127.0.0.1" {
		t.Fatalf("stored user = %+v, %v", stored, err)
	}
	stored.Username = "renamed_admin"
	if err := handler.SysUserUpdate(t.Context(), userID, stored); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.SysUserGetByUsername(t.Context(), "missing"); err == nil {
		t.Fatal("missing username accepted")
	}
	if err := handler.SysUserRemove(t.Context(), userID); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.SysUserGetByUserId(t.Context(), userID); err == nil {
		t.Fatal("removed user still exists")
	}
	if err := handler.SysUserUpdate(t.Context(), userID, stored); err == nil {
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
	if _, err := nilHandler.SysUserHasAdmin(t.Context()); err == nil {
		t.Fatal("nil dao has-admin succeeded")
	}
	if _, err := nilHandler.SysUserList(t.Context()); err == nil {
		t.Fatal("nil dao list succeeded")
	}
	if _, err := nilHandler.SysUserCreate(t.Context(), "valid_name", "GoodPass1!", true); err == nil {
		t.Fatal("nil dao create succeeded")
	}
	if _, err := nilHandler.SysUserLogin(t.Context(), "valid_name", "GoodPass1!", "", ""); err == nil {
		t.Fatal("nil dao login succeeded")
	}
	if err := nilHandler.SysUserUpdate(t.Context(), "id", model.SysUser{}); err == nil {
		t.Fatal("nil dao update succeeded")
	}
	if err := nilHandler.SysUserRemove(t.Context(), "id"); err == nil {
		t.Fatal("nil dao remove succeeded")
	}
}

func TestSchemaVersion(t *testing.T) {
	_, handler := openActionDB(t)
	if err := handler.SchemaVersionInit(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := handler.SchemaVersionInit(t.Context()); err != nil {
		t.Fatal(err)
	}
	version, err := handler.SchemaVersionGetCurrent(t.Context())
	if err != nil || version != 1 {
		t.Fatalf("version = %d, %v", version, err)
	}
	if err := handler.SchemaVersionUpdate(t.Context(), 1, 2, "next"); err != nil {
		t.Fatal(err)
	}
	if err := handler.SchemaVersionUpdate(t.Context(), 1, 3, "stale"); !errors.Is(err, gorm.ErrInvalidData) {
		t.Fatalf("stale update error = %v", err)
	}
}

func TestSeismicRecords(t *testing.T) {
	obj, handler := openActionDB(t)
	if err := obj.AutoMigrate(&model.SeisRecord{}); err != nil {
		t.Fatal(err)
	}

	if err := NewHandler(nil).SeisRecordsCreate(t.Context(), model.SeisRecord{}); err == nil {
		t.Fatal("nil dao create succeeded")
	}
	if err := NewHandler(nil).SeisRecordsQueryEach(context.Background(), time.Now(), time.Now(), func(model.SeisRecord) error { return nil }); err == nil {
		t.Fatal("nil dao query succeeded")
	}
	if err := handler.SeisRecordsQueryEach(nil, time.Now(), time.Now(), func(model.SeisRecord) error { return nil }); err == nil {
		t.Fatal("nil context accepted")
	}
	if err := handler.SeisRecordsQueryEach(context.Background(), time.Now(), time.Now(), nil); err == nil {
		t.Fatal("nil callback accepted")
	}
	end := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	start := end.Add(-2 * time.Hour)
	if _, err := handler.SeisRecordsQuery(t.Context(), start, end); err == nil {
		t.Fatal("oversized window accepted")
	}
	if _, err := handler.SeisRecordsQuery(t.Context(), end, start); err == nil {
		t.Fatal("reversed window accepted")
	}
	if err := handler.SeisRecordsPurge(t.Context(), end, start); err == nil {
		t.Fatal("reversed purge accepted")
	}

	first := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	second := first.Add(time.Minute)
	nextDay := first.Add(24 * time.Hour)
	records := []model.SeisRecord{{RecordTime: first.UnixMilli(), SampleRate: 100, ChannelData: []byte{1}}, {RecordTime: second.UnixMilli(), SampleRate: 100, ChannelData: []byte{2}}}
	if err := handler.SeisRecordsCreate(t.Context(), records...); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(t.Context(), model.SeisRecord{RecordTime: nextDay.UnixMilli(), SampleRate: 50, ChannelData: []byte{3}}); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(t.Context(), records[0]); err == nil {
		t.Fatal("duplicate record accepted")
	}

	got, err := handler.SeisRecordsQuery(t.Context(), first.Add(-time.Second), second.Add(time.Second))
	if err != nil || len(got) != 2 {
		t.Fatalf("query = %d, %v", len(got), err)
	}
	if err := handler.SeisRecordsQueryEach(t.Context(), first, second, func(model.SeisRecord) error {
		return errors.New("stop")
	}); err == nil || err.Error() != "stop" {
		t.Fatalf("callback error = %v", err)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := handler.SeisRecordsQueryEach(canceled, first, second, func(model.SeisRecord) error { return nil }); err == nil {
		t.Fatal("canceled query succeeded")
	}

	sqlDB, err := obj.Database.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(4)
	if err := handler.SeisRecordsPurge(t.Context(), first, second); err != nil {
		t.Fatal(err)
	}
	remaining, err := handler.SeisRecordsQuery(t.Context(), nextDay.Add(-time.Second), nextDay.Add(time.Second))
	if err != nil || len(remaining) != 1 {
		t.Fatalf("remaining = %d, %v", len(remaining), err)
	}
	if err := handler.SeisRecordsPurgeAll(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got, err := handler.SeisRecordsQuery(t.Context(), nextDay.Add(-time.Second), nextDay.Add(time.Second)); err != nil || len(got) != 0 {
		t.Fatalf("after purge all = %d, %v", len(got), err)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- handler.CleanupExclusive(t.Context(), func() error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started
	if err := handler.CleanupExclusive(t.Context(), func() error { return nil }); !errors.Is(err, ErrCleanupRunning) {
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

func TestActionContextCancellation(t *testing.T) {
	_, handler := openActionDB(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	now := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)

	operations := []struct {
		name string
		run  func() error
	}{
		{"SettingsGet", func() error {
			_, _, _, err := handler.SettingsGet(ctx, "station", "name")
			return err
		}},
		{"SettingsSet", func() error {
			return handler.SettingsSet(ctx, "station", "name", String, 0, "station")
		}},
		{"SettingsInit", func() error {
			_, err := handler.SettingsInit(ctx, "station", "name", String, 0, "station")
			return err
		}},
		{"SysUserHasAdmin", func() error {
			_, err := handler.SysUserHasAdmin(ctx)
			return err
		}},
		{"SysUserList", func() error {
			_, err := handler.SysUserList(ctx)
			return err
		}},
		{"SysUserGetByUsername", func() error {
			_, err := handler.SysUserGetByUsername(ctx, "valid_name")
			return err
		}},
		{"SysUserGetByUserId", func() error {
			_, err := handler.SysUserGetByUserId(ctx, "user-id")
			return err
		}},
		{"SysUserCreate", func() error {
			_, err := handler.SysUserCreate(ctx, "valid_name", "GoodPass1!", true)
			return err
		}},
		{"SysUserLogin", func() error {
			_, err := handler.SysUserLogin(ctx, "valid_name", "GoodPass1!", "", "")
			return err
		}},
		{"SysUserUpdate", func() error {
			return handler.SysUserUpdate(ctx, "user-id", model.SysUser{})
		}},
		{"SysUserRemove", func() error {
			return handler.SysUserRemove(ctx, "user-id")
		}},
		{"SchemaVersionInit", func() error { return handler.SchemaVersionInit(ctx) }},
		{"SchemaVersionGetCurrent", func() error {
			_, err := handler.SchemaVersionGetCurrent(ctx)
			return err
		}},
		{"SchemaVersionUpdate", func() error { return handler.SchemaVersionUpdate(ctx, 1, 2, "next") }},
		{"SeisRecordsCreate", func() error {
			return handler.SeisRecordsCreate(ctx, model.SeisRecord{RecordTime: now.UnixMilli()})
		}},
		{"SeisRecordsQuery", func() error {
			_, err := handler.SeisRecordsQuery(ctx, now, now.Add(time.Minute))
			return err
		}},
		{"SeisRecordsQueryEach", func() error {
			return handler.SeisRecordsQueryEach(ctx, now, now.Add(time.Minute), func(model.SeisRecord) error {
				return errors.New("canceled query called its callback")
			})
		}},
		{"SeisRecordsPurge", func() error { return handler.SeisRecordsPurge(ctx, now, now) }},
		{"SeisRecordsPurgeAll", func() error { return handler.SeisRecordsPurgeAll(ctx) }},
		{"CleanupExclusive", func() error {
			return handler.CleanupExclusive(ctx, func() error {
				return errors.New("canceled cleanup called its callback")
			})
		}},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			if err := operation.run(); !errors.Is(err, context.Canceled) {
				t.Fatalf("error = %v; want context.Canceled", err)
			}
		})
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
