package archiver

import (
	"github.com/anyshake/observer/drivers/explorer"
	"gorm.io/gorm"
)

const (
	INSERT_COUNTDOWN  = 10
	CLEANUP_COUNTDOWN = 60
)

type ArchiverService struct {
	recordBuffer     [INSERT_COUNTDOWN]explorer.ExplorerData
	insertCountDown  int
	cleanupCountDown int
	lifeCycle        int
	databaseConn     *gorm.DB
}
