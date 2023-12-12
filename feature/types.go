package feature

import (
	"sync"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/publisher"
	"gorm.io/gorm"
)

type Feature interface {
	Run(*FeatureOptions, *sync.WaitGroup)
	OnStart(*FeatureOptions, ...any)
	OnStop(*FeatureOptions, ...any)
	OnReady(*FeatureOptions, ...any)
	OnError(*FeatureOptions, error)
}

type FeatureOptions struct {
	Database *gorm.DB
	Config   *config.Conf
	Status   *publisher.Status
}
