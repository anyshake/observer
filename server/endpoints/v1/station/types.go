package station

import (
	"github.com/anyshake/observer/config"
)

type Station struct{}

type explorerInfo struct {
	Elapsed    int64   `json:"elapsed"`
	Errors     int64   `json:"errors"`
	Received   int64   `json:"received"`
	SampleRate int     `json:"sample_rate"`
	DeviceId   string  `json:"device_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Elevation  float64 `json:"elevation"`
}

type memoryInfo struct {
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type diskInfo struct {
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type osInfo struct {
	Uptime    int64  `json:"uptime"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Distro    string `json:"distro"`
	Hostname  string `json:"hostname"`
	Timestamp int64  `json:"timestamp"`
}

type cpuInfo struct {
	Model   string  `json:"model"`
	Percent float64 `json:"percent"`
}

type stationInfo struct {
	Explorer explorerInfo   `json:"explorer"`
	Station  config.Station `json:"station"`
	Stream   config.Stream  `json:"stream"`
	Sensor   config.Sensor  `json:"sensor"`
	CPU      cpuInfo        `json:"cpu"`
	Disk     diskInfo       `json:"disk"`
	Memory   memoryInfo     `json:"memory"`
	OS       osInfo         `json:"os"`
}
