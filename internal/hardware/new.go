package hardware

import (
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Entry, timeSrc *timesource.Source, actionHandler *action.Handler, explorerOptions explorer.ExplorerOptions, ntpOptions explorer.NtpOptions) (IHardware, error) {
	tr, err := transport.New(explorerOptions.Endpoint, explorerOptions.ReadTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create hardware transport: %w", err)
	}

	channelCodes, err := (&config.StationChannelCodesConfigConstraintImpl{}).Get(actionHandler)
	if err != nil {
		return nil, err
	}
	channelCodesStrArr := channelCodes.([]string)

	switch explorerOptions.Protocol {
	case "legacy":
		explorerOptions.Protocol = "v1"
		fallthrough
	case "v1":
		return &explorer.ExplorerProtoImplV1{
			ChannelCodes:    channelCodesStrArr,
			ExplorerOptions: explorerOptions,
			Transport:       tr,
			Logger:          logger,
			TimeSource:      timeSrc,
			NtpOptions:      ntpOptions,
		}, nil
	case "v2":
		return &explorer.ExplorerProtoImplV2{
			ChannelCodes:    channelCodesStrArr,
			ExplorerOptions: explorerOptions,
			Transport:       tr,
			Logger:          logger,
			TimeSource:      timeSrc,
			NtpOptions:      ntpOptions,
		}, nil
	case "v3":
		return &explorer.ExplorerProtoImplV3{
			ChannelCodes:    channelCodesStrArr,
			ExplorerOptions: explorerOptions,
			Transport:       tr,
			Logger:          logger,
			TimeSource:      timeSrc,
			NtpOptions:      ntpOptions,
		}, nil
	}

	return nil, fmt.Errorf("hardware protocol %s is not supported", explorerOptions.Protocol)
}
