package history

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/bclswl0827/mseedio"
)

func (h *History) getMiniSeedBytes(data []explorer.ExplorerData, stationCode, networkCode, locationCode, channelPrefix, channelCode string) (string, []byte, error) {
	var miniseed mseedio.MiniSeedData
	err := miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

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
	dataBytes, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

	return "query.mseed", dataBytes, nil
}
