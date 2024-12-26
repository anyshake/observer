package history

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/bclswl0827/mseedio"
)

func (h *History) handleMiniSEED(data []explorer.ExplorerData, stationCode, networkCode, locationCode, channelPrefix, channelCode string) (fileName string, dataBytes []byte, err error) {
	var miniseed mseedio.MiniSeedData
	err = miniseed.Init(mseedio.INT32, mseedio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

	startTime := time.UnixMilli(data[0].Timestamp).UTC()
	channelName := fmt.Sprintf("%s%s", channelPrefix, channelCode)

	for sequence, record := range data {
		switch channelCode {
		case explorer.EXPLORER_CHANNEL_CODE_Z:
			err = miniseed.Append(record.Z_Axis, &mseedio.AppendOptions{
				SequenceNumber: fmt.Sprintf("%06d", sequence),
				SampleRate:     float64(record.SampleRate),
				StartTime:      time.UnixMilli(record.Timestamp).UTC(),
				ChannelCode:    channelName,
				StationCode:    stationCode,
				NetworkCode:    networkCode,
				LocationCode:   locationCode,
			})
		case explorer.EXPLORER_CHANNEL_CODE_E:
			err = miniseed.Append(record.E_Axis, &mseedio.AppendOptions{
				SequenceNumber: fmt.Sprintf("%06d", sequence),
				SampleRate:     float64(record.SampleRate),
				StartTime:      time.UnixMilli(record.Timestamp).UTC(),
				ChannelCode:    channelName,
				StationCode:    stationCode,
				NetworkCode:    networkCode,
				LocationCode:   locationCode,
			})
		case explorer.EXPLORER_CHANNEL_CODE_N:
			err = miniseed.Append(record.N_Axis, &mseedio.AppendOptions{
				SequenceNumber: fmt.Sprintf("%06d", sequence),
				SampleRate:     float64(record.SampleRate),
				StartTime:      time.UnixMilli(record.Timestamp).UTC(),
				ChannelCode:    channelName,
				StationCode:    stationCode,
				NetworkCode:    networkCode,
				LocationCode:   locationCode,
			})
		}
		if err != nil {
			return "", nil, err
		}
	}
	dataBytes, err = miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

	// Return filename and bytes (e.g. 2023.193.14.22.51.0317.AS.SHAKE.00.EHZ.D.mseed)
	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.mseed",
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
	return filename, dataBytes, nil
}
