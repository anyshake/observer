package seedlink

import (
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
	"github.com/fatih/color"
)

func (s *SeedLink) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta, false)
}

func (s *SeedLink) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue, false)
}

func (s *SeedLink) OnReady(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgGreen, false)
}

func (s *SeedLink) OnError(options *feature.FeatureOptions, err error) {
	logger.Print(MODULE, err, color.FgRed, false)
}
