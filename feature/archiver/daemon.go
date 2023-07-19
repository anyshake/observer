package archiver

import (
	"fmt"
	"os"

	"com.geophone.observer/common/postgres"
	"com.geophone.observer/feature"
	"com.geophone.observer/handler"
)

func (a *Archiver) Start(options *feature.FeatureOptions) {
	if !options.Config.Archiver.Enable {
		options.OnStop(MODULE, options, "service is disabled")
		return
	}

	options.OnStart(MODULE, options, "service has started")
	pdb, err := postgres.Open(
		options.Config.Archiver.Host,
		options.Config.Archiver.Port,
		options.Config.Archiver.Username,
		options.Config.Archiver.Password,
		options.Config.Archiver.Database,
	)
	if err != nil {
		options.OnError(MODULE, options, err)
		os.Exit(1)
	}

	err = postgres.Init(pdb)
	if err != nil {
		options.OnError(MODULE, options, err)
		os.Exit(1)
	}
	options.Database = pdb

	// Archive when new message arrived
	handler.OnMessage(&options.Status.Geophone,
		func(gp *handler.Geophone) error {
			var (
				ts  = gp.TS
				ehz = gp.EHZ
				ehe = gp.EHE
				ehn = gp.EHN
			)
			err := postgres.Insert(pdb, ts, ehz, ehe, ehn)
			if err != nil {
				options.OnError(MODULE, options, err)
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
					options.OnError(MODULE, options, err)
					return err
				}
				options.Database = pdb
			}

			options.OnReady(MODULE, options)
			return nil
		},
	)

	err = fmt.Errorf("service exited with a error")
	options.OnError(MODULE, options, err)
}
