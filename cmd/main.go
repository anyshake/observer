package main

import (
	"fmt"
	"log"

	"com.geophone.observer/common/collector"
	"com.geophone.observer/common/geophone"
	"com.geophone.observer/common/ntpclient"
	"com.geophone.observer/config"
	"com.geophone.observer/helper/handler"
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
			LocationFallback: struct {
				Latitude  float64
				Longitude float64
				Altitude  float64
			}{
				Latitude:  conf.Station.Latitude,
				Longitude: conf.Station.Longitude,
				Altitude:  conf.Station.Altitude,
			},
			Geophone:     &geophone.Geophone{},
			Acceleration: &geophone.Acceleration{},
			Sensitivity:  conf.Geophone.Sensitivity,
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
						go collector.PushCollection(conn, grpc, &collector.CollectorOptions{
							Status:  &status,
							Message: message,
							OnCompleteCallback: func(r interface{}) {
								fmt.Println(r)
							},
							OnErrorCallback: func(err error) {
								fmt.Println(err)
							},
						})
					},
				}, acceleration)
			},
		},
	)
}
