package logger

import (
	"log"

	"github.com/fatih/color"
)

func Print(module string, v any, colorCode color.Attribute) {
	color.Set(colorCode)
	log.Printf("[%s] %v\n", module, v)
	color.Unset()
}
