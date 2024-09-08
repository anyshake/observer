package services

import (
	"context"
	"sync"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/utils/timesource"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Options struct {
	Config      *config.Config
	Dependency  *dig.Container
	Database    *gorm.DB
	TimeSource  *timesource.Source
	CancelToken context.Context
}

type Service interface {
	Start(*Options, *sync.WaitGroup)
	GetServiceName() string
}
