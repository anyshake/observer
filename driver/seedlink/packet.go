package seedlink

import (
	"fmt"
	"time"

	"github.com/bclswl0827/mseedio"
)

func CreateSLPacket(count []int32, ts, seq int64, network, station, channel, location string) ([]byte, error) {
	// Generate MiniSEED, send to client
	var miniseed mseedio.MiniSeedData
	// Init header fields
	miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST)

	// Append MiniSEED data
	err := miniseed.Append(count, &mseedio.AppendOptions{
		StationCode:    station,
		LocationCode:   location,
		ChannelCode:    channel,
		NetworkCode:    network,
		SampleRate:     float64(len(count)),
		StartTime:      time.UnixMilli(ts).UTC(),
		SequenceNumber: fmt.Sprintf("%06d", seq),
	})
	if err != nil {
		return nil, err
	}

	// Get MiniSEED bytes
	miniseed.Series[0].BlocketteSection.RecordLength = 9
	dataBytes, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
	if err != nil {
		return nil, err
	}

	// Return SeedLink packet
	slSeq := []byte(fmt.Sprintf("SL%06X", seq))
	return append(slSeq, dataBytes...), nil
}
