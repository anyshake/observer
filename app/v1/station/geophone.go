package station

import "github.com/anyshake/observer/config"

func getGeophone(conf *config.Conf) geophoneModel {
	return geophoneModel{
		EHZ: conf.Geophone.EHZ.Sensitivity,
		EHE: conf.Geophone.EHE.Sensitivity,
		EHN: conf.Geophone.EHN.Sensitivity,
	}
}
