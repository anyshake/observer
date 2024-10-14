package archiver

import (
	"github.com/anyshake/observer/drivers/explorer"
	"gorm.io/gorm"
)

const (
	RECORDS_INSERT_INTERVAL  = 10
	RECORDS_CLEANUP_INTERVAL = 3600
)

type ArchiverService struct {
	recordBuffer     [RECORDS_INSERT_INTERVAL]explorer.ExplorerData
	insertCountDown  int
	cleanupCountDown int
	lifeCycle        int
	databaseConn     *gorm.DB
}
