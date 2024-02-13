package station

import (
	"os"

	"github.com/shirou/gopsutil/disk"
)

func getDisk() diskModel {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	usage, err := disk.Usage(cwd)
	if err != nil {
		panic(err)
	}

	return diskModel{
		Total:   usage.Total,
		Free:    usage.Free,
		Used:    usage.Used,
		Percent: usage.UsedPercent,
	}
}
