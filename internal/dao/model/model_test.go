package model

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlite_modernc"
	"gorm.io/gorm"
)

func TestModelMetadata(t *testing.T) {
	t.Parallel()

	tables := []dao.ITable{&UserSettings{}, &SysUser{}, &SchemaVersion{}}
	for _, table := range tables {
		if table.GetModel() == nil || table.GetName("pre_") == "" || !table.UseAutoMigrate() {
			t.Fatalf("%T metadata is incomplete", table)
		}
		plugins, err := table.AddPlugins(nil, "pre_")
		if err != nil || plugins != nil {
			t.Fatalf("%T plugins = %#v, %v", table, plugins, err)
		}
	}
	record := &SeisRecord{}
	if record.GetModel() == nil || record.GetName("pre_") == "" || record.UseAutoMigrate() {
		t.Fatal("seis record metadata is incomplete")
	}

	user := SysUser{}
	if user.NewUserId() == "" {
		t.Fatal("empty user id")
	}
	hashed, err := user.GetHashedPassword("Anyshake@12#$")
	if err != nil {
		t.Fatal(err)
	}
	user.HashedPassword = hashed
	if !user.IsPasswordCorrect("Anyshake@12#$") || user.IsPasswordCorrect("nope") {
		t.Fatal("password comparison mismatch")
	}
	_ = gorm.ErrRecordNotFound
}

func TestSeisRecordCodecAndShards(t *testing.T) {
	record := SeisRecord{}
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

	database := filepath.Join(t.TempDir(), "shards.db")
	db, err := (&sqlite_modernc.SQLite{}).Open("", "", "", database, "pre_", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	plugins, err := (&SeisRecord{}).AddPlugins(db, "pre_")
	if err != nil || len(plugins) != 1 {
		t.Fatalf("plugins = %#v, %v", plugins, err)
	}
}
