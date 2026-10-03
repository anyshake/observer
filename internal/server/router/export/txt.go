package export

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
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

func (e *seismicDataEncoderTxtImpl) Encode(records seismicRecordIterator, channelCode string) ([]byte, error) {
	var buffer bytes.Buffer
	line := make([]byte, 0, 48)
	recordRange, err := forEachContinuousChannelRecord(records, channelCode, func(record model.SeisRecord, samples []int32) error {
		sampleSpanMs := 1000.0 / float64(record.SampleRate)
		for i, sample := range samples {
			timestampMs := float64(record.RecordTime) + sampleSpanMs*float64(i)
			line = strconv.AppendFloat(line[:0], timestampMs, 'f', 0, 64)
			line = append(line, ' ')
			line = strconv.AppendInt(line, int64(sample), 10)
			line = append(line, '\n')
			_, _ = buffer.Write(line)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if recordRange.Count == 0 {
		return nil, nil
	}

	return buffer.Bytes(), nil
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
