package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/hook"
	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/internal/server"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/internal/upgrade"
	"github.com/anyshake/observer/pkg/jobtracker"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/seisevent"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"

	cleaner_close_database "github.com/anyshake/observer/internal/hook/cleaner/close_database"
	cleaner_close_explorer "github.com/anyshake/observer/internal/hook/cleaner/close_explorer"
	startup_migrate_database "github.com/anyshake/observer/internal/hook/startup/migrate_database"
	startup_setup_admin "github.com/anyshake/observer/internal/hook/startup/setup_admin"
	startup_setup_station "github.com/anyshake/observer/internal/hook/startup/setup_station"

	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"

	service_archiver "github.com/anyshake/observer/internal/service/archiver"
	service_discovery "github.com/anyshake/observer/internal/service/discovery"
	service_forwarder "github.com/anyshake/observer/internal/service/forwarder"
	service_frp_client "github.com/anyshake/observer/internal/service/frp_client"
	service_helicorder "github.com/anyshake/observer/internal/service/helicorder"
	service_metrics "github.com/anyshake/observer/internal/service/metrics"
	service_miniseed "github.com/anyshake/observer/internal/service/miniseed"
	service_ntp_server "github.com/anyshake/observer/internal/service/ntp_server"
	service_quakesense "github.com/anyshake/observer/internal/service/quakesense"
	service_seedlink "github.com/anyshake/observer/internal/service/seedlink"
	service_tailscale "github.com/anyshake/observer/internal/service/tailscale"
	service_updater "github.com/anyshake/observer/internal/service/updater"
	service_watchcat "github.com/anyshake/observer/internal/service/watchcat"
	service_winston "github.com/anyshake/observer/internal/service/winston"
)

func getExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	return exePath, nil
}

