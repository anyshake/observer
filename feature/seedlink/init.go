package seedlink

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
)

func (s *SeedLink) InitClient(slClient *seedlink.SeedLinkClient) {
	slClient.StreamMode = false
}

func (s *SeedLink) InitGlobal(slGlobal *seedlink.SeedLinkGlobal, currentTime time.Time, station, network, location, bufferFile string, bufferSize int) error {
	var (
		streamEndTimeString = "9999-12-31 23:59:59"
		currentTimeString   = currentTime.Format("2006-01-02 15:04:01")
	)

	// Initialize SeedLink global states
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
	slGlobal.Streams = []seedlink.SeedLinkStream{
		{BeginTime: currentTimeString, EndTime: streamEndTimeString, SeedName: "EHZ", Location: location, Type: "D", Station: station},
		{BeginTime: currentTimeString, EndTime: streamEndTimeString, SeedName: "EHE", Location: location, Type: "D", Station: station},
		{BeginTime: currentTimeString, EndTime: streamEndTimeString, SeedName: "EHN", Location: location, Type: "D", Station: station},
	}
	slGlobal.Stations = []seedlink.SeedLinkStation{
		{BeginSequence: "000000", EndSequence: "FFFFFF", Station: station, Network: network, Description: fmt.Sprintf("%s station", network)},
	}

	// Check buffer file existence
	stat, err := os.Stat(bufferFile)
	if err != nil {
		_, err = os.Create(bufferFile)
		if err != nil {
			return err
		}
	} else if stat.IsDir() {
		return fmt.Errorf("buffer file is a directory")
	}

	// Read buffer file
	file, err := os.Open(bufferFile)
	if err != nil {
		return err
	}
	defer file.Close()
	var (
		records []publisher.Geophone
		decoder = json.NewDecoder(file)
	)
	decoder.Decode(&records)

	// Initialize ring buffer
	if len(records) >= bufferSize {
		slGlobal.SeedLinkBuffer.Data = records[len(records)-bufferSize:]
	} else {
		slGlobal.SeedLinkBuffer.Data = records
	}
	slGlobal.SeedLinkBuffer = seedlink.SeedLinkBuffer{Size: bufferSize, File: bufferFile}

	return nil
}
