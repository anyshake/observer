package ntpclient

import (
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/duration"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
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
	options.Status.ReadyTime, _ = duration.Timestamp(options.Status.System.Offset)
	logger.Print(MODULE, "time alignment succeed", color.FgGreen, false)
}

func (n *NTPClient) OnError(options *feature.FeatureOptions, err error) {
	logger.Print(MODULE, err, color.FgRed, false)
}
