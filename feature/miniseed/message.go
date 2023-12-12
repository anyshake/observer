package miniseed

import (
	"fmt"
	"time"

	"github.com/bclswl0827/mseedio"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
)

func (m *MiniSEED) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions, buffer *miniSEEDBuffer) error {
	var (
		ehz       = gp.EHZ
		ehe       = gp.EHE
		ehn       = gp.EHN
		station   = options.Config.MiniSEED.Station
		network   = options.Config.MiniSEED.Network
		timestamp = time.UnixMilli(gp.TS).UTC()
	)

	// Append EHZ channel to buffer
	buffer.EHZ.DataBuffer = append(buffer.EHZ.DataBuffer, ehz...)
	buffer.EHZ.SampleRate = (buffer.EHZ.SampleRate + int32(len(ehz))) / 2
	// Append EHE channel to buffer
	buffer.EHE.DataBuffer = append(buffer.EHE.DataBuffer, ehe...)
	buffer.EHE.SampleRate = (buffer.EHE.SampleRate + int32(len(ehe))) / 2
	// Append EHN channel to buffer
	buffer.EHN.DataBuffer = append(buffer.EHN.DataBuffer, ehn...)
	buffer.EHN.SampleRate = (buffer.EHN.SampleRate + int32(len(ehn))) / 2

	// Check if buffer is ready to write to file
	if timestamp.Sub(buffer.TimeStamp).Seconds() >= MAX_DURATION {
		// Init MiniSEED data
		var miniseed mseedio.MiniSeedData
		miniseed.Init(ENCODING_TYPE, BIT_ORDER)
		// Set basic data
		filePath := fmt.Sprintf(
			"%s/%s_%s_%s.mseed", buffer.BasePath,
			station, network, timestamp.Format("20060102"),
		)
		// Append channels to MiniSEED
		for _, v := range []string{"EHZ", "EHE", "EHN"} {
			var (
				err error
				seq = fmt.Sprintf("%06d", buffer.SeqNum)
			)
			switch v {
			case "EHZ":
				// Append EHZ channel
				err = miniseed.Append(buffer.EHZ.DataBuffer, &mseedio.AppendOptions{
					ChannelCode:    v,
					SequenceNumber: seq,
					StationCode:    station,
					NetworkCode:    network,
					StartTime:      buffer.TimeStamp,
					SampleRate:     float64(buffer.EHZ.SampleRate),
				})
			case "EHE":
				// Append EHZ channel
				err = miniseed.Append(buffer.EHE.DataBuffer, &mseedio.AppendOptions{
					ChannelCode:    v,
					SequenceNumber: seq,
					StationCode:    station,
					NetworkCode:    network,
					StartTime:      buffer.TimeStamp,
					SampleRate:     float64(buffer.EHE.SampleRate),
				})
			case "EHN":
				// Append EHZ channel
				err = miniseed.Append(buffer.EHN.DataBuffer, &mseedio.AppendOptions{
					ChannelCode:    v,
					SequenceNumber: seq,
					StationCode:    station,
					NetworkCode:    network,
					StartTime:      buffer.TimeStamp,
					SampleRate:     float64(buffer.EHN.SampleRate),
				})
			}
			if err != nil {
				m.OnError(options, err)
				return err
			} else {
				buffer.SeqNum++
			}
			// Encode record to bytes
			dataBytes, err := miniseed.Encode(mseedio.APPEND, BIT_ORDER)
			if err != nil {
				m.OnError(options, err)
				return err
			}
			// Append bytes to file
			err = miniseed.Write(filePath, mseedio.APPEND, dataBytes)
			if err != nil {
				m.OnError(options, err)
				return err
			}
		}
		// Reset buffer
		m.OnReady(options, "write")
		buffer.TimeStamp = timestamp
		buffer.EHZ.DataBuffer = []int32{}
		buffer.EHE.DataBuffer = []int32{}
		buffer.EHN.DataBuffer = []int32{}
	} else {
		m.OnReady(options, "append")
	}

	return nil
}
