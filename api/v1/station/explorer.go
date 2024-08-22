package station

import (
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/timesource"
)

func (e *explorerInfo) get(timeSource *timesource.Source, explorerDeps *explorer.ExplorerDependency) error {
	e.DeviceId = explorerDeps.Config.DeviceId
	e.Elevation = explorerDeps.Config.Elevation
	e.Errors = explorerDeps.Health.Errors
	e.Received = explorerDeps.Health.Received
	e.SampleRate = explorerDeps.Health.SampleRate

	currentTime, err := timeSource.Get()
	if err != nil {
		return err
	}
	e.Elapsed = int64(currentTime.Sub(explorerDeps.Health.StartTime).Seconds())

	e.Latitude = float64(int(explorerDeps.Config.Latitude*1000)) / 1000
	e.Longitude = float64(int(explorerDeps.Config.Longitude*1000)) / 1000

	return nil
}
