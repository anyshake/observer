package seedlink

import (
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/bclswl0827/observer/utils/text"
	"github.com/fatih/color"
)

func (s *SeedLink) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta, false)
}

func (s *SeedLink) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue, false)
}

func (s *SeedLink) OnReady(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, "1 record has been append", color.FgGreen, false)
}

func (s *SeedLink) OnError(options *feature.FeatureOptions, err error) {
	logger.Print(MODULE, err, color.FgRed, false)
}