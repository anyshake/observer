package history

import (
	"strings"

	"com.geophone.observer/config"
	"com.geophone.observer/utils/text"
)

func getNetwork(*config.Conf) string {
	return text.TruncateString(NETWORK, 8)
}

func getStation(conf *config.Conf) string {
	station := strings.ToUpper(conf.Station.UUID)
	return text.TruncateString(station, 8)
}

func getLocation(*config.Conf) string {
	return UNDEFINED
}
