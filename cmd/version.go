package main

import (
	"fmt"
	"runtime"
	"time"
)

var (
	release     = "unknown"
	version     = "Custom build"
	description = "Constructing Real-time Seismic Network Ambitiously"
)

func printVersion() {
	copyright := "© AnyShake " + fmt.Sprintf("%d", time.Now().Year()) + ". All Rights Reversed."
	version = concat(
		"Observer ", version, " (", description, ")\nRelease: ", version, "-", release, " ",
		runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n", copyright,
	)

	fmt.Println(version)
}
