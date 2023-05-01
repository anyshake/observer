package main

import (
	"log"

	"com.geophone.observer/app"
	"com.geophone.observer/common/postgres"
	"com.geophone.observer/config"
	"com.geophone.observer/features/archiver"
	"com.geophone.observer/features/collector"
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/features/ntpclient"
	"com.geophone.observer/handlers"
	"com.geophone.observer/server"
)

const apiVersion = "v1"

func main() {
	var (
		args config.Args
		conf config.Config
	)

	err := ProgramInit(&args, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		status  collector.Status
		message = collector.Message{
			Station:   conf.Station.Name,
			UUID:      conf.Station.UUID,
			Latitude:  conf.Station.Latitude,
			Longitude: conf.Station.Longitude,
			Altitude:  conf.Station.Altitude,
		}
	)

	conn, grpc, err := collector.OpenGRPC(
		conf.Collector.Host,
		conf.Collector.Port,
		conf.Collector.TLS,
		conf.Collector.Enable,
	)
	if err != nil {
		log.Fatalln(err)
	}
	if conf.Collector.Enable {
		defer collector.CloseGRPC(conn)
	}

	pdb, err := postgres.OpenPostgres(
		conf.Archiver.Host,
		conf.Archiver.Port,
		conf.Archiver.Username,
		conf.Archiver.Password,
		conf.Archiver.Database,
		conf.Collector.Enable,
	)
	if err != nil {
		log.Fatalln(err)
	}
	if conf.Archiver.Enable {
		defer postgres.ClosePostgres(pdb)
	}

	go ntpclient.ReaderDaemon(
		conf.NTPClient.Host,
		conf.NTPClient.Interval,
		ntpclient.NTPOptions{
			Port:    conf.NTPClient.Port,
			Timeout: conf.NTPClient.Timeout,
			OnErrorCallback: func(err error) {
				handlers.HandleErrors(&handlers.HandlerOptions{
					Error:  err,
					Status: &status,
				})
			},
			OnDataCallback: func(ntp *ntpclient.NTP) {
				log.Println("Read NTP server time")
				handlers.HandleMessages(&handlers.HandlerOptions{
					Status:  &status,
					Message: &message,
				}, ntp)
			},
		},
	)

	go geophone.ReaderDaemon(
		conf.Geophone.Device,
		conf.Geophone.Baud,
		geophone.GeophoneOptions{
			Geophone:     &geophone.Geophone{},
			Acceleration: &geophone.Acceleration{},
			Latitude:     conf.Station.Latitude,
			Longitude:    conf.Station.Longitude,
			Altitude:     conf.Station.Altitude,
			Sensitivity: struct {
				Vertical   float64
				EastWest   float64
				NorthSouth float64
			}{
				Vertical:   conf.Geophone.Sensitivity.Vertical,
				EastWest:   conf.Geophone.Sensitivity.EastWest,
				NorthSouth: conf.Geophone.Sensitivity.NorthSouth,
			},
			OnErrorCallback: func(err error) {
				handlers.HandleErrors(&handlers.HandlerOptions{
					Error:  err,
					Status: &status,
				})
			},
			OnDataCallback: func(acceleration *geophone.Acceleration) {
				log.Println("1 message received")
				handlers.HandleMessages(&handlers.HandlerOptions{
					Status:  &status,
					Message: &message,
					OnReadyCallback: func(message *collector.Message) {
						archiver.WriteMessage(
							pdb, &archiver.ArchiverOptions{
								Message: message,
								Enable:  conf.Archiver.Enable,
								OnCompleteCallback: func() {
									log.Println("10 message archived")
								},
								OnErrorCallback: func(err error) {
									log.Println(err)
								},
							},
						)
						go collector.PushMessage(
							conn, grpc, &collector.CollectorOptions{
								Message: message,
								Status:  &status,
								Enable:  conf.Collector.Enable,
								OnCompleteCallback: func(r any) {
									log.Println(r)
								},
								OnErrorCallback: func(err error) {
									log.Println(err)
								},
							})
					},
				}, acceleration)
			},
		},
	)

	server.ServerDaemon(
		conf.Server.Host,
		conf.Server.Port,
		&app.ServerOptions{
			WebPrefix:    "/",
			ApiPrefix:    "/api",
			Version:      apiVersion,
			Message:      &message,
			Status:       &status,
			ConnGRPC:     &grpc,
			ConnPostgres: pdb,
			Cors:         true,
			Gzip:         9,
		})
}
