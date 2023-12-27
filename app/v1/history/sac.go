package history

import (
	"fmt"
	"sort"
	"time"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/duration"
	"github.com/bclswl0827/sacio"
)

func getSACBytes(data []publisher.Geophone, channel string, options *app.ServerOptions) (string, []byte, error) {
	// Check channel
	if channel != "EHZ" && channel != "EHE" && channel != "EHN" {
		err := fmt.Errorf("no channel selected")
		return "", nil, err
	}

	// Re-sort Geophone struct array
	sort.Slice(data, func(i, j int) bool {
		return data[i].TS < data[j].TS
	})

	// Get basic info
	var (
		startTime = time.UnixMilli(data[0].TS).UTC()
		endTime   = time.UnixMilli(data[len(data)-1].TS).UTC()
		station   = getStation(options.FeatureOptions.Config)
		network   = getNetwork(options.FeatureOptions.Config)
		location  = getLocation(options.FeatureOptions.Config)
	)

	// Get sample rate
	sampleRate, err := getSampleRate(data, channel)
	if err != nil {
		return "", nil, err
	}

	// Create channel data buffer
	var channelBuffer []int32
	for _, v := range data {
		switch channel {
		case "EHZ":
			channelBuffer = append(channelBuffer, v.EHZ...)
		case "EHE":
			channelBuffer = append(channelBuffer, v.EHE...)
		case "EHN":
			channelBuffer = append(channelBuffer, v.EHN...)
		}
	}

	// Init SAC library
	var sac sacio.SACData
	sac.Init()

	// Set SAC header and body
	sac.SetTime(startTime, duration.Difference(startTime, endTime))
	sac.SetInfo(network, station, location, channel)
	sac.SetBody(int32ToFloat32(channelBuffer), sampleRate)

	// Get SAC file bytes
	sacBytes, err := sac.Encode(sacio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

	// Return filename and bytes (e.g. 2023.193.14.22.51.0317.AS.SHAKE.00.EHZ.D.sac)
	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.sac",
		startTime.Format("2006"),
		startTime.Format("002"),
		startTime.Format("15"),
		startTime.Format("04"),
		startTime.Format("05"),
		// Get the current millisecond
		startTime.Nanosecond()/1000000,
		station, network,
		location, channel,
	)
	return filename, sacBytes, nil
}

func int32ToFloat32(arr []int32) []float32 {
	floatSlice := make([]float32, len(arr))
	for i, num := range arr {
		floatSlice[i] = float32(num)
	}
	return floatSlice
}
