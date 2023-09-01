package miniseed

import (
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/bclswl0827/observer/utils/text"
	"github.com/fatih/color"
)

func (g *MiniSEED) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta)
}

func (g *MiniSEED) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue)
}

func (a *MiniSEED) OnReady(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, "1 record has been written", color.FgGreen)
}

func (m *MiniSEED) OnError(options *feature.FeatureOptions, err error) {
	logger.Print(MODULE, err, color.FgRed)
}
