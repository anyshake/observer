package main

import (
	"fmt"
	"log"
	"os"

	"com.geophone.observer/common/geophone"
	"com.geophone.observer/config"
	"github.com/common-nighthawk/go-figure"
)

var (
	args config.Args
	conf config.Config
)

func main() {
	figure.NewFigure(
		"G-Observer",
		"standard",
		true,
	).Print()
	fmt.Println()
	args.ReadFlags()

	if args.Version {
		PrintVersion()
		os.Exit(0)
	}

	err := conf.ReadConfig(args.Path)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		phone        geophone.Geophone
		acceleration geophone.Acceleration
	)

	geophone.ReaderDaemon(
		conf.Geophone.Device,
		conf.Geophone.Baud,
		geophone.ReaderOptions{
			Geophone:     &phone,
			Acceleration: &acceleration,
			Sensitivity:  conf.Geophone.Sensitivity,
			OnErrorCallback: func(err error) {
				log.Println(err)
			},
			OnDataCallback: func(acceleration *geophone.Acceleration) {
				log.Print(".")
			},
		},
	)
}
