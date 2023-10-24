package ntpclient

import (
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/bclswl0827/observer/utils/text"
	"github.com/fatih/color"
)

func (n *NTPClient) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta, false)
}

func (n *NTPClient) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue, false)
}

func (n *NTPClient) OnReady(options *feature.FeatureOptions, v ...any) {
	options.Status.System.Offset = v[0].(float64)
	options.Status.IsReady = true
	logger.Print(MODULE, "time alignment succeed", color.FgGreen, false)
}

func (n *NTPClient) OnError(options *feature.FeatureOptions, err error) {
	logger.Print(MODULE, err, color.FgRed, false)
}
