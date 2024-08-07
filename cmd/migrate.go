package main

import (
	"github.com/anyshake/observer/drivers/dao"
	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func migrate(databaseConn *gorm.DB) error {
	err := dao.Migrate(databaseConn, tables.AdcCount{})
	if err != nil {
		return err
	}

	return nil
}
