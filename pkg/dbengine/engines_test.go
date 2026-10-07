package dbengine_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/dbengine"
	"github.com/anyshake/observer/pkg/dbengine/engines/mariadb"
	"github.com/anyshake/observer/pkg/dbengine/engines/postgresql"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlite_modernc"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlite_ncruces"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlserver"
)

func TestEngines(t *testing.T) {
	engines := dbengine.New()
	for _, name := range []string{"sqlite", "sqlite3", "mysql", "mariadb", "postgres", "postgresql", "sqlserver", "mssql"} {
		if engines[name] == nil {
			t.Fatalf("missing engine %s", name)
		}
	}

	database := filepath.Join(t.TempDir(), "observer.db")
	sqlite := &sqlite_modernc.SQLite{}
	db, err := sqlite.Open("", "", "", database, "eng_", time.Second)
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
	if _, err := sqlite.Open("", "", "", filepath.Join(t.TempDir(), "missing", "observer.db"), "eng_", time.Second); err == nil {
		t.Fatal("missing sqlite directory accepted")
	}

	ncruces := &sqlite_ncruces.SQLite{}
	ncrucesDB, err := ncruces.Open("", "", "", filepath.Join(t.TempDir(), "ncruces.db"), "eng_", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ncrucesSQL, err := ncrucesDB.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := ncrucesSQL.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := ncruces.Open("", "", "", filepath.Join(t.TempDir(), "missing", "ncruces.db"), "eng_", time.Second); err == nil {
		t.Fatal("missing ncruces directory accepted")
	}

	if _, err := (&postgresql.PostgreSQL{}).Open("localhost", "user", "secret", "observer", "eng_", time.Second); err == nil {
		t.Fatal("postgres address without a port accepted")
	}
	if _, err := (&postgresql.PostgreSQL{}).Open("127.0.0.1:1", "user", "secret", "observer", "eng_", time.Second); err == nil {
		t.Fatal("closed postgres port accepted")
	}
	if _, err := (&mariadb.MariaDB{}).Open("127.0.0.1:1", "user", "secret", "observer", "eng_", time.Second); err == nil {
		t.Fatal("closed mariadb port accepted")
	}
	if _, err := (&sqlserver.SQLServer{}).Open("127.0.0.1:1", "user", "secret", "observer", "eng_", time.Second); err == nil {
		t.Fatal("closed sqlserver port accepted")
	}
}
