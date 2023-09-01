package station

import "github.com/bclswl0827/observer/config"

func getGeophone(conf *config.Conf) Geophone {
	return Geophone{
		EHZ: conf.Geophone.EHZ,
		EHE: conf.Geophone.EHE,
		EHN: conf.Geophone.EHN,
	}
}
