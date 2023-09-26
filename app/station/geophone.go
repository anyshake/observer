package station

import "github.com/bclswl0827/observer/config"

func getGeophone(conf *config.Conf) Geophone {
	return Geophone{
		EHZ: conf.Geophone.EHZ.Sensitivity,
		EHE: conf.Geophone.EHE.Sensitivity,
		EHN: conf.Geophone.EHN.Sensitivity,
	}
}
