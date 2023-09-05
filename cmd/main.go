package main

import (
	"fmt"
	"os"

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

	logger.Print("main", "Observer daemon initialized", color.FgMagenta)
	return nil
}

func init() {
	w := color.New(color.FgHiCyan).SprintFunc()
	t := figure.NewFigure("Observer", "standard", true).String()
	fmt.Println(w(t))
}

func main() {
	// Read configuration
	var conf config.Conf
	err := parseCommandLine(&conf)
	if err != nil {
		logger.Fatal("main", err, color.FgRed)
	}

	// Initialize global status
	var status publisher.Status
	publisher.Initialize(&conf, &status)

	// Register features
	featureOptions := &feature.FeatureOptions{
		Config: &conf,
		Status: &status,
	}
	features := []feature.Feature{
		&ntpclient.NTPClient{},
		&geophone.Geophone{},
		&archiver.Archiver{},
		&miniseed.MiniSEED{},
	}
	for _, s := range features {
		go s.Start(featureOptions)
	}

	// Start HTTP server
	server.ServerDaemon(
		conf.Server.Host,
		conf.Server.Port,
		&app.ServerOptions{
			Gzip:           9,
			WebPrefix:      WEB_PREFIX,
			CORS:           API_CORS,
			APIPrefix:      API_PREFIX,
			Version:        API_VERSION,
			FeatureOptions: featureOptions,
		})
}
