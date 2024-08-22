package history

import (
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/bclswl0827/sacio"
)

func (h *History) getSACBytes(data []explorer.ExplorerData, legacyMode bool, stationCode, networkCode, locationCode, channelPrefix, channelCode string) (string, []byte, error) {
	var (
		startSampleRate = data[0].SampleRate
		startTimestamp  = data[0].Timestamp
		startTime       = time.UnixMilli(startTimestamp).UTC()
		endTime         = time.UnixMilli(data[len(data)-1].Timestamp).UTC()
		channelName     = fmt.Sprintf("%s%s", channelPrefix, channelCode)
	)

	var channelBuffer []int32
	for index, record := range data {
		// Make sure timestamp is continuous
		if math.Abs(float64(record.Timestamp-startTimestamp-int64(index*1000))) >= explorer.EXPLORER_ALLOWED_JITTER_MS {
			return "", nil, fmt.Errorf(
				"timestamp is not within allowed jitter %d ms, expected %d, got %d",
				explorer.EXPLORER_ALLOWED_JITTER_MS,
				startTimestamp+int64(index*1000),
				record.Timestamp,
			)
		}

		// Make sure sample rate is the same
		if record.SampleRate != startSampleRate {
			return "", nil, fmt.Errorf("sample rate is not the same, expected %d, got %d", startSampleRate, record.SampleRate)
		}

		switch channelCode {
		case explorer.EXPLORER_CHANNEL_CODE_Z:
			channelBuffer = append(channelBuffer, record.Z_Axis...)
		case explorer.EXPLORER_CHANNEL_CODE_E:
			channelBuffer = append(channelBuffer, record.E_Axis...)
		case explorer.EXPLORER_CHANNEL_CODE_N:
			channelBuffer = append(channelBuffer, record.N_Axis...)
		}
	}

	var sac sacio.SACData
	err := sac.Init()
	if err != nil {
		return "", nil, err
	}
	sac.SetTime(startTime, endTime.Sub(startTime))
	sac.SetInfo(networkCode, stationCode, locationCode, channelName)
	sac.SetBody(h.int32ToFloat32(channelBuffer), startSampleRate)

	// Return filename and bytes (e.g. 2023.193.14.22.51.0317.AS.SHAKE.00.EHZ.D.sac)
	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.sac",
		startTime.Format("2006"),
		startTime.Format("002"),
		startTime.Format("15"),
		startTime.Format("04"),
		startTime.Format("05"),
		// Get the current millisecond
		startTime.Nanosecond()/1000000,
		stationCode, networkCode,
		locationCode, channelName,
	)
	dataBytes, err := sac.Encode(sacio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

	return filename, dataBytes, nil
}

func (h *History) int32ToFloat32(arr []int32) []float32 {
	floatSlice := make([]float32, len(arr))
	for i, num := range arr {
		floatSlice[i] = float32(num)
	}
	return floatSlice
}
