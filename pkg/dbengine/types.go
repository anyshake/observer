package dbengine

import (
	"time"

	"gorm.io/gorm"
)

type IEngine interface {
	Open(address, username, password, database, prefix string, timeout time.Duration) (*gorm.DB, error)
}
