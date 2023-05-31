package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	version     = "Custom build"
	release     = "unknown"
	description = "Constructing Real-time Seismic Network Ambitiously"
	copyright   = "© Telecom-308 " + fmt.Sprintf("%d", time.Now().Year()) + ". All Rights Reversed."
)

func ConcatString(v ...any) string {
	builder := strings.Builder{}
	for _, value := range v {
		builder.WriteString(fmt.Sprintf("%+v", value))
	}

	return builder.String()
}

func VersionStatement() []string {
	return []string{
		ConcatString(
			"G-Observer ", version, " (", description, ")\nRelease: ", release, " ",
			runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n", copyright,
		),
	}
}

func PrintVersion() {
	version := VersionStatement()
	for _, s := range version {
		fmt.Println(s)
		os.Exit(0)
	}
}
