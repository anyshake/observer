package dao

import (
	"fmt"
)

func (d *DAO) Migrate(tables ...ITable) error {
	if d.Database == nil {
		return fmt.Errorf("database is not opened")
	}

	for _, table := range tables {
		tableRecord := table.GetModel()
		tableName := table.GetName(d.prefix)

		if table.UseAutoMigrate() {
			if err := d.Database.AutoMigrate(tableRecord); err != nil {
				return fmt.Errorf("failed to auto migrate %s table: %w", tableName, err)
			}
		}

		plugins, err := table.AddPlugins(d.Database, d.prefix)
		if err != nil {
			return fmt.Errorf("failed to add plugins to %s table: %w", tableName, err)
		}

		for _, plugin := range plugins {
			if err := d.Database.Use(plugin); err != nil {
				return fmt.Errorf("failed to register plugin %s: %w", plugin.Name(), err)
			}
		}
	}

	return nil
}
