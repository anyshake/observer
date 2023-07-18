package station

import "com.geophone.observer/config"

func getLocation(conf *config.Conf) Location {
	return Location{
		Latitude:  conf.Station.Latitude,
		Longitude: conf.Station.Longitude,
		Altitude:  conf.Station.Altitude,
	}
}
