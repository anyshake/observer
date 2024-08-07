package dao

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSQL struct{}

func (p *PostgreSQL) match(engine string) bool {
	return engine == "postgres" || engine == "postgresql"
}

func (p *PostgreSQL) open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d TimeZone=Etc/GMT",
		host, port, username, password, database, int(timeout.Seconds()),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(timeout)
	return db, nil
}
