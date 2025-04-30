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
	"github.com/anyshake/observer/internal/hook"
	cleaner_close_database "github.com/anyshake/observer/internal/hook/cleaner/close_database"
	cleaner_close_explorer "github.com/anyshake/observer/internal/hook/cleaner/close_explorer"
	startup_setup_admin "github.com/anyshake/observer/internal/hook/startup/setup_admin"
	startup_setup_station "github.com/anyshake/observer/internal/hook/startup/setup_station"
	"github.com/anyshake/observer/internal/server"
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	service_archiver "github.com/anyshake/observer/internal/service/archiver"

	service_forwarder "github.com/anyshake/observer/internal/service/forwarder"
	service_helicorder "github.com/anyshake/observer/internal/service/helicorder"

	service_metrics "github.com/anyshake/observer/internal/service/metrics"
	service_miniseed "github.com/anyshake/observer/internal/service/miniseed"

	service_quakesense "github.com/anyshake/observer/internal/service/quakesense"
	service_seedlink "github.com/anyshake/observer/internal/service/seedlink"
	service_timesync "github.com/anyshake/observer/internal/service/timesync"

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

	logger.GetLogger(main).Infof("querying NTP server at %s", conf.NtpClient.Endpoint)
	ntpTimeSource, err := timesource.NewNtpClient(
		conf.NtpClient.Endpoint,
		conf.NtpClient.Retry,
		time.Duration(conf.NtpClient.Timeout)*time.Second,
	)
	if err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Infof("application time has been synced with NTP server")
	logger.GetLogger(main).Infof("current network time in UTC timezone: %s", ntpTimeSource.Get().Format(time.RFC3339))

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

	hardwareDevice, err = hardware.New(
		conf.Hardware.Endpoint,
		conf.Hardware.Protocol,
		conf.Hardware.Model,
		conf.Hardware.Timeout,
		conf.Location.Latitude,
		conf.Location.Longitude,
		conf.Location.Elevation,
		actionHandler,
		ntpTimeSource,
		logger.GetLogger("explorer_driver"),
	)
	if err != nil {
		logger.GetLogger(main).Fatalf("failed to create explorer instance: %v", err)
	}
	hardwareCtx, sendHardwareAbortSignal, err := hardwareDevice.Open(context.Background())
	if err != nil {
		logger.GetLogger(main).Fatalf("failed to open explorer instance: %v", err)
	}
	logger.GetLogger(main).Info("harware device has been connected")

	serviceMap := map[string]service.IService{
		service_archiver.ID:   service_archiver.New(hardwareDevice, actionHandler, ntpTimeSource),
		service_forwarder.ID:  service_forwarder.New(hardwareDevice, actionHandler, ntpTimeSource),
		service_helicorder.ID: service_helicorder.New(hardwareDevice, actionHandler, ntpTimeSource),
		service_metrics.ID:    service_metrics.New(hardwareDevice, actionHandler, ntpTimeSource, binaryVersion, commitHash, buildPlatform),
		service_miniseed.ID:   service_miniseed.New(hardwareDevice, actionHandler, ntpTimeSource),
		service_quakesense.ID: service_quakesense.New(hardwareDevice, actionHandler, ntpTimeSource),
		service_seedlink.ID:   service_seedlink.New(hardwareDevice, actionHandler, ntpTimeSource),
		service_timesync.ID:   service_timesync.New(ntpTimeSource),
		service_watchcat.ID:   service_watchcat.New(hardwareDevice, ntpTimeSource),
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
			TimeSource:               ntpTimeSource,
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

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-osSignal:
		logger.GetLogger(main).Warnln("interrupt signal received (e.g. Ctrl+C), shutting down...")
	case <-hardwareCtx.Done():
		logger.GetLogger(main).Warnln("fatal error detected (probably hardware connection lost), shutting down...")
	}

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

	sendHardwareAbortSignal()
	runCleanerTasks()
}
