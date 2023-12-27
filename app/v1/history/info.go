package history

import (
	"strings"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/utils/text"
)

func getNetwork(conf *config.Conf) string {
	network := strings.ToUpper(conf.Station.Network)
	return text.TruncateString(network, 2)
}

func getStation(conf *config.Conf) string {
	station := strings.ToUpper(conf.Station.Station)
	return text.TruncateString(station, 5)
}

func getLocation(conf *config.Conf) string {
	location := strings.ToUpper(conf.Station.Location)
	return text.TruncateString(location, 2)
}
