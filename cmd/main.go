package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"com.geophone.observer/app"
	"com.geophone.observer/config"
	"com.geophone.observer/feature"
	"com.geophone.observer/feature/archiver"
	"com.geophone.observer/feature/geophone"
	"com.geophone.observer/feature/miniseed"
	"com.geophone.observer/feature/ntpclient"
	"com.geophone.observer/handler"
	"com.geophone.observer/handler/callbacks"
	"com.geophone.observer/server"
	"com.geophone.observer/utils/logger"
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

	logger.Print("main", "G-Observer daemon initialized", color.FgMagenta)
	return nil
}

func init() {
	w := color.New(color.FgHiCyan).SprintFunc()
	t := figure.NewFigure("G-Observer", "standard", true).String()
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
	var status handler.Status
	handler.InitHandler(&conf, &status)

	// Register features
	featureOptions := &feature.FeatureOptions{
		Config:  &conf,
		Status:  &status,
		OnStart: callbacks.OnStart,
		OnStop:  callbacks.OnStop,
		OnReady: callbacks.OnReady,
		OnError: callbacks.OnError,
	}
	features := []feature.Feature{
		&ntpclient.NTP{},
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
