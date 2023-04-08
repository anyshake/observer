package statistics

func GetStatistics() System {
	return System{
		CPU:    GetCPU(),
		OS:     GetOS(),
		Disk:   GetDisk(),
		Memory: GetMemory(),
		Uptime: GetUptime(),
	}
}
