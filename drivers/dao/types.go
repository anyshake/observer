package dao

import (
	"time"

	"gorm.io/gorm"
)

const TIMEOUT_THRESHOLD = 5 * time.Second

type engine interface {
	match(engine string) bool
	open(host string, port int, username, password, database string, timeout time.Duration) (*gorm.DB, error)
}

type BaseModel struct {
	// This field is the primary key of the table
	PrimaryKey uint64 `gorm:"primary_key"`
	CreatedAt  int64  `gorm:"column:created_at;autoUpdateTime:milli;<-:create"`
}

type ITable[T any] interface {
	GetModel() T
	GetName() string
}
