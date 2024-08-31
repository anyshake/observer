package dao

import (
	"fmt"
	"runtime"
	"time"

	"gorm.io/gorm"
)

type _SQLite struct{}

func (s *_SQLite) match(engine string) bool {
	return engine == "sqlite3" || engine == "sqlite"
}

func (s *_SQLite) open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	return nil, fmt.Errorf("current platform %s/%s does not support SQLite", runtime.GOOS, runtime.GOARCH)
}
