package mariadb

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func (m *MariaDB) Open(address, username, password, database, prefix string, timeout time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&timeout=%ds&loc=UTC",
		username, password, address, database, int(timeout.Seconds()),
	)
	return gorm.Open(
		mysql.New(mysql.Config{DSN: dsn, SkipInitializeWithVersion: false}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: prefix,
			},
		},
	)
}
