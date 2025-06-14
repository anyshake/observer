package export

import (
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/bclswl0827/sacio"
	"github.com/samber/lo"
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

func (e *seismicDataEncoderSacImpl) Encode(records []model.SeisRecord, channelCode string) ([]byte, error) {
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

	var (
		startSampleRate = records[0].SampleRate
		startTimestamp  = records[0].RecordTime
	)

	var channelBuffer []int32
	for index, record := range records {
		_, _, channelDataArr, err := record.Decode()
		if err != nil {
			return nil, err
		}
		channelData, ok := lo.Find(channelDataArr, func(item explorer.ChannelData) bool { return item.ChannelCode == channelCode })
		if !ok {
			continue
		}

		// Make sure timestamp is continuous
		if math.Abs(float64(record.RecordTime-startTimestamp-int64(index*1000))) >= explorer.ALLOWED_JITTER_MS {
			return nil, fmt.Errorf(
				"timestamp is not within allowed jitter %d ms, expected %d, got %d",
				explorer.ALLOWED_JITTER_MS,
				startTimestamp+int64(index*1000),
				record.RecordTime,
			)
		}

		// Make sure sample rate is the same
		if record.SampleRate != startSampleRate {
			return nil, fmt.Errorf("sample rate is not the same, expected %d, got %d", startSampleRate, record.SampleRate)
		}

		channelBuffer = append(channelBuffer, channelData.Data...)
	}

	var sac sacio.SACData
	if err = sac.Init(); err != nil {
		return nil, err
	}
	startTime, _, _, err := records[0].Decode()
	if err != nil {
		return nil, err
	}
	endTime, _, _, err := records[len(records)-1].Decode()
	if err != nil {
		return nil, err
	}
	sac.SetTime(startTime.UTC(), endTime.Sub(startTime))
	sac.SetInfo(networkCodeStr, stationCodeStr, locationCodeStr, channelCode)
	sac.SetBody(e.int32ToFloat32(channelBuffer), startSampleRate)

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

func (e *seismicDataEncoderSacImpl) int32ToFloat32(arr []int32) []float32 {
	floatSlice := make([]float32, len(arr))
	for i, num := range arr {
		floatSlice[i] = float32(num)
	}
	return floatSlice
}
