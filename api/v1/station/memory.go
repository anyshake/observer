package station

import "github.com/shirou/gopsutil/v4/mem"

func (m *memoryInfo) get() error {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	m.Total = vmStat.Total
	m.Free = vmStat.Free
	m.Used = vmStat.Used
	m.Percent = vmStat.UsedPercent
	return nil
}
