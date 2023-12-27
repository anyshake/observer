package station

import "github.com/shirou/gopsutil/mem"

func getMemory() memoryModel {
	vmStat, _ := mem.VirtualMemory()

	return memoryModel{
		Total:   vmStat.Total,
		Free:    vmStat.Free,
		Used:    vmStat.Used,
		Percent: vmStat.UsedPercent,
	}
}
