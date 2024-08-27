package explorer

import (
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/drivers/transport"
	"github.com/anyshake/observer/startups"
	"github.com/anyshake/observer/utils/logger"
	"go.uber.org/dig"
)

func (t *ExplorerStartupTask) Provide(container *dig.Container, options *startups.Options) error {
	// Open AnyShake Explorer device
	explorerDsn := &transport.TransportDependency{
		DSN:    options.Config.Explorer.DSN,
		Engine: options.Config.Explorer.Engine,
	}
	explorerTransport, err := transport.New(explorerDsn)
	if err != nil {
		return err
	}
	err = explorerTransport.Open(explorerDsn)
	if err != nil {
		return err
	}

	logger.GetLogger(t.GetTaskName()).Infoln("device has been opened successfully")
	return container.Provide(func() *explorer.ExplorerDependency {
		var explorerConfig explorer.ExplorerConfig
		explorerConfig.SetLegacyMode(options.Config.Explorer.Legacy)
		explorerConfig.SetLatitude(options.Config.Location.Latitude)
		explorerConfig.SetLongitude(options.Config.Location.Longitude)
		explorerConfig.SetElevation(options.Config.Location.Elevation)

		return &explorer.ExplorerDependency{
			FallbackTime: options.TimeSource,
			CancelToken:  t.CancelToken,
			Transport:    explorerTransport,
			Config:       &explorerConfig,
			Health:       &explorer.ExplorerHealth{},
		}
	})
}
