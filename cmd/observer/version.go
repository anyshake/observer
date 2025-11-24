package main

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/unibuild"
)

func concatText(v ...any) string {
	builder := strings.Builder{}
	for _, value := range v {
		builder.WriteString(fmt.Sprintf("%+v", value))
	}

	return builder.String()
}

func printMap(m map[string]any, indentLevel int) {
	indentRepeat := func(level int) string {
		return string(make([]byte, level*2))
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := m[key]
		fmt.Printf(" %s%s:\n", indentRepeat(indentLevel), key)
		switch v := value.(type) {
		case map[string]any:
			printMap(v, indentLevel+1)
		default:
			fmt.Printf("%s  %v\n", indentRepeat(indentLevel+1), v)
		}
	}
	fmt.Println()
}

func PrintVersion(ver *semver.Version, build *unibuild.UniBuild, verbose bool) {
	binaryVersion := ver.String()
	buildCommit := build.GetCommit()
	buildChannel := build.GetChannel()
	buildTime := build.GetTime()
	toolchainId := build.GetToolchainId()

	buildTimeStr := fmt.Sprintf("-%d", buildTime.Unix())
	if buildTime.UnixMilli() == 0 {
		buildTimeStr = ""
	}

	fmt.Println(concatText(
		"AnyShake Observer ", binaryVersion, " (", startupDescription, ")\nRelease: ", binaryVersion, "-", buildCommit, buildTimeStr, " ",
		runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n",
	))
	if verbose {
		detailedMap := map[string]any{
			"Build Version": binaryVersion,
			"Build Commit":  buildCommit,
			"Build Channel": buildChannel,
		}
		if toolchainId != "" {
			detailedMap["Build Toolchain"] = toolchainId
		}
		if buildTime.UnixMilli() != 0 {
			detailedMap["Build Time"] = buildTime
		}
		printMap(detailedMap, 2)
	}
}
