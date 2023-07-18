package station

import "github.com/shirou/gopsutil/mem"

func getMemory() Memory {
	vmStat, _ := mem.VirtualMemory()

	return Memory{
		Total:   vmStat.Total,
		Free:    vmStat.Free,
		Used:    vmStat.Used,
		Percent: vmStat.UsedPercent,
	}
}
