package callbacks

import (
	"com.geophone.observer/common/postgres"
	"com.geophone.observer/feature"
	"com.geophone.observer/utils/logger"
	"github.com/fatih/color"
)

func OnError(module string, options *feature.FeatureOptions, err error) {
	switch module {
	case "geophone":
		options.Status.System.Errors++
	case "archiver":
		postgres.Close(options.Database)
	}

	logger.Print(module, err, color.FgRed)
}
