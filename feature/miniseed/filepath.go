package miniseed

import (
	"fmt"
	"time"
)

func getFilePath(basePath, station, network, location, channel string, timestamp time.Time) string {
	// e.g. /path/to/miniseed/AS.SHAKE.00.EHZ.D.2023.208.mseed
	return fmt.Sprintf("%s/%s.%s.%s.%s.D.%s.%s.mseed",
		basePath,
		network, station, location, channel,
		timestamp.Format("2006"),
		timestamp.Format("002"),
	)
}
