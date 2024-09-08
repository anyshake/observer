package database

import (
	"github.com/anyshake/observer/cleaners"
	"github.com/anyshake/observer/drivers/dao"
	"github.com/anyshake/observer/utils/logger"
)

func (d *DatabaseCleanerTask) Execute(options *cleaners.Options) {
	logger.GetLogger(d.GetTaskName()).Info("closing connection to database")
	dao.Close(options.Database)
}
