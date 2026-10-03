package export

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/bclswl0827/sacio"
)

type seismicDataEncoderSacImpl struct {
	actionHandler      *action.Handler
	stationCodeConfig  config.StationStationCodeConfigConstraintImpl
	locationCodeConfig config.StationLocationCodeConfigConstraintImpl
	networkCodeConfig  config.StationNetworkCodeConfigConstraintImpl
}

func (e *seismicDataEncoderSacImpl) GetName() string {
	return "SAC"
}

func (e *seismicDataEncoderSacImpl) Encode(records seismicRecordIterator, channelCode string) ([]byte, error) {
	stationCode, err := e.stationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return nil, err
	}
	locationCode, err := e.locationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return nil, err
	}
	networkCode, err := e.networkCodeConfig.Get(e.actionHandler)
	if err != nil {
		return nil, err
	}

	stationCodeStr := stationCode.(string)
	locationCodeStr := locationCode.(string)
	networkCodeStr := networkCode.(string)

	var channelBuffer []float32
	recordRange, err := forEachContinuousChannelRecord(records, channelCode, func(_ model.SeisRecord, samples []int32) error {
		for _, sample := range samples {
			channelBuffer = append(channelBuffer, float32(sample))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if recordRange.Count == 0 {
		return nil, nil
	}

	var sac sacio.SACData
	if err = sac.Init(); err != nil {
		return nil, err
	}
	startTime := time.UnixMilli(recordRange.StartTimestamp)
	endTime := time.UnixMilli(recordRange.EndTimestamp)
	sac.SetTime(startTime.UTC(), endTime.Sub(startTime))
	sac.SetInfo(networkCodeStr, stationCodeStr, locationCodeStr, channelCode)
	sac.SetBody(channelBuffer, recordRange.SampleRate)

	dataBytes, err := sac.Encode(sacio.MSBFIRST)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (e *seismicDataEncoderSacImpl) GetFileName(startTime time.Time, channelCode string) (string, error) {
	stationCode, err := e.stationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}
	locationCode, err := e.stationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}
	networkCode, err := e.networkCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.sac",
		startTime.UTC().Format("2006"),
		startTime.UTC().Format("002"),
		startTime.UTC().Format("15"),
		startTime.UTC().Format("04"),
		startTime.UTC().Format("05"),
		startTime.UTC().Nanosecond()/1000000,
		stationCode, networkCode,
		locationCode, channelCode,
	)
	return filename, nil
}
