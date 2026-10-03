package export

import (
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
)

const LOG_PREFIX = "restful_api_export"

type seismicRecordIterator func(func(model.SeisRecord) error) error

type seismicDataEncoder interface {
	GetName() string
	Encode(records seismicRecordIterator, channel string) ([]byte, error)
	GetFileName(startTime time.Time, channelCode string) (string, error)
}

type channelRecordRange struct {
	Count          int
	SampleRate     int
	StartTimestamp int64
	EndTimestamp   int64
}

func forEachContinuousChannelRecord(
	records seismicRecordIterator,
	channelCode string,
	callback func(record model.SeisRecord, samples []int32) error,
) (channelRecordRange, error) {
	var result channelRecordRange
	err := records(func(record model.SeisRecord) error {
		recordIndex := result.Count
		if recordIndex == 0 {
			result.SampleRate = record.SampleRate
			result.StartTimestamp = record.RecordTime
		}
		result.Count++
		result.EndTimestamp = record.RecordTime

		_, _, channelData, err := record.Decode()
		if err != nil {
			return err
		}

		var samples []int32
		channelFound := false
		for i := range channelData {
			if channelData[i].ChannelCode == channelCode {
				channelFound = true
				samples = channelData[i].Data
				break
			}
		}
		if !channelFound {
			return nil
		}

		expectedTimestamp := result.StartTimestamp + int64(recordIndex*1000)
		if math.Abs(float64(record.RecordTime-expectedTimestamp)) >= explorer.ALLOWED_JITTER_MS_NTP {
			return fmt.Errorf(
				"timestamp is not within allowed jitter %d ms, expected %d, got %d",
				explorer.ALLOWED_JITTER_MS_NTP,
				expectedTimestamp,
				record.RecordTime,
			)
		}
		if record.SampleRate != result.SampleRate {
			return fmt.Errorf("sample rate is not the same, expected %d, got %d", result.SampleRate, record.SampleRate)
		}

		return callback(record, samples)
	})
	return result, err
}
