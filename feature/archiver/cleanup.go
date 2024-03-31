package archiver

import (
	"time"

	"github.com/anyshake/observer/driver/dao"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/duration"
	"gorm.io/gorm"
)

func (a *Archiver) handleCleanup(status *publisher.Status, db *gorm.DB, lifeCycle int) {
	for {
		// Wait until system is ready
		if status.ReadyTime.IsZero() {
			time.Sleep(time.Second)
			continue
		}

		// Get start and end time
		currentTime, _ := duration.Timestamp(status.System.Offset)
		endTime := currentTime.Add(-time.Duration(lifeCycle) * time.Hour * 24)

		// Remove expired records
		err := dao.Delete(db, 0, endTime.UnixMilli())
		if err != nil {
			a.OnError(nil, err)
		}

		time.Sleep(time.Hour)
	}
}
