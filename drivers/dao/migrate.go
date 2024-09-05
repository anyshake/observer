package dao

import (
	"gorm.io/gorm"
)

func Migrate(dbObj *gorm.DB, tableImpl Table[any]) error {
	tableRecord := tableImpl.GetModel()
	return dbObj.Table(tableImpl.GetName()).AutoMigrate(&tableRecord)
}
