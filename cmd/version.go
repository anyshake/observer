package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/bclswl0827/observer/utils/text"
	"github.com/fatih/color"
)

var (
	release     = "unknown"
	version     = "Custom build"
	description = "Constructing Real-time Seismic Network Ambitiously"
)

func printVersion() {
	var (
		copyright = "© Telecom-308 " + fmt.Sprintf("%d", time.Now().Year()) + ". All Rights Reversed."
		version   = text.Concat(
			"G-Observer ", version, " (", description, ")\nRelease: ", release, " ",
			runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n", copyright,
		)
	)

	w := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println(w(version))
	os.Exit(0)
}
