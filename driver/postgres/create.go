package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS counts (
            id SERIAL PRIMARY KEY,
            ts TIMESTAMPTZ,
            ehz INTEGER [],
            ehn INTEGER [],
            ehe INTEGER []
        )
    `)
	if err != nil {
		return err
	}

	return nil
}
