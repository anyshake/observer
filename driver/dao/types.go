package dao

import (
	"time"

	"github.com/bclswl0827/observer/publisher"
	"gorm.io/gorm"
)

const (
	DB_TIMEZONE  = "Etc/GMT"
	DB_TIMEOUT   = 5 * time.Second
	DB_TABLENAME = "geophone_records"
)

type dbRecord struct {
	ID uint `gorm:"primarykey"`
	publisher.Geophone
}

type dbEngine interface {
	isCompatible(engine string) bool
	openDBConn(host string, port int, username, password, database string) (*gorm.DB, error)
}
