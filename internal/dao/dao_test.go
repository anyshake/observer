package dao_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/model"
)

func TestDAOLifecycle(t *testing.T) {
	database := filepath.Join(t.TempDir(), "observer.db")
	obj, err := dao.New("sqlite3://localhost", "user", "secret", "app_", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if obj.GetPrefix() != "app_" {
		t.Fatalf("prefix = %q", obj.GetPrefix())
	}
	if err := obj.Close(); err == nil {
		t.Fatal("close before open succeeded")
	}
	if err := obj.AutoMigrate(&model.UserSettings{}); err == nil {
		t.Fatal("migrate before open succeeded")
	}
	if err := obj.Open(database); err != nil {
		t.Fatal(err)
	}
	if err := obj.Open(database); err == nil {
		t.Fatal("second open succeeded")
	}
	if err := obj.AutoMigrate(&model.UserSettings{}, &model.SysUser{}, &model.SchemaVersion{}); err != nil {
		t.Fatal(err)
	}
	if err := obj.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDAORejectsBadEndpoints(t *testing.T) {
	t.Parallel()

	if _, err := dao.New("://bad", "", "", "", time.Second); err == nil {
		t.Fatal("invalid endpoint accepted")
	}
	if _, err := dao.New("oracle://localhost", "", "", "", time.Second); err == nil {
		t.Fatal("unknown engine accepted")
	}
}
