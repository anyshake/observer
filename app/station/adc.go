package station

import "github.com/anyshake/observer/config"

func getADC(conf *config.Conf) ADC {
	return ADC{
		Resolution: conf.ADC.Resolution,
		FullScale:  conf.ADC.FullScale,
	}
}
