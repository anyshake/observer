package sqlite_modernc

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bclswl0827/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func (s *SQLite) Open(address, username, password, database, prefix string, timeout time.Duration) (*gorm.DB, error) {
	baseDir := filepath.Dir(database)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", baseDir)
	}

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?_pragma=synchronous(NORMAL)&_pragma=cache_size(-20000)&_pragma=temp_store(MEMORY)&_pragma=foreign_keys(ON)&_pragma=auto_vacuum(FULL)", database)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: prefix,
		},
	})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db, err
}
