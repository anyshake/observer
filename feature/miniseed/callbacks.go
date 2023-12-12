package miniseed

import (
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
	"github.com/fatih/color"
)

func (m *MiniSEED) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta, false)
}

func (m *MiniSEED) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue, false)
}

func (m *MiniSEED) OnReady(options *feature.FeatureOptions, v ...any) {
	switch v[0].(string) {
	case "append":
		logger.Print(MODULE, "1 record has been append", color.FgGreen, false)
	case "write":
		logger.Print(MODULE, "1 record has been written", color.FgGreen, false)
	}
}

func (m *MiniSEED) OnError(options *feature.FeatureOptions, err error) {
	logger.Print(MODULE, err, color.FgRed, false)
}
