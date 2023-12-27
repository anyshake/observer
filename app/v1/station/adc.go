package station

import "github.com/anyshake/observer/config"

func getADC(conf *config.Conf) adcModel {
	return adcModel{
		Resolution: conf.ADC.Resolution,
		FullScale:  conf.ADC.FullScale,
	}
}
