package station

import "com.geophone.observer/features/collector"

func Getstation(status *collector.Status) System {
	return System{
		CPU:    GetCPU(),
		OS:     GetOS(),
		Disk:   GetDisk(),
		Memory: GetMemory(),
		Uptime: GetUptime(),
		Status: *status,
	}
}
