package main

import (
	"flag"
)

type arguments struct {
	configPath  string
	showVersion bool
	upgrade     bool
}

func parseCommandLine() (args arguments) {
	flag.BoolVar(&args.showVersion, "version", false, "Print version information")
	flag.StringVar(&args.configPath, "config", "./config.json", "Path to config file")

	flag.Parse()

	return args
}
