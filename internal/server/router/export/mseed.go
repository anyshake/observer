package export

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/bclswl0827/mseedio"
	"github.com/samber/lo"
)

type seismicDataEncoderMseedImpl struct {
	actionHandler      *action.Handler
	stationCodeConfig  config.StationStationCodeConfigConstraintImpl
	locationCodeConfig config.StationLocationCodeConfigConstraintImpl
	networkCodeConfig  config.StationNetworkCodeConfigConstraintImpl
}

func (e *seismicDataEncoderMseedImpl) GetName() string {
	return "MiniSEED"
}

func (e *seismicDataEncoderMseedImpl) Encode(records []model.SeisRecord, channelCode string) ([]byte, error) {
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

	var miniseed mseedio.MiniSeedData
	err = miniseed.Init(mseedio.INT32, mseedio.MSBFIRST)
	if err != nil {
		return nil, err
	}

	for sequence, record := range records {
		tm, sampleRate, channelDataArr, err := record.Decode()
		if err != nil {
			return nil, err
		}
		channelData, ok := lo.Find(channelDataArr, func(item explorer.ChannelData) bool { return item.ChannelCode == channelCode })
		if !ok {
			continue
		}
		err = miniseed.Append(channelData.Data, &mseedio.AppendOptions{
			SequenceNumber: fmt.Sprintf("%06d", sequence),
			SampleRate:     float64(sampleRate),
			StartTime:      tm.UTC(),
			ChannelCode:    channelCode,
			StationCode:    stationCodeStr,
			NetworkCode:    networkCodeStr,
			LocationCode:   locationCodeStr,
		})
		if err != nil {
			return nil, err
		}
	}

	dataBytes, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (e *seismicDataEncoderMseedImpl) GetFileName(startTime time.Time, channelCode string) (string, error) {
	stationCode, err := e.stationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}
	locationCode, err := e.locationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}
	networkCode, err := e.networkCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}

	stationCodeStr := stationCode.(string)
	locationCodeStr := locationCode.(string)
	networkCodeStr := networkCode.(string)

	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.mseed",
		startTime.UTC().Format("2006"),
		startTime.UTC().Format("002"),
		startTime.UTC().Format("15"),
		startTime.UTC().Format("04"),
		startTime.UTC().Format("05"),
		startTime.UTC().Nanosecond()/1000000,
		stationCodeStr, networkCodeStr,
		locationCodeStr, channelCode,
	)
	return filename, nil
}
