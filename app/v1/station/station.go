package station

import "github.com/anyshake/observer/config"

func getStation(conf *config.Conf) stationModel {
	return stationModel{
		UUID:     conf.Station.UUID,
		Name:     conf.Station.Name,
		Station:  conf.Station.Station,
		Network:  conf.Station.Network,
		Location: conf.Station.Location,
	}
}
