package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSQL struct{}

func (p *PostgreSQL) isCompatible(engine string) bool {
	return engine == "postgres" || engine == "postgresql"
}

func (p *PostgreSQL) openDBConn(host string, port int, username, password, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d TimeZone=Etc/GMT",
		host, port, username, password, database, int(DB_TIMEOUT.Seconds()),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true, // Disable transaction to improve performance
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(DB_TIMEOUT)
	return db, nil
}

type MariaDB struct{}

func (m *MariaDB) isCompatible(engine string) bool {
	return engine == "mysql" || engine == "mariadb"
}

func (m *MariaDB) openDBConn(host string, port int, username, password, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds&loc=UTC",
		username, password, host, port, database,
		int(DB_TIMEOUT.Seconds()),
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true, // Disable transaction to improve performance
	})
}

type SQLServer struct{}

func (s *SQLServer) isCompatible(engine string) bool {
	return engine == "sqlserver" || engine == "mssql"
}

func (s *SQLServer) openDBConn(host string, port int, username, password, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		username, password, host, port, database,
	)
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true, // Disable transaction to improve performance
	})
}
