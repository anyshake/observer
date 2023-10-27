package miniseed

import (
	"fmt"
	"sync"

	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
)

func (m *MiniSEED) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	if !options.Config.MiniSEED.Enable {
		m.OnStop(options, "service is disabled")
		return
	}

	// Get MiniSEED info & options
	var (
		basePath  = options.Config.MiniSEED.Path
		station   = options.Config.MiniSEED.Station
		network   = options.Config.MiniSEED.Network
		lifeCycle = options.Config.MiniSEED.LifeCycle
	)

	// Start cleanup routine if life cycle bigger than 0
	if lifeCycle > 0 {
		go m.handleCleanup(basePath, station, network, lifeCycle)
	}

	// Init sequence number
	var seqNumber int
	m.OnStart(options, "service has started")

	// Append and write when new message arrived
	publisher.Subscribe(
		&options.Status.Geophone,
		func(gp *publisher.Geophone) error {
			return m.handleMessage(gp, options, seqNumber, basePath, station, network)
		},
	)

	err := fmt.Errorf("service exited with an error")
	m.OnError(options, err)
}
