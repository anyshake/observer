package export

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/pkg/seekbuf"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

const FILTER_NUM_TAPS = 101

type seismicDataEncoderWavImpl struct {
	outputSampleRate   int
	actionHandler      *action.Handler
	stationCodeConfig  config.StationStationCodeConfigConstraintImpl
	locationCodeConfig config.StationLocationCodeConfigConstraintImpl
	networkCodeConfig  config.StationNetworkCodeConfigConstraintImpl
}

func (e *seismicDataEncoderWavImpl) GetName() string {
	return "WAV (Audio)"
}

func (e *seismicDataEncoderWavImpl) Encode(records seismicRecordIterator, channelCode string) ([]byte, error) {
	var channelBuffer []int32
	recordRange, err := forEachContinuousChannelRecord(records, channelCode, func(_ model.SeisRecord, samples []int32) error {
		channelBuffer = append(channelBuffer, samples...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if recordRange.Count == 0 {
		return nil, nil
	}

	audioData := e.normalizeToInt16(channelBuffer)

	timeDiff := time.UnixMilli(recordRange.EndTimestamp).Sub(time.UnixMilli(recordRange.StartTimestamp)).Seconds()
	if timeDiff == 0 {
		return nil, errors.New("invalid time difference")
	}

	originalSampleRate := int(float64(len(audioData)) / (timeDiff / 30.0))
	filterKernel := e.getLowPassFilter(originalSampleRate, 40, FILTER_NUM_TAPS)
	audioData = e.applyFilter(audioData, filterKernel)

	interpolatedData := e.linearInterpolate(audioData, originalSampleRate, e.outputSampleRate)

	dataBytes, err := e.saveToWavBytes(interpolatedData, e.outputSampleRate)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (e *seismicDataEncoderWavImpl) GetFileName(startTime time.Time, channelCode string) (string, error) {
	stationCode, err := e.stationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}
	locationCode, err := e.locationCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}
	networkCode, err := e.networkCodeConfig.Get(e.actionHandler)
	if err != nil {
		return "", err
	}

	stationCodeStr := stationCode.(string)
	locationCodeStr := locationCode.(string)
	networkCodeStr := networkCode.(string)

	filename := fmt.Sprintf("%s.%s.%s.%s.%s.%04d.%s.%s.%s.%s.D.wav",
		startTime.UTC().Format("2006"),
		startTime.UTC().Format("002"),
		startTime.UTC().Format("15"),
		startTime.UTC().Format("04"),
		startTime.UTC().Format("05"),
		startTime.UTC().Nanosecond()/1000000,
		stationCodeStr, networkCodeStr,
		locationCodeStr, channelCode,
	)
	return filename, nil
}

func (e *seismicDataEncoderWavImpl) normalizeToInt16(data []int32) []int16 {
	var absMaxVal float64
	for _, value := range data {
		absValue := math.Abs(float64(value))
		if absValue > absMaxVal {
			absMaxVal = absValue
		}
	}
	if absMaxVal == 0 {
		return nil
	}
	scaleFactor := float64(math.MaxInt16) / absMaxVal
	normalized := make([]int16, len(data))
	for i, value := range data {
		normalized[i] = int16(float64(value) * scaleFactor)
	}
	return normalized
}

func (e *seismicDataEncoderWavImpl) linearInterpolate(data []int16, oldRate, newRate int) []int16 {
	if oldRate == newRate || oldRate == 0 {
		return data
	}

	scale := float64(oldRate) / float64(newRate)
	newLength := int(float64(len(data)) * float64(newRate) / float64(oldRate))
	interpolated := make([]int16, newLength)

	for i := 0; i < newLength; i++ {
		origIndex := float64(i) * scale
		index := int(origIndex)
		fraction := origIndex - float64(index)

		if index+1 < len(data) {
			interpolated[i] = int16(float64(data[index])*(1-fraction) + float64(data[index+1])*fraction)
		} else {
			interpolated[i] = data[index]
		}
	}

	return interpolated
}

func (e *seismicDataEncoderWavImpl) getLowPassFilter(sampleRate int, cutoffFreq float64, numTaps int) []float64 {
	normalizedCutoff := cutoffFreq / float64(sampleRate)
	coeffs := make([]float64, numTaps)

	center := numTaps / 2
	for i := 0; i < numTaps; i++ {
		n := float64(i - center)
		if i == center {
			coeffs[i] = 2 * normalizedCutoff
		} else {
			coeffs[i] = math.Sin(2*math.Pi*normalizedCutoff*n) / (math.Pi * n)
		}

		// Hamming window
		coeffs[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(numTaps-1))
	}

	// Normalize gain
	var sum float64
	for _, c := range coeffs {
		sum += c
	}
	for i := range coeffs {
		coeffs[i] /= sum
	}

	return coeffs
}

func (e *seismicDataEncoderWavImpl) applyFilter(data []int16, kernel []float64) []int16 {
	numTaps := len(kernel)
	filtered := make([]int16, len(data))

	for i := range data {
		sum := 0.0
		for j := 0; j < numTaps; j++ {
			if i-j >= 0 {
				sum += float64(data[i-j]) * kernel[j]
			}
		}
		filtered[i] = int16(sum)
	}

	return filtered
}

func (e *seismicDataEncoderWavImpl) saveToWavBytes(data []int16, sampleRate int) ([]byte, error) {
	var buf seekbuf.Buffer
	encoder := wav.NewEncoder(&buf, sampleRate, 16, 1, 1)
	audioSamples := make([]int, len(data))
	for i, value := range data {
		audioSamples[i] = int(value)
	}
	buffer := &audio.IntBuffer{
		Data:   audioSamples,
		Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1},
	}

	if err := encoder.Write(buffer); err != nil {
		return nil, err
	}

	if err := encoder.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
