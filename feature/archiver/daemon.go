package archiver

import (
	"fmt"
	"os"

	"github.com/bclswl0827/observer/driver/postgres"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
)

func (a *Archiver) Start(options *feature.FeatureOptions) {
	if !options.Config.Archiver.Enable {
		a.OnStop(options, "service is disabled")
		return
	}

	a.OnStart(options, "service has started")
	pdb, err := postgres.Open(
		options.Config.Archiver.Host,
		options.Config.Archiver.Port,
		options.Config.Archiver.Username,
		options.Config.Archiver.Password,
		options.Config.Archiver.Database,
	)
	if err != nil {
		a.OnError(options, err)
		os.Exit(1)
	}

	err = postgres.Init(pdb)
	if err != nil {
		a.OnError(options, err)
		os.Exit(1)
	}
	options.Database = pdb

	// Archive when new message arrived
	publisher.Subscribe(
		&options.Status.Geophone,
		func(gp *publisher.Geophone) error {
			var (
				ts  = gp.TS
				ehz = gp.EHZ
				ehe = gp.EHE
				ehn = gp.EHN
			)
			err := postgres.Insert(pdb, ts, ehz, ehe, ehn)
			if err != nil {
				a.OnError(options, err)
				postgres.Close(pdb)

				// Reconnect to PostgreSQL
				pdb, err := postgres.Open(
					options.Config.Archiver.Host,
					options.Config.Archiver.Port,
					options.Config.Archiver.Username,
					options.Config.Archiver.Password,
					options.Config.Archiver.Database,
				)
				if err != nil {
					a.OnError(options, err)
					return err
				}
				options.Database = pdb
			}

			a.OnReady(options)
			return nil
		},
	)

	err = fmt.Errorf("service exited with a error")
	a.OnError(options, err)
}
