package testsupport

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
)

// OpenDAO opens a temporary SQLite database and migrates the settings, user, and schema tables.
func OpenDAO(t *testing.T) (*dao.DAO, *action.Handler) {
	t.Helper()

	database := filepath.Join(t.TempDir(), "test.db")
	obj, err := dao.New("sqlite3://localhost", "", "", "cov_", time.Second)
	if err != nil {
		t.Fatalf("dao.New: %v", err)
	}
	if err := obj.Open(database); err != nil {
		t.Fatalf("dao.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = obj.Close()
	})
	if err := obj.AutoMigrate(&model.UserSettings{}, &model.SysUser{}, &model.SchemaVersion{}); err != nil {
		t.Fatalf("dao.AutoMigrate: %v", err)
	}
	return obj, action.NewHandler(obj)
}
