package dao

import (
	"github.com/anyshake/observer/publisher"
	"gorm.io/gorm"
)

func Insert(db *gorm.DB, gp *publisher.Geophone) error {
	return db.Table(DB_TABLENAME).Create(&dbRecord{Geophone: *gp}).Error
}
