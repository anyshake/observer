package station

import (
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/duration"
)

func getSystem(options *feature.FeatureOptions) System {
	_, ts := duration.Timestamp(options.Status.System.Offset)
	return System{
		Status:    options.Status.System,
		Location:  getLocation(options.Config),
		Geophone:  getGeophone(options.Config),
		Station:   getStation(options.Config),
		ADC:       getADC(options.Config),
		Memory:    getMemory(),
		Uptime:    getUptime(),
		Disk:      getDisk(),
		CPU:       getCPU(),
		OS:        getOS(),
		Timestamp: ts,
	}
}
