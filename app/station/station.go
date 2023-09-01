package station

import (
	"github.com/bclswl0827/observer/feature"
)

func getStation(options *feature.FeatureOptions) System {
	return System{
		Status:   options.Status.System,
		Station:  options.Config.Station.Name,
		UUID:     options.Config.Station.UUID,
		Location: getLocation(options.Config),
		Geophone: getGeophone(options.Config),
		ADC:      getADC(options.Config),
		OS:       getOS(),
		CPU:      getCPU(),
		Disk:     getDisk(),
		Memory:   getMemory(),
		Uptime:   getUptime(),
	}
}
