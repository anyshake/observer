package main

import (
	"log"

	"com.geophone.observer/app"
	"com.geophone.observer/common/redis"
	"com.geophone.observer/config"
	"com.geophone.observer/features/archiver"
	"com.geophone.observer/features/collector"
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/features/ntpclient"
	"com.geophone.observer/handler"
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
			Station: conf.Station.Name,
			UUID:    conf.Station.UUID,
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

	rdb, err := redis.OpenRedis(
		conf.Archiver.Host,
		conf.Archiver.Port,
		conf.Archiver.Password,
		conf.Archiver.Database,
		conf.Archiver.Enable,
	)
	if err != nil {
		log.Fatalln(err)
	}
	if conf.Archiver.Enable {
		defer redis.CloseRedis(rdb)
	}

	go ntpclient.ReaderDaemon(
		conf.NTPClient.Host,
		conf.NTPClient.Interval,
		ntpclient.NTPOptions{
			Port:    conf.NTPClient.Port,
			Timeout: conf.NTPClient.Timeout,
			OnErrorCallback: func(err error) {
				handler.HandleErrors(&handler.HandlerOptions{
					Error:  err,
					Status: &status,
				})
			},
			OnDataCallback: func(ntp *ntpclient.NTP) {
				handler.HandleMessages(&handler.HandlerOptions{
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
			Sensitivity:  conf.Geophone.Sensitivity,
			LocationFallback: struct {
				Latitude  float64
				Longitude float64
				Altitude  float64
			}{
				Latitude:  conf.Station.Latitude,
				Longitude: conf.Station.Longitude,
				Altitude:  conf.Station.Altitude,
			},
			OnErrorCallback: func(err error) {
				handler.HandleErrors(&handler.HandlerOptions{
					Error:  err,
					Status: &status,
				})
			},
			OnDataCallback: func(acceleration *geophone.Acceleration) {
				handler.HandleMessages(&handler.HandlerOptions{
					Status:  &status,
					Message: &message,
					OnReadyCallback: func(message *collector.Message) {
						archiver.WriteMessage(
							rdb, &archiver.ArchiverOptions{
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
			WebPrefix: "/",
			ApiPrefix: "/api",
			Version:   apiVersion,
			Message:   &message,
			Status:    &status,
			ConnGRPC:  &grpc,
			ConnRedis: rdb,
			Cors:      true,
			Gzip:      9,
		})
}
