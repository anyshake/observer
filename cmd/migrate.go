package main

import (
	"github.com/anyshake/observer/drivers/dao"
	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func migrate(databaseConn *gorm.DB) error {
	appTables := []dao.Table{
		&tables.AdcCount{},
		&tables.SysUser{},
	}
	for _, table := range appTables {
		err := dao.Migrate(databaseConn, table)
		if err != nil {
			return err
		}
	}

	return nil
}
