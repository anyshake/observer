package sqlite_ncruces

import (
	"path/filepath"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	if _, err := (&SQLite{}).Open("", "", "", filepath.Join(t.TempDir(), "missing", "observer.db"), "t_", time.Second); err == nil {
		t.Fatal("missing directory was accepted")
	}
	db, err := (&SQLite{}).Open("", "", "", filepath.Join(t.TempDir(), "observer.db"), "t_", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
}
