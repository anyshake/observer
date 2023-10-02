package logger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func Print(module string, v any, colorCode color.Attribute, carriage bool) {
	color.Set(colorCode)
	if carriage {
		fmt.Print("\r")
	}
	log.Printf("[%s] %v\n", module, v)
	color.Unset()
}
