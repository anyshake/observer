package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Open(host string, port int, username, password, database string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		host, port, username, password, database,
	))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	Init(db)
	return db, nil
}

func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}

	return db.Close()
}
