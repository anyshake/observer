package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/dig"

	"github.com/anyshake/observer/cleaners"
	cleaner_database "github.com/anyshake/observer/cleaners/database"
	cleaner_explorer "github.com/anyshake/observer/cleaners/explorer"
	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/drivers/dao"
	"github.com/anyshake/observer/server"
	"github.com/anyshake/observer/services"
	service_archiver "github.com/anyshake/observer/services/archiver"

	service_miniseed "github.com/anyshake/observer/services/miniseed"
	service_watchdog "github.com/anyshake/observer/services/watchdog"
	"github.com/anyshake/observer/startups"
	startup_explorer "github.com/anyshake/observer/startups/explorer"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/timesource"
	"github.com/beevik/ntp"
	"github.com/common-nighthawk/go-figure"
)

func parseCommandLine() (args arguments) {
	flag.StringVar(&args.Path, "config", "./config.json", "Path to config file")
	flag.BoolVar(&args.Version, "version", false, "Print version information")
	flag.Parse()

	if args.Version {
		printVersion()
		os.Exit(0)
	}

	return args
}

func setupLogger(level, dumpPath string) {
	var err error
	switch level {
	case "info":
		err = logger.SetLevel(logger.INFO)
	case "warn":
		err = logger.SetLevel(logger.WARN)
	case "error":
		err = logger.SetLevel(logger.ERROR)
	default:
		err = logger.SetLevel(logger.INFO)
	}

	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}

	if len(dumpPath) != 0 {
		logger.SetFile(dumpPath)
	}
}

func init() {
	t := figure.NewFigure("Observer", "standard", true).String()
	fmt.Println(t)
	logger.Initialize()
}

// @BasePath /api/v1
// @title AnyShake Observer APIv1
// @description This is APIv1 documentation for AnyShake Observer, please set `server_settings.debug` to `false` in `config.json` when deploying to production environment in case of any security issues.
func main() {
	args := parseCommandLine()
	printVersion()

	var conf config.Config
	err := conf.Read(args.Path)
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	err = conf.Validate()
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}

	// Setup logger with given configuration
	setupLogger(conf.Logger.Level, conf.Logger.Dump)
	logger.GetLogger(main).Info("global configuration has been loaded")

	// Create time source with NTP server
	logger.GetLogger(main).Infof("querying NTP server at %s:%d", conf.NtpClient.Host, conf.NtpClient.Port)
	res, err := ntp.QueryWithOptions(conf.NtpClient.Host, ntp.QueryOptions{
		Port: conf.NtpClient.Port, Timeout: 10 * time.Second,
	})
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	timeSource := timesource.New(time.Now(), res.Time)
	logger.GetLogger(main).Info("time source has been created")

	// Connect to database
	databaseConn, err := dao.Open(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Engine,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Database,
	)
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Info("database connection has been established")
	err = migrate(databaseConn)
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Info("database schema has been migrated")

	// Create dependency injection container
	depsContainer := dig.New()

	// Setup context for graceful shutdown
	cancelToken, abortSignal := context.WithCancel(context.Background())

	// Setup cleaner tasks for graceful shutdown
	cleanerTasks := []cleaners.CleanerTask{
		&cleaner_explorer.ExplorerCleanerTask{},
		&cleaner_database.DatabaseCleanerTask{},
	}
	cleanerOptions := &cleaners.Options{
		Config:     &conf,
		Database:   databaseConn,
		Dependency: depsContainer,
		TimeSource: timeSource,
	}
	runCleanerTasks := func() {
		for _, t := range cleanerTasks {
			taskName := t.GetTaskName()
			logger.GetLogger(taskName).Infof("running cleaner task for %s", taskName)
			t.Execute(cleanerOptions)
		}
	}
	defer runCleanerTasks()

	// Setup startup tasks and provide dependencies
	startupTasks := []startups.StartupTask{
		&startup_explorer.ExplorerStartupTask{CancelToken: cancelToken},
	}
	startupOptions := &startups.Options{
		Config:     &conf,
		Database:   databaseConn,
		TimeSource: timeSource,
	}
	for _, t := range startupTasks {
		taskName := t.GetTaskName()
		err := t.Provide(depsContainer, startupOptions)
		if err != nil {
			logger.GetLogger(taskName).Errorln(err)
			runCleanerTasks()
			os.Exit(1)
		}
		err = t.Execute(depsContainer, startupOptions)
		if err != nil {
			logger.GetLogger(taskName).Errorln(err)
			runCleanerTasks()
			os.Exit(1)
		}
	}

	// Setup background services
	regServices := []services.Service{
		&service_watchdog.WatchdogService{},
		&service_archiver.ArchiverService{},
		&service_miniseed.MiniSeedService{},
	}
	serviceOptions := &services.Options{
		Config:      &conf,
		Database:    databaseConn,
		Dependency:  depsContainer,
		TimeSource:  timeSource,
		CancelToken: cancelToken,
	}
	var waitGroup sync.WaitGroup
	for _, s := range regServices {
		waitGroup.Add(1)
		go s.Start(serviceOptions, &waitGroup)
	}

	// Start HTTP server
	go server.Serve(
		conf.Server.Host,
		conf.Server.Port,
		&server.Options{
			CORS:            conf.Server.CORS,
			DebugMode:       conf.Server.Debug,
			GzipLevel:       GZIP_LEVEL,
			RateFactor:      conf.Server.Rate,
			WebPrefix:       WEB_PREFIX,
			ApiPrefix:       API_PREFIX,
			ServicesOptions: serviceOptions,
		})
	logger.GetLogger(main).Infof("web server is listening on %s:%d", conf.Server.Host, conf.Server.Port)

	// Receive interrupt signals
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	<-osSignal

	// Stop services gracefully
	logger.GetLogger(main).Info("services are shutting down, please wait")
	abortSignal()
	waitGroup.Wait()
}
