package dao

import (
	"github.com/anyshake/observer/publisher"
	"gorm.io/gorm"
)

func Query(db *gorm.DB, start, end int64) ([]publisher.Geophone, error) {
	var records []dbRecord
	err := db.Table(DB_TABLENAME).Select("ts, ehz, ehe, ehn").Where("ts >= ? AND ts <= ?", start, end).Scan(&records).Error

	var result []publisher.Geophone
	for _, v := range records {
		result = append(result, publisher.Geophone{
			TS: v.TS, EHZ: v.EHZ, EHE: v.EHE, EHN: v.EHN,
		})
	}

	return result, err
}
