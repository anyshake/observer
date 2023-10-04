package history

import (
	"fmt"
	"sort"
	"time"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/publisher"
	"github.com/bclswl0827/observer/utils/duration"
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
	sacBytes, err := sac.GetBytes(sacio.MSBFIRST)
	if err != nil {
		return "", nil, err
	}

	// Return filename and bytes
	filename := fmt.Sprintf("%s_%s_%s.sac", station, channel, network)
	return filename, sacBytes, nil
}

func int32ToFloat32(arr []int32) []float32 {
	floatSlice := make([]float32, len(arr))
	for i, num := range arr {
		floatSlice[i] = float32(num)
	}
	return floatSlice
}
