package dao

import (
	"fmt"
	"time"

	"github.com/bclswl0827/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLite struct{}

func (s *SQLite) match(engine string) bool {
	return engine == "sqlite3" || engine == "sqlite"
}

func (s *SQLite) open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file://%s?cache=shared&mode=rwc&_pragma=busy_timeout(%d)", database, int(timeout.Seconds()))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db, err
}
