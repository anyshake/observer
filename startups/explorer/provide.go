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
		return &explorer.ExplorerDependency{
			FallbackTime: options.TimeSource,
			CancelToken:  t.CancelToken,
			Transport:    explorerTransport,
			Config: explorer.ExplorerConfig{
				Latitude:   options.Config.Location.Latitude,
				Longitude:  options.Config.Location.Longitude,
				Elevation:  options.Config.Location.Elevation,
				LegacyMode: options.Config.Explorer.Legacy,
			},
		}
	})
}
