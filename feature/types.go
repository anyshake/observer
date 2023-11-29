package feature

import (
	"sync"

	"github.com/bclswl0827/observer/config"
	"github.com/bclswl0827/observer/publisher"
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
