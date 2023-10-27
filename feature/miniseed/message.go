package miniseed

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bclswl0827/mseedio"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
)

func (m *MiniSEED) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions, seqNumber int, basePath, station, network string) error {
	var (
		ehz = gp.EHZ
		ehe = gp.EHE
		ehn = gp.EHN
		ts  = time.UnixMilli(gp.TS).UTC()
	)

	// Init MiniSEED library
	var miniseed mseedio.MiniSeedData
	miniseed.Init(mseedio.INT32, mseedio.MSBFIRST)

	// Get file name by date
	filePath := fmt.Sprintf(
		"%s/%s_%s_%s.mseed",
		basePath, station, network,
		ts.Format("20060102"),
	)

	// If file exists, check sequence number
	_, err := os.Stat(filePath)
	if err == nil && seqNumber == 0 {
		// Read MiniSEED file
		var ms mseedio.MiniSeedData
		err := ms.Read(filePath)
		if err != nil {
			m.OnError(options, err)
			return err
		}

		// Get last sequence number
		recordLength := len(ms.Series)
		if recordLength > 0 {
			lastRecord := ms.Series[recordLength-1]
			n, err := strconv.Atoi(lastRecord.FixedSection.SequenceNumber)
			if err != nil {
				m.OnError(options, err)
				return err
			}

			// Set current sequence number
			seqNumber = n
		}
	}

	// Increments sequence number by 1
	if seqNumber >= 999999 {
		seqNumber = 0
	} else {
		seqNumber++
	}
	seqNumberString := fmt.Sprintf("%06d", seqNumber)

	// Append 3 channels
	for i, v := range [][]int32{ehz, ehe, ehn} {
		var err error
		switch i {
		case 0:
			err = miniseed.Append(v, mseedio.AppendOptions{
				StartTime:      ts,
				ChannelCode:    "EHZ",
				StationCode:    station,
				NetworkCode:    network,
				SequenceNumber: seqNumberString,
				SampleRate:     float64(len(ehz) - 1),
			})
		case 1:
			err = miniseed.Append(v, mseedio.AppendOptions{
				StartTime:      ts,
				ChannelCode:    "EHE",
				StationCode:    station,
				NetworkCode:    network,
				SequenceNumber: seqNumberString,
				SampleRate:     float64(len(ehe) - 1),
			})
		case 2:
			err = miniseed.Append(v, mseedio.AppendOptions{
				StartTime:      ts,
				ChannelCode:    "EHN",
				StationCode:    station,
				NetworkCode:    network,
				SequenceNumber: seqNumberString,
				SampleRate:     float64(len(ehn) - 1),
			})
		}
		if err != nil {
			m.OnError(options, err)
			return err
		}

		// Encode record to bytes
		dataBytes, err := miniseed.Encode(mseedio.APPEND, mseedio.MSBFIRST)
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

	m.OnReady(options)
	return nil
}
