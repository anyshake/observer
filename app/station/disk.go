package station

import "github.com/shirou/gopsutil/disk"

func GetDisk() Disk {
	partitions, err := disk.Partitions(false)
	if err != nil {
		panic(err)
	}

	disks := make([]Disk, 0)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			panic(err)
		}

		disks = append(disks, Disk{
			Total:   usage.Total,
			Free:    usage.Free,
			Used:    usage.Used,
			Percent: usage.UsedPercent,
		})
	}

	return disks[0]
}
