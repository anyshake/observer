package main

import (
	"log"

	"com.geophone.observer/common/handler"
	"com.geophone.observer/config"
	"com.geophone.observer/features/archiver"
	"com.geophone.observer/features/collector"
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/features/ntpclient"
)

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

	conn, grpc, err := collector.OpenGrpc(
		conf.Collector.Host,
		conf.Collector.Port,
		conf.Collector.TLS,
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer collector.CloseGrpc(conn)

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

	geophone.ReaderDaemon(
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
							conf.Archiver.Path,
							&archiver.ArchiverOptions{
								Message: message,
								Status:  &status,
								Enable:  conf.Archiver.Enable,
								Name:    conf.Archiver.Name,
								OnCompleteCallback: func() {
									log.Println("Message archived")
								},
								OnErrorCallback: func(err error) {
									log.Println(err)
								},
							},
						)
						go collector.PushMessage(
							conn, grpc,
							&collector.CollectorOptions{
								Status:  &status,
								Message: message,
								OnCompleteCallback: func(r interface{}) {
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

}
