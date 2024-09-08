package station

import (
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/timesource"
)

func (e *explorerInfo) get(timeSource *timesource.Source, explorerDeps *explorer.ExplorerDependency) error {
	e.DeviceId = explorerDeps.Config.GetDeviceId()
	e.Elevation = explorerDeps.Config.GetElevation()
	e.Errors = explorerDeps.Health.GetErrors()
	e.Received = explorerDeps.Health.GetReceived()
	e.SampleRate = explorerDeps.Health.GetSampleRate()

	currentTime := timeSource.Get()
	e.Elapsed = int64(currentTime.Sub(explorerDeps.Health.GetStartTime()).Seconds())

	e.Latitude = float64(int(explorerDeps.Config.GetLatitude()*1000)) / 1000
	e.Longitude = float64(int(explorerDeps.Config.GetLongitude()*1000)) / 1000

	return nil
}
