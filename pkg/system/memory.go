package system

import (
	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryPercent() (float64, error) {
	memUsage, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return memUsage.UsedPercent, nil
}
