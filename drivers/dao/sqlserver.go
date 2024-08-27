package dao

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type _SQLServer struct{}

func (s *_SQLServer) match(engine string) bool {
	return engine == "sqlserver" || engine == "mssql"
}

func (s *_SQLServer) open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		username, password, host, port, database,
	)
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
