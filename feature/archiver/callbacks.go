package archiver

import (
	"github.com/bclswl0827/observer/driver/postgres"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/bclswl0827/observer/utils/text"
	"github.com/fatih/color"
)

func (a *Archiver) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta, false)
}

func (a *Archiver) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue, false)
}

func (a *Archiver) OnReady(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, "1 message has been archived", color.FgGreen, false)
}

func (a *Archiver) OnError(options *feature.FeatureOptions, err error) {
	postgres.Close(options.Database)
	logger.Print(MODULE, err, color.FgRed, false)
}
