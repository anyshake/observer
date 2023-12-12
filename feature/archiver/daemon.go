package archiver

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anyshake/observer/driver/dao"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/logger"
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

	// Connect to database
	a.OnStart(options, "service has started")
	pdb, err := dao.Open(
		options.Config.Archiver.Host,
		options.Config.Archiver.Port,
		options.Config.Archiver.Engine,
		options.Config.Archiver.Username,
		options.Config.Archiver.Password,
		options.Config.Archiver.Database,
	)
	if err != nil {
		a.OnError(options, err)
		os.Exit(1)
	}

	// Migrate database
	err = dao.Migrate(pdb)
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
		os.Exit(1)
	}()

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signals
	<-sigCh
	logger.Print(MODULE, "closing database connection", color.FgBlue, true)
	dao.Close(pdb)
}
