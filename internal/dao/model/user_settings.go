package model

import (
	"fmt"

	"github.com/anyshake/observer/internal/dao"
	"gorm.io/gorm"
)

type UserSettings struct {
	dao.BaseTable
	Namespace string `gorm:"column:namespace;index;not null"`
	Key       string `gorm:"column:key;index;not null"`
	Value     []byte `gorm:"column:value;not null"`
	Type      string `gorm:"column:type;not null"`
	Version   int    `gorm:"column:version;not null"`
	UpdatedAt int64  `gorm:"column:update_at;autoUpdateTime:milli;<-:update"`
}

func (t *UserSettings) GetModel() any {
	return &UserSettings{}
}

func (t *UserSettings) GetName(tablePrefix string) string {
	// return "user_settings"
	return fmt.Sprintf("%s%s", tablePrefix, "user_settings")
}

func (t *UserSettings) UseAutoMigrate() bool {
	return true
}

func (t *UserSettings) AddPlugins(dbObj *gorm.DB, tablePrefix string) ([]gorm.Plugin, error) {
	return nil, nil
}
