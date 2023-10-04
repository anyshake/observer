package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fatih/color"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/config"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/feature/archiver"
	"github.com/bclswl0827/observer/feature/geophone"
	"github.com/bclswl0827/observer/feature/miniseed"
	"github.com/bclswl0827/observer/feature/ntpclient"
	"github.com/bclswl0827/observer/publisher"
	"github.com/bclswl0827/observer/server"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/common-nighthawk/go-figure"
)

func parseCommandLine(conf *config.Conf) error {
	var args config.Args
	args.Read()
	if args.Version {
		printVersion()
		os.Exit(0)
	}

	err := conf.Read(args.Path)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	w := color.New(color.FgHiCyan).SprintFunc()
	t := figure.NewFigure("Observer", "standard", true).String()
	fmt.Println(w(t))
}

// @BasePath /api/v1
// @title Observer RESTful API documentation
// @description This is Observer RESTful API documentation, please set `server_settings.debug` to `false` in `config.json` when deploying to production environment in case of any security issues.
func main() {
	// Read configuration
	var conf config.Conf
	err := parseCommandLine(&conf)
	if err != nil {
		logger.Fatal("main", err, color.FgRed)
	} else {
		logger.Print("main", "main daemon has initialized", color.FgMagenta, false)
	}

	// Initialize global status
	var status publisher.Status
	publisher.Initialize(&conf, &status)

	// Register features
	features := []feature.Feature{
		&ntpclient.NTPClient{},
		&geophone.Geophone{},
		&archiver.Archiver{},
		&miniseed.MiniSEED{},
	}
	featureOptions := &feature.FeatureOptions{
		Config: &conf,
		Status: &status,
	}
	featureWaitGroup := new(sync.WaitGroup)
	for _, s := range features {
		go s.Run(featureOptions, featureWaitGroup)
	}

	// Start HTTP server
	go server.StartDaemon(
		conf.Server.Host,
		conf.Server.Port,
		&app.ServerOptions{
			Gzip:           9,
			WebPrefix:      WEB_PREFIX,
			APIPrefix:      API_PREFIX,
			Version:        API_VERSION,
			FeatureOptions: featureOptions,
			CORS:           conf.Server.CORS,
		})

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Wait for all features to stop
	logger.Print("main", "main daemon is shutting down", color.FgMagenta, true)
	featureWaitGroup.Wait()
}
