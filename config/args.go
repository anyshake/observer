package config

import "flag"

func (a *Args) Read() {
	flag.StringVar(&a.Path, "config", "./config.json", "Path to config file")
	flag.BoolVar(&a.Version, "version", false, "Print version information")
	flag.Parse()
}
