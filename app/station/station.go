package station

import (
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/duration"
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
