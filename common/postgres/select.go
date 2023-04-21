package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func SelectData(db *sql.DB, timestamp int64) ([]map[string]interface{}, error) {
	var (
		start = timestamp - 35000
		end   = timestamp + 35000
	)

	rows, err := db.Query(fmt.Sprintf(`
		SELECT * FROM acceleration
		WHERE timestamp >= to_timestamp(%d / 1000.0)
		AND timestamp <= to_timestamp(%d / 1000.0)
	`, start, end))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resMap []map[string]interface{}
	for rows.Next() {
		var (
			id      int
			ts      time.Time
			station string
			uuid    string
			data    string
		)
		err := rows.Scan(&id, &ts, &station, &uuid, &data)
		if err != nil {
			return nil, err
		}

		resMap = append(resMap, map[string]interface{}{
			"id":        id,
			"timestamp": ts.UnixMilli(),
			"station":   station,
			"uuid":      uuid,
			"data":      data,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resMap, nil
}
