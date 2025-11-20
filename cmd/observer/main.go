package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/unibuild"
	"github.com/common-nighthawk/go-figure"
)

func init() {
	copyright := fmt.Sprintf("© SensePlex Limited %d. All Rights Reserved.", time.Now().Year())
	anyshakeLogo := figure.NewFigure("AnyShake", "standard", true).String()
	observerLogo := figure.NewFigure("Observer", "standard", true).String()
	fmt.Printf("%s%s\n%s\n\n", anyshakeLogo, observerLogo, copyright)
	logger.Init()
}

func main() {
	build := unibuild.New(buildToolchain, buildChannel, buildCommit, buildTimestamp)
	ver := semver.New(versionMajor, versionMinor, versionPatch)
	args := parseCommandLine()

	PrintVersion(ver, build, args.showVersion)
	if args.showVersion {
		os.Exit(0)
	}

	if args.upgrade {
		os.Exit(0)
	}

	appStart(ver, build, args)
}
