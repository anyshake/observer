package dao

import (
	"time"

	"github.com/bclswl0827/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type _SQLite struct{}

func (s *_SQLite) match(engine string) bool {
	return engine == "sqlite3" || engine == "sqlite"
}

func (s *_SQLite) open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db, err
}
