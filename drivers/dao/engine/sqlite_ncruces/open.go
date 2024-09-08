package sqlite_ncruces

import (
	"time"

	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (s *SQLite) Open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(gormlite.Open(database), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db, err
}
