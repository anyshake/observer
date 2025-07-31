package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/hook"
	cleaner_close_database "github.com/anyshake/observer/internal/hook/cleaner/close_database"
	cleaner_close_explorer "github.com/anyshake/observer/internal/hook/cleaner/close_explorer"
	startup_setup_admin "github.com/anyshake/observer/internal/hook/startup/setup_admin"
	startup_setup_station "github.com/anyshake/observer/internal/hook/startup/setup_station"
	"github.com/anyshake/observer/internal/server"
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	service_archiver "github.com/anyshake/observer/internal/service/archiver"

	service_forwarder "github.com/anyshake/observer/internal/service/forwarder"
	service_frp_client "github.com/anyshake/observer/internal/service/frp_client"
	service_helicorder "github.com/anyshake/observer/internal/service/helicorder"

	service_metrics "github.com/anyshake/observer/internal/service/metrics"
	service_miniseed "github.com/anyshake/observer/internal/service/miniseed"

	service_quakesense "github.com/anyshake/observer/internal/service/quakesense"
	service_seedlink "github.com/anyshake/observer/internal/service/seedlink"

	_ "net/http/pprof"

	service_watchcat "github.com/anyshake/observer/internal/service/watchcat"

	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/seisevent"
	"github.com/anyshake/observer/pkg/timesource"
)

