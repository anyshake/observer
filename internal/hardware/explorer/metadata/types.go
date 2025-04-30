package metadata

import "time"

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

type IMetadata interface {
	SeisCompXML() string
	StationXML() string
}
