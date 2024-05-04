package seedlink

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
	"github.com/fatih/color"
)

func (s *SeedLink) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	if !options.Config.SeedLink.Enable {
		s.OnStop(options, "service is disabled")
		return
	}

	// Increase wait group counter
	waitGroup.Add(1)
	defer waitGroup.Done()

	// Create TCP server and listen
	host, port := options.Config.SeedLink.Host, options.Config.SeedLink.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		s.OnError(options, err)
		os.Exit(1)
	}
	defer listener.Close()

	// Init SeedLink global state
	var (
		slGlobal         seedlink.SeedLinkGlobal
		station          = text.TruncateString(options.Config.Station.Station, 5)
		network          = text.TruncateString(options.Config.Station.Network, 2)
		location         = text.TruncateString(options.Config.Station.Location, 2)
		bufferDuration   = options.Config.SeedLink.Duration
		currentLocalTime = time.Now().UTC()
	)
	err = s.InitGlobal(&slGlobal, currentLocalTime, station, network, location, bufferDuration)
	if err != nil {
		s.OnError(options, err)
		return
	}
	defer slGlobal.SeedLinkBuffer.Database.Close()

	// Accept incoming connections
	s.OnStart(options, "service has started")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			// Handle seedlink from client
			var slClient seedlink.SeedLinkClient
			s.InitClient(&slClient)
			go s.handleCommand(options, &slGlobal, &slClient, conn)
		}
	}()

	// Subscribe to publisher to append buffer
	expressionForSubscribe := true
	go publisher.Subscribe(&options.Status.Geophone, &expressionForSubscribe, func(gp *publisher.Geophone) error {
		return s.handleBuffer(gp, &slGlobal.SeedLinkBuffer)
	})

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signals
	<-sigCh
	logger.Print(MODULE, "closing buffer area", color.FgBlue, true)
}
