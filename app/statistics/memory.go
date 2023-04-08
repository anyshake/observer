package statistics

import "github.com/shirou/gopsutil/mem"

func GetMemory() Memory {
	vmStat, _ := mem.VirtualMemory()

	return Memory{
		Total:   vmStat.Total,
		Free:    vmStat.Free,
		Used:    vmStat.Used,
		Percent: vmStat.UsedPercent,
	}
}
