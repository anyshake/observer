package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func OpenPostgres(host string, port int, username, password, database string, enable bool) (*sql.DB, error) {
	if !enable {
		return nil, nil
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		host, port, username, password, database,
	))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	InitPostgres(db)
	return db, nil
}

func ClosePostgres(db *sql.DB) error {
	return db.Close()
}

func InitPostgres(db *sql.DB) error {
	err := CreateTable(db)
	if err != nil {
		return err
	}

	rows, err := db.Query(`
        SELECT column_name
        FROM information_schema.columns
        WHERE table_name = 'acceleration'
		AND column_name IN ('id', 'timestamp', 'station', 'uuid', 'data')
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	columnMap := make(map[string]bool)
	for rows.Next() {
		var columnName string
		err := rows.Scan(&columnName)
		if err != nil {
			return err
		}
		columnMap[columnName] = true
	}

	for _, v := range []string{
		"id", "timestamp", "station", "uuid", "data",
	} {
		_, ok := columnMap[v]
		if !ok {
			_, err := db.Exec("DROP TABLE acceleration")
			if err != nil {
				return err
			}

			err = CreateTable(db)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}
