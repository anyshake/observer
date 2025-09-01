package metadata

import "time"

const (
	SEISCOMP_TEMPLATE   = "seiscomp.tpl"
	STATIONXML_TEMPLATE = "stationxml.tpl"
)

type Options struct {
	StartTime  time.Time
	SampleRate int

	Latitude  float64
	Longitude float64
	Elevation float64

	NetworkCode  string
	StationCode  string
	LocationCode string
	ChannelCodes []string

	StationPlace       string
	StationCountry     string
	StationAffiliation string
	StationDescription string
}
