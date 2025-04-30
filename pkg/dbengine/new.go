//go:build !mips && !mips64 && !mipsle && !mips64le && !windows && !openbsd
// +build !mips,!mips64,!mipsle,!mips64le,!windows,!openbsd

package dbengine

import (
	"github.com/anyshake/observer/pkg/dbengine/engines/mariadb"
	"github.com/anyshake/observer/pkg/dbengine/engines/postgresql"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlite_modernc"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlserver"
)

func New() map[string]IEngine {
	return map[string]IEngine{
		"sqlite3":    &sqlite_modernc.SQLite{},
		"sqlite":     &sqlite_modernc.SQLite{},
		"mysql":      &mariadb.MariaDB{},
		"mariadb":    &mariadb.MariaDB{},
		"sqlserver":  &sqlserver.SQLServer{},
		"mssql":      &sqlserver.SQLServer{},
		"postgres":   &postgresql.PostgreSQL{},
		"postgresql": &postgresql.PostgreSQL{},
	}
}
