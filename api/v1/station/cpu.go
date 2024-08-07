package station

import (
	"github.com/shirou/gopsutil/cpu"
)

func (c *cpuInfo) get() error {
	eachCoreInfo, _ := cpu.Info()
	eachPercent, _ := cpu.Percent(0, true)

	if len(eachCoreInfo) != 0 {
		totalPercent := 0.0
		currentModel := eachCoreInfo[0].ModelName
		for _, v := range eachPercent {
			totalPercent += v
		}
		c.Model = currentModel
		c.Percent = totalPercent / float64(len(eachPercent))
	} else {
		c.Model = "Unknown"
		c.Percent = 0
	}

	return nil
}
