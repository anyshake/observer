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

func (m *MiniSEED) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions, buffer *miniSEEDBuffer) error {
	var (
		ehz       = gp.EHZ
		ehe       = gp.EHE
		ehn       = gp.EHN
		basePath  = options.Config.MiniSEED.Path
		timestamp = time.UnixMilli(gp.TS).UTC()
		station   = text.TruncateString(options.Config.Station.Station, 5)
		network   = text.TruncateString(options.Config.Station.Network, 2)
		location  = text.TruncateString(options.Config.Station.Location, 2)
	)

	// Append EHZ channel to buffer
	buffer.EHZ.DataBuffer = append(buffer.EHZ.DataBuffer, ehz...)
	buffer.EHZ.Samples += int32(len(ehz))
	// Append EHE channel to buffer
	buffer.EHE.DataBuffer = append(buffer.EHE.DataBuffer, ehe...)
	buffer.EHE.Samples += int32(len(ehe))
	// Append EHN channel to buffer
	buffer.EHN.DataBuffer = append(buffer.EHN.DataBuffer, ehn...)
	buffer.EHN.Samples += int32(len(ehn))

	// Check if buffer is ready to write to file
	timeDiffSec := timestamp.Sub(buffer.TimeStamp).Seconds()
	if timeDiffSec >= MAX_DURATION {
		// Init MiniSEED data
		var miniseed mseedio.MiniSeedData
		miniseed.Init(ENCODING_TYPE, BIT_ORDER)
		// Get sequence number in string
		seqNum := fmt.Sprintf("%06d", buffer.SeqNum)
		buffer.SeqNum++
		// Append channels to MiniSEED
		for i, v := range map[string]*channelBuffer{
			"EHZ": buffer.EHZ,
			"EHE": buffer.EHE,
			"EHN": buffer.EHN,
		} {
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
			filePath := getFilePath(basePath, station, network, location, timestamp)
			err = miniseed.Write(filePath, mseedio.APPEND, dataBytes)
			if err != nil {
				m.OnError(options, err)
				return err
			}
		}
		// Reset buffer
		m.OnReady(options, "write")
		buffer.TimeStamp = timestamp
		buffer.EHZ = &channelBuffer{}
		buffer.EHE = &channelBuffer{}
		buffer.EHN = &channelBuffer{}
	} else {
		m.OnReady(options, "append")
	}

	return nil
}
