package station

import "github.com/anyshake/observer/publisher"

type Station struct{}

type adcModel struct {
	Resolution int     `json:"resolution"`
	FullScale  float64 `json:"fullscale"`
}

type geophoneModel struct {
	Sensitivity float64 `json:"sensitivity"`
	Frequency   float64 `json:"frequency"`
}

type memoryModel struct {
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type diskModel struct {
	Total   uint64  `json:"total"`
	Free    uint64  `json:"free"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type osModel struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Distro   string `json:"distro"`
	Hostname string `json:"hostname"`
}

type cpuModel struct {
	Model   string  `json:"model"`
	Percent float64 `json:"percent"`
}

type positionModel struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
}

type stationModel struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Station  string `json:"station"`
	Network  string `json:"network"`
	Location string `json:"location"`
}

type System struct {
	Timestamp int64             `json:"timestamp"`
	Uptime    int64             `json:"uptime"`
	Station   stationModel      `json:"station"`
	Memory    memoryModel       `json:"memory"`
	Disk      diskModel         `json:"disk"`
	ADC       adcModel          `json:"adc"`
	OS        osModel           `json:"os"`
	CPU       cpuModel          `json:"cpu"`
	Geophone  geophoneModel     `json:"geophone"`
	Location  positionModel     `json:"position"`
	Status    *publisher.System `json:"status"`
}
