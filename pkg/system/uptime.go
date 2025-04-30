package system

import (
	"github.com/shirou/gopsutil/v4/host"
)

func GetOsUptime() (int64, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return 0, err
	}

	return int64(uptime * 1000), nil
}
