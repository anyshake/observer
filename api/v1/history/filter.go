package history

import (
	"fmt"
	"time"

	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/anyshake/observer/drivers/explorer"
)

func (h *History) filterHistory(startTime, endTime int64, maxDuration time.Duration, resolver *v1.Resolver) ([]explorer.ExplorerData, error) {
	if endTime-startTime > maxDuration.Milliseconds() {
		return nil, fmt.Errorf("duration is too large")
	}

	var (
		adcCountModel tables.AdcCount
		adcCountData  []tables.AdcCount
	)
	err := resolver.Database.
		Table(adcCountModel.GetName()).
		Where("timestamp >= ? AND timestamp <= ?", startTime, endTime).
		Order("timestamp ASC").
		Find(&adcCountData).
		Error
	if err != nil {
		return nil, err
	}

	var explorerData []explorer.ExplorerData
	for _, record := range adcCountData {
		explorerData = append(explorerData, explorer.ExplorerData{
			Timestamp:  record.Timestamp,
			SampleRate: record.SampleRate,
			Z_Axis:     record.Z_Axis,
			E_Axis:     record.E_Axis,
			N_Axis:     record.N_Axis,
		})
	}

	if len(explorerData) == 0 {
		return nil, fmt.Errorf("no data available for the given time range")
	}
	return explorerData, nil
}
