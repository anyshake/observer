package main

import (
	"github.com/anyshake/observer/drivers/dao"
	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func migrate(databaseConn *gorm.DB) error {
	return dao.Migrate(databaseConn, &tables.AdcCount{})
}
