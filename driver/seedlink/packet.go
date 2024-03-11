package seedlink

import (
	"fmt"
	"net"
	"time"

	"github.com/bclswl0827/mseedio"
)

func SendSLPacket(conn net.Conn, client *SeedLinkClient, data SeedLinkPacket) error {
	// Create data chunks to adapt to SeedLink packet size
	var countGroup [][]int32
	if len(data.Count) > CHUNK_SIZE {
		for i := 0; i < len(data.Count); i += CHUNK_SIZE {
			if i+CHUNK_SIZE > len(data.Count) {
				countGroup = append(countGroup, data.Count[i:])
			} else {
				countGroup = append(countGroup, data.Count[i:i+CHUNK_SIZE])
			}
		}
	} else {
		countGroup = append(countGroup, data.Count)
	}

	dataSpanMs := 1000 / float64(len(data.Count))
	for i, c := range countGroup {
		// Generate MiniSEED record
		var miniseed mseedio.MiniSeedData
		miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST)
		err := miniseed.Append(c, &mseedio.AppendOptions{
			ChannelCode:    data.Channel,
			StationCode:    client.Station,
			LocationCode:   client.Location,
			NetworkCode:    client.Network,
			SampleRate:     float64(len(data.Count)),
			SequenceNumber: fmt.Sprintf("%06d", client.Sequence),
			StartTime:      time.UnixMilli(data.Timestamp + int64(float64(i*CHUNK_SIZE)*dataSpanMs)).UTC(),
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
		slSeq := []byte(fmt.Sprintf("SL%06X", client.Sequence))
		_, err = conn.Write(slSeq)
		if err != nil {
			return err
		}

		// Send SeedLink packet data
		_, err = conn.Write(slData)
		if err != nil {
			return err
		}

		client.Sequence++
	}

	return nil
}
