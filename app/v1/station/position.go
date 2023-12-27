package station

import "github.com/anyshake/observer/config"

func getLocation(conf *config.Conf) positionModel {
	return positionModel{
		Latitude:  conf.Station.Latitude,
		Longitude: conf.Station.Longitude,
		Elevation: conf.Station.Elevation,
	}
}
