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

func New(endpoint, protocol, model string, timeout int, latitude, longitude, elevation float64, actionHandler *action.Handler, timeSource *timesource.Source, logger *logrus.Entry) (IHardware, error) {
	tr, err := transport.New(endpoint, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create hardware transport: %w", err)
	}

	channelCodes, err := (&config.StationChannelCodesConfigConstraintImpl{}).Get(actionHandler)
	if err != nil {
		return nil, err
	}
	channelCodesStrArr := channelCodes.([]string)

	switch protocol {
	case "legacy":
		fallthrough
	case "v1":
		return &explorer.ExplorerProtoImplV1{
			ChannelCodes:      channelCodesStrArr,
			Model:             model,
			Transport:         tr,
			Logger:            logger,
			TimeSource:        timeSource,
			FallbackLatitude:  latitude,
			FallbackLongitude: longitude,
			FallbackElevation: elevation,
		}, nil
	case "v2":
		return &explorer.ExplorerProtoImplV2{
			ChannelCodes:      channelCodesStrArr,
			Model:             model,
			Transport:         tr,
			Logger:            logger,
			TimeSource:        timeSource,
			FallbackLatitude:  latitude,
			FallbackLongitude: longitude,
			FallbackElevation: elevation,
		}, nil
	case "v3":
		return &explorer.ExplorerProtoImplV3{
			ChannelCodes:      channelCodesStrArr,
			Model:             model,
			Transport:         tr,
			Logger:            logger,
			TimeSource:        timeSource,
			FallbackLatitude:  latitude,
			FallbackLongitude: longitude,
			FallbackElevation: elevation,
		}, nil
	}

	return nil, fmt.Errorf("hardware protocol %s is not supported", protocol)
}
