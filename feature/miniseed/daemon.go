package miniseed

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/duration"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
	"github.com/bclswl0827/mseedio"
	"github.com/fatih/color"
)

func (m *MiniSEED) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	if !options.Config.MiniSEED.Enable {
		m.OnStop(options, "service is disabled")
		return
	}

	// Get MiniSEED info & options
	var (
		basePath  = options.Config.MiniSEED.Path
		lifeCycle = options.Config.MiniSEED.LifeCycle
		station   = text.TruncateString(options.Config.Station.Station, 5)
		network   = text.TruncateString(options.Config.Station.Network, 2)
		location  = text.TruncateString(options.Config.Station.Location, 2)
	)

	// Start cleanup routine if life cycle bigger than 0
	if lifeCycle > 0 {
		go m.handleCleanup(basePath, station, network, lifeCycle)
	}

	// Wait for time syncing
	for !options.Status.IsReady {
		logger.Print(MODULE, "waiting for time alignment", color.FgYellow, false)
		time.Sleep(50 * time.Millisecond)
	}

	// Get sequence number if file exists
	var seqNumber int64
	currentTime, _ := duration.Timestamp(options.Status.System.Offset)
	filePath := getFilePath(basePath, station, network, location, currentTime)
	_, err := os.Stat(filePath)
	if err == nil {
		// Get last sequence number
		logger.Print(MODULE, "starting from last record", color.FgYellow, false)

		// Read MiniSEED file
		var ms mseedio.MiniSeedData
		err := ms.Read(filePath)
		if err != nil {
			m.OnError(options, err)
			return
		}

		// Get last sequence number
		recordLength := len(ms.Series)
		if recordLength > 0 {
			lastRecord := ms.Series[recordLength-1]
			lastSeqNum := lastRecord.FixedSection.SequenceNumber
			n, err := strconv.Atoi(lastSeqNum)
			if err != nil {
				m.OnError(options, err)
				return
			}
			// Set current sequence number
			seqNumber = int64(n)
		}
	} else {
		// Create new file with sequence number 0
		logger.Print(MODULE, "starting from a new file", color.FgYellow, false)
	}

	// Init MiniSEED archiving buffer
	buffer := &miniSEEDBuffer{
		TimeStamp: currentTime,
		SeqNum:    seqNumber,
		EHZ:       &channelBuffer{},
		EHE:       &channelBuffer{},
		EHN:       &channelBuffer{},
		BasePath:  options.Config.MiniSEED.Path,
		Station:   options.Config.Station.Station,
		Network:   options.Config.Station.Network,
	}
	m.OnStart(options, "service has started")

	// Append and write when new message arrived
	publisher.Subscribe(
		&options.Status.Geophone,
		func(gp *publisher.Geophone) error {
			return m.handleMessage(gp, options, buffer)
		},
	)

	err = fmt.Errorf("service exited with an error")
	m.OnError(options, err)
}
