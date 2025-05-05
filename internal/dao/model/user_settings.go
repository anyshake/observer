package model

import (
	"fmt"

	"github.com/anyshake/observer/internal/dao"
	"gorm.io/gorm"
)

type UserSettings struct {
	dao.BaseTable
	Namespace   string `gorm:"column:namespace;index;not null"`
	ConfigKey   string `gorm:"column:config_key;index;not null"`
	ConfigValue []byte `gorm:"column:config_val;not null"`
	ConfigType  string `gorm:"column:config_type;not null"`
	Version     int    `gorm:"column:version;not null"`
	UpdatedAt   int64  `gorm:"column:update_at;autoUpdateTime:milli;<-:update"`
}

func (t *UserSettings) GetModel() any {
	return &UserSettings{}
}

func (t *UserSettings) GetName(tablePrefix string) string {
	return fmt.Sprintf("%s%s", tablePrefix, "user_settings")
}

func (t *UserSettings) UseAutoMigrate() bool {
	return true
}

func (t *UserSettings) AddPlugins(dbObj *gorm.DB, tablePrefix string) ([]gorm.Plugin, error) {
	return nil, nil
}
