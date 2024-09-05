//go:build !mips && !mips64 && !mipsle && !mips64le && !windows && !openbsd
// +build !mips,!mips64,!mipsle,!mips64le,!windows,!openbsd

package dao

import (
	"github.com/anyshake/observer/drivers/dao/engine/mariadb"
	"github.com/anyshake/observer/drivers/dao/engine/postgresql"
	"github.com/anyshake/observer/drivers/dao/engine/sqlite_modernc"
	"github.com/anyshake/observer/drivers/dao/engine/sqlserver"
)

func createEngines() map[string]Engine {
	return map[string]Engine{
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
