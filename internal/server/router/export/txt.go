package export

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
)

type seismicDataEncoderTxtImpl struct {
	actionHandler      *action.Handler
	stationCodeConfig  config.StationStationCodeConfigConstraintImpl
	locationCodeConfig config.StationLocationCodeConfigConstraintImpl
	networkCodeConfig  config.StationNetworkCodeConfigConstraintImpl
}

func (e *seismicDataEncoderTxtImpl) GetName() string {
	return "TXT"
}

func (e *seismicDataEncoderTxtImpl) Encode(records []model.SeisRecord, channelCode string) ([]byte, error) {
	var builder strings.Builder
	builder.Grow(1024 * len(records))

	var (
		startSampleRate = records[0].SampleRate
		startTimestamp  = records[0].Timestamp
	)

	for index, record := range records {
		_, _, channelDataArr, err := record.Decode()
		if err != nil {
			return nil, err
		}

		var channelData *explorer.ChannelData
		for i := range channelDataArr {
			if channelDataArr[i].ChannelCode == channelCode {
				channelData = &channelDataArr[i]
				break
			}
		}
		if channelData == nil {
			continue
		}

		// Make sure timestamp is continuous
		if math.Abs(float64(record.Timestamp-startTimestamp-int64(index*1000))) >= explorer.ALLOWED_JITTER_MS {
			return nil, fmt.Errorf(
				"timestamp is not within allowed jitter %d ms, expected %d, got %d",
				explorer.ALLOWED_JITTER_MS,
				startTimestamp+int64(index*1000),
				record.Timestamp,
			)
		}

		// Make sure sample rate is the same
		if record.SampleRate != startSampleRate {
			return nil, fmt.Errorf("sample rate is not the same, expected %d, got %d", startSampleRate, record.SampleRate)
		}

		sampleSpanMs := 1000.0 / float64(record.SampleRate)
		for i, v := range channelData.Data {
			timestampMs := float64(record.Timestamp) + sampleSpanMs*float64(i)
			builder.WriteString(strconv.FormatFloat(timestampMs, 'f', 0, 64))
			builder.WriteByte(' ')
			builder.WriteString(strconv.Itoa(int(v)))
			builder.WriteByte('\n')
		}
	}

	return []byte(builder.String()), nil
}

func (e *seismicDataEncoderTxtImpl) GetFileName(startTime time.Time, channelCode string) (string, error) {
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

	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.txt",
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
