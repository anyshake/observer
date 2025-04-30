package dao

import (
	"errors"
	"fmt"
)

func (d *DAO) Open(database string) error {
	if d.driver == nil {
		return errors.New("database driver is not set")
	}

	if d.Database != nil {
		return errors.New("database is already opened")
	}

	dbObj, err := d.driver.Open(d.address, d.username, d.password, database, d.prefix, d.timeout)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	d.Database = dbObj
	return nil
}
