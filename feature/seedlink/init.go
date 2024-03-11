package seedlink

import (
	"fmt"
	"log"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/dgraph-io/badger/v4"
	c "github.com/ostafen/clover/v2"
	badgerstore "github.com/ostafen/clover/v2/store/badger"
)

func (s *SeedLink) InitClient(slClient *seedlink.SeedLinkClient) {
	slClient.Streaming = false
}

func (s *SeedLink) InitGlobal(slGlobal *seedlink.SeedLinkGlobal, currentTime time.Time, station, network, location string, bufferDuration int) error {
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

	// Create buffer store
	const collectionName = "observer"
	badgerstoreOptions := badger.DefaultOptions("").WithInMemory(true)
	badgerstoreOptions.Logger = nil
	store, err := badgerstore.OpenWithOptions(badgerstoreOptions)
	if err != nil {
		return err
	}
	db, err := c.OpenWithStore(store)
	if err != nil {
		return err
	}

	// Create collection
	collectionExists, err := db.HasCollection(collectionName)
	if err != nil {
		log.Fatalln(err)
	}
	if !collectionExists {
		db.CreateCollection(collectionName)
	}

	// Initialize ring buffer
	duration := time.Duration(bufferDuration) * time.Second
	slGlobal.SeedLinkBuffer = seedlink.SeedLinkBuffer{Collection: collectionName, Duration: duration, Database: db}

	return nil
}
