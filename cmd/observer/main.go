package main

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/logger"
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
	appStart(parseCommandLine())
}
