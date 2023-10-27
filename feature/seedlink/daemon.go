package seedlink

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
)

func (s *SeedLink) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	if !options.Config.SeedLink.Enable {
		s.OnStop(options, "service is disabled")
		return
	}

	// Create TCP server and listen
	host, port := options.Config.SeedLink.Host, options.Config.SeedLink.Port
	li, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		s.OnError(options, err)
		os.Exit(1)
	}
	defer li.Close()

	// Subscribe to geophone publisher
	go func() {
		publisher.Subscribe(
			&options.Status.Geophone,
			func(gp *publisher.Geophone) error {
				return s.handleMessage(gp, options)
			},
		)

		err := fmt.Errorf("service exited with an error")
		s.OnError(options, err)
	}()

	// Accept incoming connections
	s.OnStart(options, "service has started")
	for {
		conn, err := li.Accept()
		if err != nil {
			s.OnError(options, err)
			continue
		}

		go s.handleConnection(options, conn)
	}
}
