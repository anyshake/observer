package station

import "com.geophone.observer/config"

func getADC(conf *config.Conf) ADC {
	return ADC{
		Resolution: conf.ADC.Resolution,
		FullScale:  conf.ADC.FullScale,
	}
}
