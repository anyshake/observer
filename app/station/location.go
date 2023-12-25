package station

import "github.com/anyshake/observer/config"

func getLocation(conf *config.Conf) Location {
	return Location{
		Latitude:  conf.Station.Latitude,
		Longitude: conf.Station.Longitude,
		Elevation: conf.Station.Elevation,
	}
}
