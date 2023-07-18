package feature

import (
	"database/sql"

	"com.geophone.observer/config"
	"com.geophone.observer/handler"
)

type Feature interface {
	Start(*FeatureOptions)
}

type FeatureOptions struct {
	Database *sql.DB
	Config   *config.Conf
	Status   *handler.Status
	OnStart  func(module string, options *FeatureOptions, v ...any)
	OnStop   func(module string, options *FeatureOptions, v ...any)
	OnReady  func(module string, options *FeatureOptions, v ...any)
	OnError  func(module string, options *FeatureOptions, err error)
}
