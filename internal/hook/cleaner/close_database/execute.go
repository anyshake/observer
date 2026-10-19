package close_database

import (
	"context"

	"github.com/anyshake/observer/pkg/logger"
)

func (d *CloseDatabaseCleanerImpl) Execute(context.Context) error {
	if d.DAO != nil {
		logger.GetLogger(d.GetName()).Infoln("closing connection to database")
		defer logger.GetLogger(d.GetName()).Infoln("database connection has been closed")
		return d.DAO.Close()
	}

	return nil
}
