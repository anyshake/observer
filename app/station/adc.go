package station

import "github.com/bclswl0827/observer/config"

func getADC(conf *config.Conf) ADC {
	return ADC{
		Resolution: conf.ADC.Resolution,
		FullScale:  conf.ADC.FullScale,
	}
}
