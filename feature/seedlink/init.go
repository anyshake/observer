package seedlink

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
)

func (s *SeedLink) InitClient(slClient *seedlink.SeedLinkClient) {
	slClient.StreamMode = false
}

func (s *SeedLink) InitGlobal(slGlobal *seedlink.SeedLinkGlobal, currentTime time.Time, station, network, location string) {
	var (
		streamEndTimeString = "9999-12-31 23:59:59"
		currentTimeString   = currentTime.Format("2006-01-02 15:04:01")
	)

	slGlobal.SeedLinkState = seedlink.SeedLinkState{
		Organization: seedlink.ORGANIZATION,
		StartTime:    currentTimeString,
		Software:     seedlink.RELEASE,
	}
	slGlobal.Capabilities = []seedlink.SeedLinkCapability{
		{Name: "info:all"}, {Name: "info:gaps"}, {Name: "info:streams"},
		{Name: "dialup"}, {Name: "info:id"}, {Name: "multistation"},
		{Name: "window-extraction"}, {Name: "info:connections"},
		{Name: "info:capabilities"}, {Name: "info:stations"},
	}
	// Station field are not used by SeedLink, but are required by the protocol to differentiate between stations
	slGlobal.SeedLinkBuffer = seedlink.SeedLinkBuffer{Size: SEEDLINK_BUFFERSIZE}
	slGlobal.Streams = []seedlink.SeedLinkStream{
		{BeginTime: currentTimeString, EndTime: streamEndTimeString, SeedName: "EHZ", Location: location, Type: "D", Station: station},
		{BeginTime: currentTimeString, EndTime: streamEndTimeString, SeedName: "EHE", Location: location, Type: "D", Station: station},
		{BeginTime: currentTimeString, EndTime: streamEndTimeString, SeedName: "EHN", Location: location, Type: "D", Station: station},
	}
	slGlobal.Stations = []seedlink.SeedLinkStation{
		{BeginSequence: "000000", EndSequence: "FFFFFF", Station: station, Network: network, Description: fmt.Sprintf("%s station", network)},
	}
}
