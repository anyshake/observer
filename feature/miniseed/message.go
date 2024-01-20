package miniseed

import (
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/text"
	"github.com/bclswl0827/mseedio"
)

func (m *MiniSEED) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions, buffer *publisher.SegmentBuffer) error {
	var (
		basePath   = options.Config.MiniSEED.Path
		timestamp  = time.UnixMilli(gp.TS).UTC()
		station    = text.TruncateString(options.Config.Station.Station, 5)
		network    = text.TruncateString(options.Config.Station.Network, 2)
		location   = text.TruncateString(options.Config.Station.Location, 2)
		channelMap = map[string]publisher.Int32Array{
			"EHZ": gp.EHZ, "EHE": gp.EHE, "EHN": gp.EHN,
		}
	)

	// Append geophone channel data to buffer
	for i, v := range buffer.ChannelBuffer {
		channelData, ok := channelMap[i]
		if ok {
			v.DataBuffer = append(v.DataBuffer, channelData...)
			v.Samples += int32(len(channelData))
		}
	}

	// Check if buffer is ready to write to file
	timeDiffSec := timestamp.Sub(buffer.TimeStamp).Seconds()
	if timeDiffSec >= MAX_DURATION {
		// Append channels to MiniSEED
		for i, v := range buffer.ChannelBuffer {
			// Init MiniSEED data
			var miniseed mseedio.MiniSeedData
			miniseed.Init(ENCODING_TYPE, BIT_ORDER)
			// Get sequence number in string
			seqNum := fmt.Sprintf("%06d", v.SeqNum)
			v.SeqNum++
			// Get sample rate
			sampleRate := math.Round(float64(v.Samples) / timeDiffSec)
			// Append channel data
			err := miniseed.Append(v.DataBuffer, &mseedio.AppendOptions{
				ChannelCode:    i,
				SequenceNumber: seqNum,
				StationCode:    station,
				NetworkCode:    network,
				LocationCode:   location,
				StartTime:      timestamp,
				SampleRate:     sampleRate,
			})
			if err != nil {
				m.OnError(options, err)
				return err
			}
			// Encode record to bytes
			dataBytes, err := miniseed.Encode(mseedio.APPEND, BIT_ORDER)
			if err != nil {
				m.OnError(options, err)
				return err
			}
			// Append bytes to file
			filePath := getFilePath(basePath, station, network, location, i, timestamp)
			err = miniseed.Write(filePath, mseedio.APPEND, dataBytes)
			if err != nil {
				m.OnError(options, err)
				return err
			}
		}

		// Reset buffer
		m.OnReady(options, "write")
		buffer.TimeStamp = timestamp
		for _, v := range buffer.ChannelBuffer {
			v.DataBuffer = []int32{}
			v.Samples = 0
		}
	} else {
		m.OnReady(options, "append")
	}

	return nil
}
