package dao

import "gorm.io/gorm"

func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	sqlDB.Close()
	return err
}
