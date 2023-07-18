package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

func Query(db *sql.DB, start, end int64) ([]map[string]any, error) {
	rows, err := db.Query(`
        SELECT * FROM counts WHERE
		ts >= to_timestamp($1 / 1000.0)
		AND ts <= to_timestamp($2 / 1000.0)
	`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resMap []map[string]any
	for rows.Next() {
		var (
			id  int
			ehz []int32
			ehe []int32
			ehn []int32
			ts  time.Time
		)
		err := rows.Scan(&id, &ts,
			(*pq.Int32Array)(&ehz),
			(*pq.Int32Array)(&ehe),
			(*pq.Int32Array)(&ehn),
		)
		if err != nil {
			return nil, err
		}

		resMap = append(resMap, map[string]any{
			"id":  id,
			"ehz": ehz,
			"ehe": ehe,
			"ehn": ehn,
			"ts":  ts.UnixMilli(),
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resMap, nil
}
