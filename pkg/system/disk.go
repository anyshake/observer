package system

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiskPercent() (float64, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("failed to get current executable path: %w", err)
	}
	diskObj, err := disk.Usage(cwd)
	if err != nil {
		return 0, err
	}

	return diskObj.UsedPercent, nil
}