func appStart(args arguments) {
	conf := &config.BaseConfig{}
	if err := conf.Parse(args.configPath, "json"); err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Info("global configuration has been loaded")

	if conf.Server.Debug {
		go func() {
			logger.GetLogger(main).Infoln("pprof server running on http://localhost:6060/debug/pprof/")
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	logPath, err := setupLogger(
		conf.Logger.Level,
		conf.Logger.Path,
		conf.Logger.Size,
		conf.Logger.Rotation,
		conf.Logger.LifeCycle,
	)
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	if len(logPath) != 0 {
		logger.GetLogger(main).Infof("logs will be saved to: %s", logPath)
	}

	daoObj, err := dao.New(
		conf.Database.Endpoint,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Prefix,
		time.Duration(conf.Database.Timeout)*time.Second,
	)
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	if err = daoObj.Open(conf.Database.Database); err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Info("database connection has been established")

	if err = daoObj.Migrate(
		&model.SeisRecord{},
		&model.SysUser{},
		&model.UserSettings{},
	); err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Info("database schema has been configured")
	actionHandler := action.New(daoObj)

	var hardwareDevice hardware.IHardware

	runCleanerTasks := func() {
		cleanerTasks := []hook.IHook{
			&cleaner_close_explorer.CloseExplorerCleanerImpl{HardwareDev: hardwareDevice},
			&cleaner_close_database.CloseDatabaseCleanerImpl{DAO: daoObj},
		}
		for _, t := range cleanerTasks {
			taskName := t.GetName()
			logger.GetLogger(taskName).Infof("running cleaner task: %s", taskName)
			if err = t.Execute(); err != nil {
				logger.GetLogger(taskName).Errorf("failed to run cleaner task %s: %v", taskName, err)
			}
		}
	}

	stationConfigConstraints := config.NewStationConstraints()
	startupTasks := []hook.IHook{
		&startup_setup_station.SetupStationStartupImpl{
			ActionHandler:            actionHandler,
			StationConfigConstraints: stationConfigConstraints,
		},
		&startup_setup_admin.SetupAdminStartupImpl{ActionHandler: actionHandler},
	}
	for _, t := range startupTasks {
		taskName := t.GetName()
		logger.GetLogger(taskName).Infof("running startup task: %s", taskName)
		if err = t.Execute(); err != nil {
			logger.GetLogger(taskName).Errorf("failed to run startup task %s: %v", taskName, err)
			runCleanerTasks()
			os.Exit(1)
		}
	}

	timeSrc := timesource.New()
	hardwareDevice, err = hardware.New(
		logger.GetLogger("explorer_driver"),
		timeSrc,
		actionHandler,
		explorer.ExplorerOptions{
			Endpoint:    conf.Hardware.Endpoint,
			Protocol:    conf.Hardware.Protocol,
			Model:       conf.Hardware.Model,
			Latitude:    conf.Location.Latitude,
			Longitude:   conf.Location.Longitude,
			Elevation:   conf.Location.Elevation,
			ReadTimeout: conf.Hardware.Timeout,
		},
		explorer.NtpOptions{
			Endpoint:    conf.NtpClient.Endpoint,
			Retry:       conf.NtpClient.Retry,
			ReadTimeout: conf.NtpClient.Timeout,
		},
	)
	if err != nil {
		logger.GetLogger(main).Fatalf("failed to create explorer instance: %v", err)
	}
	hardwareCtx, sendHardwareAbortSignal, err := hardwareDevice.Open(context.Background())
	if err != nil {
		logger.GetLogger(main).Fatalf("failed to open explorer instance: %v", err)
	}
	logger.GetLogger(main).Infof("hardware device has been connected, current time in UTC: %s", timeSrc.Now().Format(time.RFC3339))

	serviceMap := map[string]service.IService{
		service_archiver.ID:   service_archiver.New(hardwareDevice, actionHandler, timeSrc),
		service_forwarder.ID:  service_forwarder.New(hardwareDevice, actionHandler, timeSrc),
		service_frp_client.ID: service_frp_client.New(conf.Server.Listen, actionHandler, timeSrc),
		service_helicorder.ID: service_helicorder.New(hardwareDevice, actionHandler, timeSrc),
		service_metrics.ID:    service_metrics.New(hardwareDevice, actionHandler, timeSrc, binaryVersion, commitHash, buildPlatform),
		service_miniseed.ID:   service_miniseed.New(hardwareDevice, actionHandler, timeSrc),
		service_quakesense.ID: service_quakesense.New(hardwareDevice, actionHandler, timeSrc),
		service_seedlink.ID:   service_seedlink.New(hardwareDevice, actionHandler, timeSrc),
		service_watchcat.ID:   service_watchcat.New(hardwareDevice, timeSrc),
	}

	for serviceName, serviceObj := range serviceMap {
		if err = serviceObj.Init(); err != nil {
			logger.GetLogger(serviceName).Errorf("failed to initialize service %s: %v", serviceName, err)
			continue
		}
		if enabled := serviceObj.IsEnabled(); !enabled {
			logger.GetLogger(serviceName).Infof("service %s is disabled by configuration", serviceName)
			continue
		}
		if err = serviceObj.Start(); err != nil {
			logger.GetLogger(serviceName).Errorf("failed to start service %s: %v", serviceName, err)
			continue
		}
		logger.GetLogger(serviceName).Infof("service %s has been started", serviceName)
	}

	httpServer := server.New(
		conf.Server.Debug,
		conf.Server.CORS,
		&graph_resolver.Resolver{
			ServiceMap:               serviceMap,
			HardwareDev:              hardwareDevice,
			TimeSource:               timeSrc,
			ActionHandler:            actionHandler,
			SeisEventSource:          seisevent.New(30 * time.Second),
			StationConfigConstraints: stationConfigConstraints,
		},
		logger.GetLogger("http_server"),
	)

	if err = httpServer.Setup(conf.Server.Listen); err != nil {
		logger.GetLogger(main).Errorln(err)
		runCleanerTasks()
		os.Exit(1)
	}
	go func() {
		if err := httpServer.Start(); err != nil {
			logger.GetLogger(main).Errorln(err)
			runCleanerTasks()
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(signalChan)

	exitWithError := false
	select {
	case <-signalChan:
		logger.GetLogger(main).Warnln("interrupt signal received (e.g. Ctrl+C), shutting down...")
	case <-hardwareCtx.Done():
		exitWithError = true
		logger.GetLogger(main).Warnln("fatal error detected (probably hardware connection lost), shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)

		if err = httpServer.Stop(); err != nil {
			logger.GetLogger(main).Errorf("failed to stop http server: %v", err)
		}

		for serviceName, serviceObj := range serviceMap {
			if !serviceObj.GetStatus().GetIsRunning() {
				continue
			}
			if err = serviceObj.Stop(); err != nil {
				logger.GetLogger(serviceName).Errorf("failed to stop service %s: %v", serviceName, err)
			} else {
				logger.GetLogger(serviceName).Infof("service %s has been stopped", serviceName)
			}
		}
	}()

	handleExit := func(reason string, warn bool) {
		sendHardwareAbortSignal()
		runCleanerTasks()
		if warn {
			logger.GetLogger(main).Warn(reason)
		} else {
			logger.GetLogger(main).Info(reason)
		}
		if warn || exitWithError {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	select {
	case <-done:
		handleExit("program exited successfully, goodbye", false)
	case <-shutdownCtx.Done():
		handleExit("shutdown timed out, forcing exit", true)
	}
}
