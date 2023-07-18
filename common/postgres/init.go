package postgres

import "database/sql"

func Init(db *sql.DB) error {
	err := CreateTable(db)
	if err != nil {
		return err
	}

	rows, err := db.Query(`
        SELECT column_name
        FROM information_schema.columns
        WHERE table_name = 'counts'
        AND column_name IN ('id', 'ts', 'ehz', 'ehe', 'ehn')
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
		"id", "ts", "ehz", "ehe", "ehn",
	} {
		_, ok := columnMap[v]
		if !ok {
			_, err := db.Exec("DROP TABLE counts")
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
