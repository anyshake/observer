package dao

import (
	"fmt"

	"gorm.io/gorm"
)

func Open(host string, port int, engineName, username, password, database string) (*gorm.DB, error) {
	engines := []engine{
		&_PostgreSQL{},
		&_MariaDB{},
		&_SQLServer{},
		&_SQLite{},
	}
	for _, e := range engines {
		if e.match(engineName) {
			return e.open(host, port, username, password, database, TIMEOUT_THRESHOLD)
		}
	}

	err := fmt.Errorf("database engine %s is unsupported", engineName)
	return nil, err
}
