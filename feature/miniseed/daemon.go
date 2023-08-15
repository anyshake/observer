package miniseed

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"com.geophone.observer/feature"
	"com.geophone.observer/handler"
	"github.com/bclswl0827/mseedio"
)

func (m *MiniSEED) Start(options *feature.FeatureOptions) {
	if !options.Config.MiniSEED.Enable {
		options.OnStop(MODULE, options, "service is disabled")
		return
	}

	// Init sequence number
	var seqNumber int
	options.OnStart(MODULE, options, "service has started")

	// Append and write when new message arrived
	handler.OnMessage(&options.Status.Geophone,
		func(gp *handler.Geophone) error {
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
				"%s/%s.mseed",
				options.Config.MiniSEED.Path,
				ts.Format("2006-01-02"),
			)

			// If file exists, check sequence number
			_, err := os.Stat(filePath)
			if err == nil && seqNumber == 0 {
				// Read MiniSEED file
				var ms mseedio.MiniSeedData
				err := ms.Read(filePath)
				if err != nil {
					options.OnError(MODULE, options, err)
					return err
				}

				// Get last sequence number
				recordLength := len(ms.Series)
				if recordLength > 0 {
					lastRecord := ms.Series[recordLength-1]
					n, err := strconv.Atoi(lastRecord.FixedSection.SequenceNumber)
					if err != nil {
						options.OnError(MODULE, options, err)
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
			var (
				station = options.Config.MiniSEED.Station
				network = options.Config.MiniSEED.Network
			)
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
						SampleRate:     float64(len(ehz)) - 0.1,
					})
				case 1:
					err = miniseed.Append(v, mseedio.AppendOptions{
						StartTime:      ts,
						ChannelCode:    "EHE",
						StationCode:    station,
						NetworkCode:    network,
						SequenceNumber: seqNumberString,
						SampleRate:     float64(len(ehe)) - 0.1,
					})
				case 2:
					err = miniseed.Append(v, mseedio.AppendOptions{
						StartTime:      ts,
						ChannelCode:    "EHN",
						StationCode:    station,
						NetworkCode:    network,
						SequenceNumber: seqNumberString,
						SampleRate:     float64(len(ehn)) - 0.1,
					})
				}
				if err != nil {
					options.OnError(MODULE, options, err)
					return err
				}

				// Encode record to bytes
				dataBytes, err := miniseed.Encode(mseedio.APPEND, mseedio.MSBFIRST)
				if err != nil {
					options.OnError(MODULE, options, err)
					return err
				}

				// Append bytes to file
				err = miniseed.Write(filePath, mseedio.APPEND, dataBytes)
				if err != nil {
					options.OnError(MODULE, options, err)
					return err
				}
			}

			options.OnReady(MODULE, options)
			return nil
		},
	)

	err := fmt.Errorf("service exited with a error")
	options.OnError(MODULE, options, err)
}
