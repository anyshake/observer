package dao

import "errors"

func (t *DAO) Close() error {
	if t.Database == nil {
		return errors.New("database is not opened")
	}

	sqlDB, err := t.Database.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
