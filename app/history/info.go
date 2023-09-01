package history

import (
	"strings"

	"github.com/bclswl0827/observer/config"
	"github.com/bclswl0827/observer/utils/text"
)

func getNetwork(*config.Conf) string {
	return NETWORK
}

func getStation(conf *config.Conf) string {
	station := strings.ToUpper(conf.Station.UUID)
	return text.TruncateString(station, 8)
}

func getLocation(*config.Conf) string {
	return UNDEFINED
}
