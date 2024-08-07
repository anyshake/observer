package dao

import "gorm.io/gorm"

func Close(dbObj *gorm.DB) error {
	if dbObj == nil {
		return nil
	}

	sqlDB, err := dbObj.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
