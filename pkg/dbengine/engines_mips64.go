package dbengine

import (
	"github.com/anyshake/observer/pkg/dbengine/engines/mariadb"
	"github.com/anyshake/observer/pkg/dbengine/engines/postgresql"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlite_ncruces"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlserver"
)

func New() map[string]IEngine {
	return map[string]IEngine{
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
