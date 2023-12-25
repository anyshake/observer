package station

import "github.com/anyshake/observer/publisher"

type Station struct{}

type ADC struct {
	Resolution int     `json:"resolution"`
	FullScale  float64 `json:"fullscale"`
}

type Geophone struct {
	EHZ float64 `json:"ehz"`
	EHE float64 `json:"ehe"`
	EHN float64 `json:"ehn"`
}

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

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
}

type System struct {
	UUID      string            `json:"uuid"`
	Station   string            `json:"station"`
	Timestamp int64             `json:"timestamp"`
	Uptime    int64             `json:"uptime"`
	Memory    Memory            `json:"memory"`
	Disk      Disk              `json:"disk"`
	ADC       ADC               `json:"adc"`
	OS        OS                `json:"os"`
	CPU       CPU               `json:"cpu"`
	Geophone  Geophone          `json:"geophone"`
	Location  Location          `json:"location"`
	Status    *publisher.System `json:"status"`
}
