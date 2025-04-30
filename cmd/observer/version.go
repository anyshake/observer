package main

import (
	"fmt"
	"runtime"
	"strings"
)

func concatText(v ...any) string {
	builder := strings.Builder{}
	for _, value := range v {
		builder.WriteString(fmt.Sprintf("%+v", value))
	}

	return builder.String()
}

func PrintVersion() {
	fmt.Println(concatText(
		"AnyShake Observer ", binaryVersion, " (", appDescription, ")\nRelease: ", binaryVersion, "-", commitHash, " ",
		runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n",
	))
}
