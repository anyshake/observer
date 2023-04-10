package station

import "com.geophone.observer/features/collector"

func GetStation(uuid, station string, status collector.Status) System {
	return System{
		UUID:    uuid,
		Station: station,
		CPU:     GetCPU(),
		OS:      GetOS(),
		Disk:    GetDisk(),
		Memory:  GetMemory(),
		Uptime:  GetUptime(),
		Status:  status,
	}
}
