package dao

import (
	"time"

	"github.com/anyshake/observer/pkg/dbengine"
	"gorm.io/gorm"
)

type DAO struct {
	address  string
	username string
	password string
	prefix   string
	timeout  time.Duration

	driver   dbengine.IEngine
	Database *gorm.DB
}

type ITable interface {
	GetModel() any
	UseAutoMigrate() bool
	GetName(tablePrefix string) string
	AddPlugins(db *gorm.DB, tablePrefix string) ([]gorm.Plugin, error)
}

type BaseTable struct {
	PrimaryKey uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt  int64  `gorm:"column:created_at;autoUpdateTime:milli;<-:create"`
}
