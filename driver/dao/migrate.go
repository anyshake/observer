package dao

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.Table(DB_TABLENAME).AutoMigrate(&dbRecord{})
}
