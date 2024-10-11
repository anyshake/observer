package miniseed

import (
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/bclswl0827/mseedio"
)

func (m *MiniSeedService) handleWrite() error {
	var (
		startTimestamp  = m.miniseedBuffer[0].Timestamp
		startSampleRate = m.miniseedBuffer[0].SampleRate
	)

	for i := int64(1); i < int64(len(m.miniseedBuffer)); i++ {
		// Make sure timestamp is increasing by 1000 ms with allowed jitter
		if math.Abs(float64(m.miniseedBuffer[i].Timestamp-startTimestamp-(i*time.Second.Milliseconds()))) >= explorer.EXPLORER_ALLOWED_JITTER_MS {
			return fmt.Errorf(
				"timestamp is not within allowed jitter %d ms, expected %d, got %d",
				explorer.EXPLORER_ALLOWED_JITTER_MS,
				startTimestamp+(i*time.Second.Milliseconds()),
				m.miniseedBuffer[i].Timestamp,
			)
		}

		// Make sure sample rate is the same
		if m.miniseedBuffer[i].SampleRate != startSampleRate {
			return fmt.Errorf("sample rate is not the same, expected %d, got %d", startSampleRate, m.miniseedBuffer[i].SampleRate)
		}
	}

	// Write data to file by channels
	for _, channelCode := range []string{
		explorer.EXPLORER_CHANNEL_CODE_Z,
		explorer.EXPLORER_CHANNEL_CODE_E,
		explorer.EXPLORER_CHANNEL_CODE_N,
	} {
		var channelBuffer []int32
		for i := 0; i < len(m.miniseedBuffer); i++ {
			switch channelCode {
			case explorer.EXPLORER_CHANNEL_CODE_Z:
				channelBuffer = append(channelBuffer, m.miniseedBuffer[i].Z_Axis...)
			case explorer.EXPLORER_CHANNEL_CODE_E:
				channelBuffer = append(channelBuffer, m.miniseedBuffer[i].E_Axis...)
			case explorer.EXPLORER_CHANNEL_CODE_N:
				channelBuffer = append(channelBuffer, m.miniseedBuffer[i].N_Axis...)
			}
		}
		var miniseed mseedio.MiniSeedData
		err := miniseed.Init(MINISEED_ENCODE_TYPE, MINISEED_BIT_ORDER)
		if err != nil {
			return err
		}
		channelName := fmt.Sprintf("%s%s", m.channelPrefix, channelCode)
		startTime := time.UnixMilli(startTimestamp).UTC()
		err = miniseed.Append(channelBuffer, &mseedio.AppendOptions{
			SequenceNumber: fmt.Sprintf("%06d", m.miniseedSequence[channelCode]),
			SampleRate:     float64(startSampleRate),
			StartTime:      startTime,
			ChannelCode:    channelName,
			StationCode:    m.stationCode,
			NetworkCode:    m.networkCode,
			LocationCode:   m.locationCode,
		})
		if err != nil {
			return err
		}
		dataBytes, err := miniseed.Encode(mseedio.APPEND, MINISEED_BIT_ORDER)
		if err != nil {
			return err
		}
		filePath := m.getFilePath(channelName, startTime)
		err = miniseed.Write(filePath, mseedio.APPEND, dataBytes)
		if err != nil {
			return err
		}

		m.miniseedSequence[channelCode]++
	}

	return nil
}
