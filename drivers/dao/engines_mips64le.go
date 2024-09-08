package dao

import (
	"github.com/anyshake/observer/drivers/dao/engine/mariadb"
	"github.com/anyshake/observer/drivers/dao/engine/postgresql"
	"github.com/anyshake/observer/drivers/dao/engine/sqlite_ncruces"
	"github.com/anyshake/observer/drivers/dao/engine/sqlserver"
)


func createEngines() map[string]Engine {
	return map[string]Engine{
		"sqlite3":    &sqlite_ncruces.SQLite{},
		"sqlite":     &sqlite_ncruces.SQLite{},
		"mysql":      &mariadb.MariaDB{},
		"mariadb":    &mariadb.MariaDB{},
		"sqlserver":  &sqlserver.SQLServer{},
		"mssql":      &sqlserver.SQLServer{},
		"postgres":   &postgresql.PostgreSQL{},
		"postgresql": &postgresql.PostgreSQL{},
	}
}

