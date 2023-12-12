package archiver

import (
	"github.com/anyshake/observer/driver/dao"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
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
	dao.Close(options.Database)
	logger.Print(MODULE, err, color.FgRed, false)
}
