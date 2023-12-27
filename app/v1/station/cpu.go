package station

import (
	"github.com/shirou/gopsutil/cpu"
)

func getCPU() cpuModel {
	eachCoreInfo, _ := cpu.Info()
	eachPercent, _ := cpu.Percent(0, true)

	if len(eachCoreInfo) != 0 {
		totalPercent := 0.0
		currentModel := eachCoreInfo[0].ModelName
		for _, v := range eachPercent {
			totalPercent += v
		}
		return cpuModel{
			Model:   currentModel,
			Percent: totalPercent / float64(len(eachPercent)),
		}
	} else {
		return cpuModel{
			Model:   "Unknown",
			Percent: 0,
		}
	}
}
