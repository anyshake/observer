package main

import (
	"flag"
	"os"
)

type arguments struct {
	configPath  string
	showVersion bool
}

func parseCommandLine() (args arguments) {
	flag.BoolVar(&args.showVersion, "version", false, "Print version information")
	flag.StringVar(&args.configPath, "config", "./config.json", "Path to config file")

	flag.Parse()

	if args.showVersion {
		PrintVersion()
		os.Exit(0)
	}

	return args
}
