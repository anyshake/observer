package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InsertData(db *sql.DB, timestamp int64, station, uuid string, data []byte) error {
	_, err := db.Exec(fmt.Sprintf(`
		INSERT INTO acceleration (timestamp, station, uuid, data)
		VALUES (to_timestamp(%d / 1000.0), TEXT '%s', TEXT '%s', '%s')
	`, timestamp, station, uuid, string(data)))
	if err != nil {
		return err
	}

	return nil
}
