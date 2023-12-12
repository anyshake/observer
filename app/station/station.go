package station

import (
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/duration"
)

func getStation(options *feature.FeatureOptions) System {
	_, ts := duration.Timestamp(options.Status.System.Offset)
	return System{
		Timestamp: ts,
		Status:    options.Status.System,
		Station:   options.Config.Station.Name,
		UUID:      options.Config.Station.UUID,
		Location:  getLocation(options.Config),
		Geophone:  getGeophone(options.Config),
		ADC:       getADC(options.Config),
		OS:        getOS(),
		CPU:       getCPU(),
		Disk:      getDisk(),
		Memory:    getMemory(),
		Uptime:    getUptime(),
	}
}
