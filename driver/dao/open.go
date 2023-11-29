package dao

import (
	"fmt"

	"gorm.io/gorm"
)

func Open(host string, port int, engine, username, password, database string) (*gorm.DB, error) {
	engines := []dbEngine{
		&PostgreSQL{}, &MariaDB{},
		&SQLite{}, &SQLServer{},
	}
	for _, e := range engines {
		if e.isCompatible(engine) {
			return e.openDBConn(host, port, username, password, database)
		}
	}

	err := fmt.Errorf("database engine %s is unsupported", engine)
	return nil, err
}
