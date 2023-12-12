package history

import (
	"strings"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/utils/text"
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
