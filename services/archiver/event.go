package archiver

import (
	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/logger"
)

func (a *ArchiverService) handleExplorerEvent(data *explorer.ExplorerData) {
	var adcCountModel tables.AdcCount

	a.recordBuffer[len(a.recordBuffer)-a.insertCountDown] = *data
	a.insertCountDown--
	a.cleanupCountDown--

	if a.insertCountDown == 0 {
		records := make([]tables.AdcCount, len(a.recordBuffer))
		for i := 0; i < len(a.recordBuffer); i++ {
			records[i] = tables.AdcCount{
				Timestamp:  a.recordBuffer[i].Timestamp,
				SampleRate: a.recordBuffer[i].SampleRate,
				Z_Axis:     a.recordBuffer[i].Z_Axis,
				E_Axis:     a.recordBuffer[i].E_Axis,
				N_Axis:     a.recordBuffer[i].N_Axis,
			}
		}
		err := a.databaseConn.
			Table(adcCountModel.GetName()).
			Create(records).
			Error
		if err != nil {
			logger.GetLogger(a.GetServiceName()).Warnln(err)
		} else {
			logger.GetLogger(a.GetServiceName()).Infof("%d record(s) has been inserted to database", len(a.recordBuffer))
		}
		a.insertCountDown = INSERT_COUNTDOWN
	}

	if a.cleanupCountDown == 0 {
		err := a.databaseConn.
			Table(adcCountModel.GetName()).
			Where("timestamp < ?", data.Timestamp-int64(a.lifeCycle*86400*1000)).
			Delete(&tables.AdcCount{}).
			Error
		if err != nil {
			logger.GetLogger(a.GetServiceName()).Warnln(err)
		}
		a.cleanupCountDown = CLEANUP_COUNTDOWN
	}
}
