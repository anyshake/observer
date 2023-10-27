package archiver

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bclswl0827/observer/driver/postgres"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/fatih/color"
)

func (a *Archiver) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	if !options.Config.Archiver.Enable {
		a.OnStop(options, "service is disabled")
		return
	} else {
		waitGroup.Add(1)
		defer waitGroup.Done()
	}

	// Connect to PostgreSQL
	a.OnStart(options, "service has started")
	pdb, err := postgres.Open(
		options.Config.Archiver.Host,
		options.Config.Archiver.Port,
		options.Config.Archiver.Username,
		options.Config.Archiver.Password,
		options.Config.Archiver.Database,
	)
	if err != nil {
		a.OnError(options, err)
		os.Exit(1)
	}

	// Initialize PostgreSQL
	err = postgres.Init(pdb)
	if err != nil {
		a.OnError(options, err)
		os.Exit(1)
	}
	options.Database = pdb

	// Archive when new message arrived
	go func() {
		publisher.Subscribe(
			&options.Status.Geophone,
			func(gp *publisher.Geophone) error {
				return a.handleMessage(gp, options, pdb)
			},
		)

		err = fmt.Errorf("service exited with an error")
		a.OnError(options, err)
	}()

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signals
	<-sigCh
	logger.Print(MODULE, "closing database connection", color.FgBlue, true)
	postgres.Close(pdb)
}
