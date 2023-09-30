package station

import (
	"github.com/mackerelio/go-osstat/uptime"
)

func getUptime() int64 {
	up, err := uptime.Get()
	if err != nil {
		return -1
	}

	return int64(up.Seconds())
}
