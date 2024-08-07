package startups

import (
	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/utils/timesource"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Options struct {
	Config     *config.Config
	Database   *gorm.DB
	TimeSource timesource.Source
}

type StartupTask interface {
	Provide(*dig.Container, *Options) error
	Execute(*dig.Container, *Options) error
	GetTaskName() string
}
