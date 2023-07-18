package station

import "github.com/shirou/gopsutil/cpu"

func getCPU() CPU {
	info, _ := cpu.Info()
	percent, _ := cpu.Percent(0, true)

	cpus := make([]CPU, len(info))
	for i, cpuInfo := range info {
		cpus[i] = CPU{
			Model:   cpuInfo.ModelName,
			Percent: percent[i],
		}
	}

	return cpus[0]
}
