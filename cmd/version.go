package main

import (
	"fmt"
	"runtime"
	"time"
)

var (
	release     = "unknown"
	version     = "custombuild"
	description = "Constructing Real-time Seismic Network Ambitiously"
)

func printVersion() {
	copyright := "© AnyShake " + fmt.Sprintf("%d", time.Now().Year()) + ". All Rights Reversed."
	version = concat(
		"AnyShake Observer ", version, " (", description, ")\nRelease: ", version, "-", release, " ",
		runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n", copyright, "\n",
	)

	fmt.Println(version)
}
