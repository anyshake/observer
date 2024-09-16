package station

import (
	"os"

	"github.com/shirou/gopsutil/v4/disk"
)

func (d *diskInfo) get() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	usage, err := disk.Usage(cwd)
	if err != nil {
		return err
	}

	d.Total = usage.Total
	d.Free = usage.Free
	d.Used = usage.Used
	d.Percent = usage.UsedPercent
	return nil
}
