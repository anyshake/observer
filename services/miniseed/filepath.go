package miniseed

import (
	"fmt"
	"time"
)

func (m *MiniSeedService) getMiniSEEDFilePath(channelName string, t time.Time) string {
	// e.g. /path/to/miniseed/AS.SHAKE.00.EHZ.D.2023.208.mseed
	return fmt.Sprintf("%s/%s.%s.%s.%s.D.%s.%s.mseed",
		m.basePath,
		m.networkCode, m.stationCode, m.locationCode, channelName,
		t.UTC().Format("2006"),
		t.UTC().Format("002"),
	)
}

func (m *MiniSeedService) getSequenceFilePath() string {
	// e.g. .AS.SHAKE.00.json
	return fmt.Sprintf(
		"%s/.%s.%s.%s.json",
		m.basePath,
		m.networkCode,
		m.stationCode,
		m.locationCode,
	)
}
