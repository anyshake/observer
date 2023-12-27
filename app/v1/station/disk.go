package station

import "github.com/shirou/gopsutil/disk"

func getDisk() diskModel {
	partitions, err := disk.Partitions(false)
	if err != nil {
		panic(err)
	}

	disks := make([]diskModel, 0)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			panic(err)
		}

		disks = append(disks, diskModel{
			Total:   usage.Total,
			Free:    usage.Free,
			Used:    usage.Used,
			Percent: usage.UsedPercent,
		})
	}

	return disks[0]
}
