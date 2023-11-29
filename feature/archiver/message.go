package archiver

import (
	"github.com/bclswl0827/observer/driver/dao"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
	"gorm.io/gorm"
)

func (a *Archiver) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions, pdb *gorm.DB) error {
	err := dao.Insert(pdb, gp)
	if err != nil {
		a.OnError(options, err)
		dao.Close(pdb)
		return err
	}

	a.OnReady(options)
	return nil
}
