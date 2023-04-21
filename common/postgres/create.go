package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS acceleration (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP,
			station TEXT,
			uuid TEXT,
			data JSON
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
