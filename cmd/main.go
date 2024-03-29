package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fatih/color"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/feature/archiver"
	"github.com/anyshake/observer/feature/geophone"
	"github.com/anyshake/observer/feature/miniseed"
	"github.com/anyshake/observer/feature/ntpclient"
	"github.com/anyshake/observer/feature/seedlink"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/server"
	"github.com/anyshake/observer/utils/logger"
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
		&seedlink.SeedLink{},
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
			FeatureOptions: featureOptions,
			CORS:           conf.Server.CORS,
			RateFactor:     conf.Server.Rate,
		})

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Wait for all features to stop
	logger.Print("main", "main daemon is shutting down", color.FgMagenta, true)
	featureWaitGroup.Wait()
}
