package station

import "com.geophone.observer/features/collector"

type Station struct{}

type Memory struct {
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type Disk struct {
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type OS struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Distro   string `json:"distro"`
	Hostname string `json:"hostname"`
}

type CPU struct {
	Model   string  `json:"model"`
	Percent float64 `json:"percent"`
}

type System struct {
	Uptime int64            `json:"uptime"`
	Memory Memory           `json:"memory"`
	Disk   Disk             `json:"disk"`
	OS     OS               `json:"os"`
	CPU    CPU              `json:"cpu"`
	Status collector.Status `json:"status"`
}
