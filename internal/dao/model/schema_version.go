package model

import (
	"fmt"

	"github.com/anyshake/observer/internal/dao"
	"gorm.io/gorm"
)

type SchemaVersion struct {
	dao.BaseTable

	Version     int    `gorm:"column:version;not null"`
	Description string `gorm:"column:description;not null"`
	AppliedAt   int64  `gorm:"column:applied_at;not null"`
}

func (t *SchemaVersion) GetModel() any {
	return &SchemaVersion{}
}

func (t *SchemaVersion) GetName(tablePrefix string) string {
	return fmt.Sprintf("%s%s", tablePrefix, "schema_version")
}

func (t *SchemaVersion) UseAutoMigrate() bool {
	return true
}

func (t *SchemaVersion) AddPlugins(dbObj *gorm.DB, tablePrefix string) ([]gorm.Plugin, error) {
	return nil, nil
}