func appStart(ver *semver.Version, build *unibuild.UniBuild, args arguments) {
	conf := &config.BaseConfig{}
	if err := conf.Parse(args.configPath, "json"); err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	if err := conf.Migrate(logger.GetLogger(main)); err != nil {
		logger.GetLogger(main).Fatalln(err)
	}
	logger.GetLogger(main).Infoln("global configuration has been loaded")

	var pprofServer *http.Server
	if conf.Server.Debug {
		pprofServer = &http.Server{Addr: "localhost:6060", Handler: http.DefaultServeMux}
		go func() {
			logger.GetLogger(main).Infoln("pprof server running on http://localhost:6060/debug/pprof/")
			if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.GetLogger(main).Errorf("pprof server stopped unexpectedly: %v", err)
			}
		}()
	}

	logBuffer, logPath, err := setupLogger(
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

	exePath, err := getExecutablePath()
	if err != nil {
		logger.GetLogger(main).Fatalf("failed to get current executable path: %v", err)
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
	logger.GetLogger(main).Infoln("database connection has been established")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(signalChan)

	if err = daoObj.AutoMigrate(
		&model.SchemaVersion{},
		&model.SeisRecord{},
		&model.SysUser{},
		&model.FDSNAccessKey{},
		&model.UserSettings{},
	); err != nil {
		logger.GetLogger(main).Errorln(err)
		_ = daoObj.Close()
		os.Exit(1)
	}
	logger.GetLogger(main).Infoln("database schema has been configured")
	actionHandler := action.NewHandler(daoObj)
	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	var hardwareDevice hardware.IHardware

	runCleanerTasks := func(ctx context.Context) error {
		var errs []error
		cleanerTasks := []hook.IHook{
			&cleaner_close_explorer.CloseExplorerCleanerImpl{HardwareDev: hardwareDevice},
			&cleaner_close_database.CloseDatabaseCleanerImpl{DAO: daoObj},
		}
		for _, t := range cleanerTasks {
			taskName := t.GetName()
			logger.GetLogger(taskName).Infof("running cleaner task: %s", taskName)
			if err := t.Execute(ctx); err != nil {
				logger.GetLogger(taskName).Errorf("failed to run cleaner task %s: %v", taskName, err)
				errs = append(errs, fmt.Errorf("%s: %w", taskName, err))
			}
		}
		return errors.Join(errs...)
	}

	stationConfigConstraints := config.NewStationConstraints()
	startupTasks := []hook.IHook{
		&startup_migrate_database.MigrateDatabaseStartupImpl{
			ActionHandler: actionHandler,
		},
		&startup_setup_station.SetupStationStartupImpl{
			ActionHandler:            actionHandler,
			StationConfigConstraints: stationConfigConstraints,
		},
		&startup_setup_admin.SetupAdminStartupImpl{ActionHandler: actionHandler},
	}
	for _, t := range startupTasks {
		taskName := t.GetName()
		logger.GetLogger(taskName).Infof("running startup task: %s", taskName)
		if err = t.Execute(appCtx); err != nil {
			logger.GetLogger(taskName).Errorf("failed to run startup task %s: %v", taskName, err)
			runCleanerTasks(appCtx)
			os.Exit(1)
		}
	}

	timeSrc := timesource.New(nil)
	hardwareDevice, err = hardware.New(
		appCtx,
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
			Pool:        conf.NtpClient.Pool,
			Retry:       conf.NtpClient.Retry,
			ReadTimeout: conf.NtpClient.Timeout,
		},
	)
	if err != nil {
		logger.GetLogger(main).Errorf("failed to create explorer instance: %v", err)
		runCleanerTasks(appCtx)
		os.Exit(1)
	}
	hardwareCtx, sendHardwareAbortSignal, interrupted, err := openHardware(signalChan, hardwareDevice.Open)
	if err != nil {
		hardwareDevice = nil
	} else if interrupted {
		sendHardwareAbortSignal()
	}
	if interrupted {
		logger.GetLogger(main).Warnln("interrupt signal received during startup, shutting down...")
		runCleanerTasks(appCtx)
		os.Exit(0)
	}
	if err != nil {
		logger.GetLogger(main).Errorf("failed to open explorer instance: %v", err)
		runCleanerTasks(appCtx)
		os.Exit(1)
	}
	logger.GetLogger(main).Infof("hardware device has been connected, current time in UTC: %s", timeSrc.Now().Format(time.RFC3339))

	var (
		restartOnce sync.Once
		restartChan = make(chan struct{}, 1)
	)
	notificationHub := notification.NewHub(32)

	serviceMap := map[string]service.IService{
		service_archiver.ID:   service_archiver.New(appCtx, hardwareDevice, actionHandler, timeSrc),
		service_forwarder.ID:  service_forwarder.New(appCtx, hardwareDevice, actionHandler, timeSrc),
		service_frp_client.ID: service_frp_client.New(appCtx, conf.Server.Listen, actionHandler, timeSrc),
		service_helicorder.ID: service_helicorder.New(appCtx, hardwareDevice, actionHandler, timeSrc),
		service_metrics.ID:    service_metrics.New(appCtx, hardwareDevice, actionHandler, timeSrc, ver, build),
		service_miniseed.ID:   service_miniseed.New(appCtx, hardwareDevice, actionHandler, timeSrc),
		service_ntp_server.ID: service_ntp_server.New(appCtx, actionHandler, timeSrc),
		service_quakesense.ID: service_quakesense.New(appCtx, hardwareDevice, actionHandler, timeSrc, notificationHub),
		service_seedlink.ID:   service_seedlink.New(appCtx, hardwareDevice, actionHandler, timeSrc),
		service_tailscale.ID: service_tailscale.New(
			appCtx,
			conf.Server.Listen,
			actionHandler,
			timeSrc,
		),
		service_watchcat.ID: service_watchcat.New(appCtx, hardwareDevice, timeSrc),
		service_winston.ID:  service_winston.New(appCtx, hardwareDevice, actionHandler, timeSrc),
	}
	var upgradeHelper *upgrade.Helper
	if !ver.IsPreRelease() && build.GetChannel() == officialBuildChannel {
		upgradeHelper = upgrade.NewHelper(exePath, ver, build)
		serviceMap[service_updater.ID] = service_updater.New(appCtx, actionHandler, timeSrc, upgradeHelper, exePath, restartChan)
	}
	serviceMap[service_discovery.ID] = service_discovery.New(appCtx, conf.Server.Listen, serviceMap, actionHandler, timeSrc)

	seisEventSource, closeSeisEventSource, err := seisevent.New(timeSrc, 30*time.Second)
	if err != nil {
		logger.GetLogger(main).Errorf("failed to create seis event source: %v", err)
	}

	graphqlResolver := &graph_resolver.Resolver{
		AppContext:               appCtx,
		RestartChan:              restartChan,
		ExecutablePath:           exePath,
		CurrentVersion:           ver,
		CurrentBuild:             build,
		UpgradeHelper:            upgradeHelper,
		ServiceMap:               serviceMap,
		HardwareDev:              hardwareDevice,
		TimeSource:               timeSrc,
		ActionHandler:            actionHandler,
		LogBuffer:                logBuffer,
		SeisEventSource:          seisEventSource,
		StationConfigConstraints: stationConfigConstraints,
		DataPurgeJob:             jobtracker.New(appCtx, "data_purge_job"),
	}
	httpServer := server.New(appCtx, conf.Server.Debug, conf.Server.CORS, graphqlResolver, notificationHub, logger.GetLogger("http_server"))
	if err = httpServer.Setup(conf.Server.Listen); err != nil {
		logger.GetLogger(main).Errorln(err)
		sendHardwareAbortSignal()
		runCleanerTasks(appCtx)
		os.Exit(1)
	}
	if err = httpServer.Listen(); err != nil {
		logger.GetLogger(main).Errorln(err)
		sendHardwareAbortSignal()
		runCleanerTasks(appCtx)
		os.Exit(1)
	}

	for serviceName, serviceObj := range serviceMap {
		if err = serviceObj.Init(appCtx); err != nil {
			logger.GetLogger(serviceName).Errorf("failed to initialize service %s: %v", serviceName, err)
			continue
		}
		if enabled := serviceObj.IsEnabled(appCtx); !enabled {
			logger.GetLogger(serviceName).Infof("service %s is disabled by configuration", serviceName)
			continue
		}
		if err = serviceObj.Start(appCtx); err != nil {
			logger.GetLogger(serviceName).Errorf("failed to start service %s: %v", serviceName, err)
			continue
		}
		logger.GetLogger(serviceName).Infof("service %s has been started", serviceName)
	}

	httpErrChan := make(chan error, 1)
	go func() {
		if err := httpServer.Start(); err != nil {
			httpErrChan <- err
		}
	}()

	var exitWithReason ExitReason
	select {
	case <-signalChan:
		exitWithReason = ExitInterrupt
		logger.GetLogger(main).Warnln("interrupt signal received (e.g. Ctrl+C), shutting down...")
	case <-restartChan:
		exitWithReason = ExitRestart
		logger.GetLogger(main).Infoln("application restart request received, shutting down...")
	case <-hardwareCtx.Done():
		exitWithReason = ExitError
		logger.GetLogger(main).Warnln("fatal error detected (probably hardware connection lost), shutting down...")
	case err := <-httpErrChan:
		exitWithReason = ExitError
		logger.GetLogger(main).Errorf("http server stopped unexpectedly, shutting down: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	graphqlResolver.BeginShutdown()
	cancelApp()

	shutdown := &applicationShutdown{
		resolver: graphqlResolver,
		stopHTTP: httpServer.Stop,
		services: serviceMap,
		cleanup: func(ctx context.Context) error {
			sendHardwareAbortSignal()
			var errs []error
			if closeSeisEventSource != nil {
				errs = append(errs, closeSeisEventSource())
			}
			if upgradeHelper != nil {
				errs = append(errs, upgradeHelper.Close())
			}
			return errors.Join(append(errs, runCleanerTasks(ctx))...)
		},
	}
	if pprofServer != nil {
		shutdown.stopPprof = pprofServer.Shutdown
	}
	done := make(chan error, 1)
	go func() { done <- shutdown.Run(shutdownCtx) }()

	handleExit := func(reason string, warn bool, callback func()) {
		if warn {
			logger.GetLogger(main).Warnln(reason)
		} else {
			logger.GetLogger(main).Infoln(reason)
		}
		if callback != nil {
			callback()
		}
		switch {
		case warn:
			os.Exit(1)
		case exitWithReason == ExitInterrupt:
			os.Exit(0)
		case exitWithReason == ExitRestart:
			os.Exit(0)
		case exitWithReason == ExitError:
			os.Exit(1)
		}
	}

	select {
	case err := <-done:
		if shutdownCtx.Err() != nil {
			handleExit("shutdown timed out, forcing exit", true, nil)
		}
		if err != nil {
			handleExit(fmt.Sprintf("shutdown completed with errors: %v", err), true, nil)
		}
		handleExit("program exited successfully, goodbye", false, func() {
			if exitWithReason == ExitRestart {
				restartOnce.Do(func() { executeBinary(exePath) })
			}
		})
	case <-shutdownCtx.Done():
		for _, name := range shutdown.Pending() {
			logger.GetLogger(main).Warnf("shutdown incomplete: %s", name)
		}
		handleExit("shutdown timed out, forcing exit", true, nil)
	case <-signalChan:
		for _, name := range shutdown.Pending() {
			logger.GetLogger(main).Warnf("shutdown incomplete: %s", name)
		}
		handleExit("interrupt signal received again, forcing exit", true, nil)
	}
}
