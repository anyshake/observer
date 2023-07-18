package postgres

import (
	"database/sql"

	"github.com/lib/pq"
)

func Insert(db *sql.DB, ts int64, ehz, ehe, ehn []int32) error {
	var (
		ehzArray = pq.Array(ehz)
		eheArray = pq.Array(ehe)
		ehnArray = pq.Array(ehn)
	)

	_, err := db.Exec(`
        INSERT INTO counts (ts, ehz, ehe, ehn)
        VALUES (to_timestamp($1 / 1000.0), $2, $3, $4)
	`, ts, ehzArray, eheArray, ehnArray)
	if err != nil {
		return err
	}

	return nil
}
