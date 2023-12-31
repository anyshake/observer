package logger

import (
	"os"

	"github.com/fatih/color"
)

func Fatal(module string, v any, colorCode color.Attribute) {
	Print(module, v, colorCode, false)
	os.Exit(1)
}
