package station

import "com.geophone.observer/features/collector"

func GetStation(message collector.Message, status collector.Status) System {
	return System{
		UUID:     message.UUID,
		Station:  message.Station,
		Location: GetLocation(message),
		CPU:      GetCPU(),
		OS:       GetOS(),
		Disk:     GetDisk(),
		Memory:   GetMemory(),
		Uptime:   GetUptime(),
		Status:   status,
	}
}
