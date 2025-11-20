package main

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

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
	fmt.Println(concatText(
		"AnyShake Observer ", binaryVersion, " (", description, ")\nRelease: ", binaryVersion, "-", build.Commit, "-", build.Time.Unix(), " ",
		runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, "\n",
	))
	if verbose {
		detailedMap := map[string]any{
			"Build Version": binaryVersion,
			"Build Commit":  build.Commit,
			"Build Channel": build.Channel,
		}
		if build.ToolchainId != "" {
			detailedMap["Build Toolchain"] = build.ToolchainId
			detailedMap["Build Time"] = build.Time.Format(time.RFC3339)
		}
		printMap(detailedMap, 2)
	}
}
