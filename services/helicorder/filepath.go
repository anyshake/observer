package helicorder

import (
	"fmt"
	"time"
)

func (m *HelicorderService) getFilePath(channelName string, t time.Time) string {
	// e.g. /path/to/helicorder/AS.SHAKE.00.EHZ.2023.208.svg
	return fmt.Sprintf("%s/%s.%s.%s.%s.%s.%s.svg",
		m.basePath,
		m.networkCode, m.stationCode, m.locationCode, channelName,
		t.UTC().Format("2006"),
		t.UTC().Format("002"),
	)
}
