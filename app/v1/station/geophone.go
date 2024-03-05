package station

import "github.com/anyshake/observer/config"

func getGeophone(conf *config.Conf) geophoneModel {
	return geophoneModel{
		Sensitivity: conf.Geophone.Sensitivity,
		Frequency:   conf.Geophone.Frequency,
	}
}
