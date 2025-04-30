package system

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCpuModel() (string, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return "", err
	}

	if len(cpuInfo) == 0 || cpuInfo[0].ModelName == "" {
		return "null", nil
	}
	return cpuInfo[0].ModelName, nil
}

func GetCpuPercent() (float64, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		return 0, err
	}

	return cpuPercent[0], nil
}
