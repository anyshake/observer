package cleaners

import (
	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/utils/timesource"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Options struct {
	Config     *config.Config
	Database   *gorm.DB
	Dependency *dig.Container
	TimeSource *timesource.Source
}

type CleanerTask interface {
	Execute(*Options)
	GetTaskName() string
}
