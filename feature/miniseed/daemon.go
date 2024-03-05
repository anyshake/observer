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
	for options.Status.ReadyTime.IsZero() {
		logger.Print(MODULE, "waiting for time alignment", color.FgYellow, false)
		time.Sleep(1 * time.Second)
	}

	// Init MiniSEED archiving buffer
	currentTime, _ := duration.Timestamp(options.Status.System.Offset)
	miniSEEDBuffer := &publisher.SegmentBuffer{
		TimeStamp: currentTime,
		ChannelBuffer: map[string]*publisher.ChannelSegmentBuffer{
			"EHZ": {}, "EHE": {}, "EHN": {},
		},
	}

	// Get sequence number if file exists
	for i, v := range miniSEEDBuffer.ChannelBuffer {
		filePath := getFilePath(basePath, station, network, location, i, currentTime)
		_, err := os.Stat(filePath)
		if err == nil {
			// Get last sequence number
			logger.Print(MODULE, fmt.Sprintf("starting %s from last record", i), color.FgYellow, false)

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
				v.SeqNum = int64(n)
			}
		} else {
			// Create new file with sequence number 0
			logger.Print(MODULE, fmt.Sprintf("starting %s from a new file", i), color.FgYellow, false)
		}
	}
	m.OnStart(options, "service has started")

	// Append and write when new message arrived
	expressionForSubscribe := true
	publisher.Subscribe(
		&options.Status.Geophone,
		&expressionForSubscribe,
		func(gp *publisher.Geophone) error {
			return m.handleMessage(gp, options, miniSEEDBuffer)
		},
	)

	err := fmt.Errorf("service exited with an error")
	m.OnError(options, err)
}
