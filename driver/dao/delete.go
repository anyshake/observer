package dao

import "gorm.io/gorm"

func Delete(db *gorm.DB, start, end int64) error {
	var records []dbRecord
	return db.Table(DB_TABLENAME).Where("ts >= ? AND ts <= ?", start, end).Delete(&records).Error
}
