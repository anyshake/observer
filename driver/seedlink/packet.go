package seedlink

import (
	"fmt"
	"net"
	"time"

	"github.com/bclswl0827/mseedio"
)

func SendSLPacket(conn net.Conn, count []int32, ts int64, seq *int64, network, station, channel, location string) error {
	// Create data chunks to adapt to SeedLink packet size
	var countGroup [][]int32
	if len(count) > CHUNK_SIZE {
		for i := 0; i < len(count); i += CHUNK_SIZE {
			if i+CHUNK_SIZE > len(count) {
				countGroup = append(countGroup, count[i:])
			} else {
				countGroup = append(countGroup, count[i:i+CHUNK_SIZE])
			}
		}
	} else {
		countGroup = append(countGroup, count)
	}

	dataSpanMs := 1000 / float64(len(count))
	for i, c := range countGroup {
		// Generate MiniSEED record
		var miniseed mseedio.MiniSeedData
		miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST)
		err := miniseed.Append(c, &mseedio.AppendOptions{
			StationCode:    station,
			LocationCode:   location,
			ChannelCode:    channel,
			NetworkCode:    network,
			SampleRate:     float64(len(count)),
			SequenceNumber: fmt.Sprintf("%06d", *seq),
			StartTime:      time.UnixMilli(ts + int64(float64(i*CHUNK_SIZE)*dataSpanMs)).UTC(),
		})
		if err != nil {
			return err
		}

		// Get MiniSEED data bytes always in 512 bytes
		miniseed.Series[0].BlocketteSection.RecordLength = 9
		slData, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
		if err != nil {
			return err
		}

		// Prepend and send SeedLink sequence number
		slSeq := []byte(fmt.Sprintf("SL%06X", *seq))
		_, err = conn.Write(slSeq)
		if err != nil {
			return err
		}

		// Send SeedLink packet data
		_, err = conn.Write(slData)
		if err != nil {
			return err
		}

		*seq++
	}

	return nil
}
