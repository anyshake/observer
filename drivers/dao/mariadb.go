package dao

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MariaDB struct{}

func (m *MariaDB) match(engine string) bool {
	return engine == "mysql" || engine == "mariadb"
}

func (m *MariaDB) open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds&loc=UTC",
		username, password, host, port, database, int(timeout.Seconds()),
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
