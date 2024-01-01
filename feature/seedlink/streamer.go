package seedlink

import (
	"fmt"
	"net"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
	"github.com/bclswl0827/mseedio"
)

func (s *SeedLink) Streamer(gp *publisher.Geophone, conn net.Conn, channel, network, station, location string, mode *int) {
	lastTime := time.Now().UTC().UnixMilli()

	for seqNum := 0; *mode == seedlink.WORKINGMODE_STREAM; seqNum++ {
		var (
			ehz   = gp.EHZ
			ehe   = gp.EHE
			ehn   = gp.EHN
			chMap = map[string]publisher.Int32Array{
				"EHZ": ehz, "EHE": ehe, "EHN": ehn,
			}
		)

		// Check if channel exists
		data, ok := chMap[channel]
		if !ok {
			return
		}

		// Check if data is ready
		if gp.TS != lastTime && len(ehz) > 0 && len(ehe) > 0 && len(ehn) > 0 {
			// Generate MiniSEED, send to client
			var miniseed mseedio.MiniSeedData
			// Init header fields
			miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST)
			// Append MiniSEED data
			err := miniseed.Append(data, &mseedio.AppendOptions{
				StationCode:    station,
				LocationCode:   location,
				ChannelCode:    channel,
				NetworkCode:    network,
				SampleRate:     float64(len(data) - 1),
				SequenceNumber: fmt.Sprintf("%06d", seqNum),
				StartTime:      time.Now().UTC().Add(time.Second),
			})
			if err != nil {
				return
			}
			// Get MiniSEED bytes
			dataBytes, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
			if err != nil {
				return
			}
			// Send data with SL sequence number
			dataBytes = append([]byte(fmt.Sprintf("SL%06x", seqNum)), dataBytes...)
			_, err = conn.Write(dataBytes)
			if err != nil {
				return
			}

			lastTime = gp.TS
		}

		time.Sleep(10 * time.Millisecond)
	}
}
