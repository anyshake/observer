package miniseed

import (
	"fmt"
	"time"
)

func getFilePath(basePath, station, network, location string, timestamp time.Time) string {
	// e.g. /path/to/miniseed/AS.SHAKE.00.D.2023.208.mseed
	return fmt.Sprintf("%s/%s.%s.%s.D.%s.%s.mseed",
		basePath,
		network, station, location,
		timestamp.Format("2006"),
		timestamp.Format("002"),
	)
}
