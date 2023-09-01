package feature

import (
	"database/sql"

	"github.com/bclswl0827/observer/config"
	"github.com/bclswl0827/observer/publisher"
)

type Feature interface {
	Start(*FeatureOptions)
	OnStart(*FeatureOptions, ...any)
	OnStop(*FeatureOptions, ...any)
	OnReady(*FeatureOptions, ...any)
	OnError(*FeatureOptions, error)
}

type FeatureOptions struct {
	Database *sql.DB
	Config   *config.Conf
	Status   *publisher.Status
}
