package archiver

import (
	"database/sql"

	"github.com/bclswl0827/observer/driver/postgres"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
)

func (a *Archiver) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions, pdb *sql.DB) error {
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
}
