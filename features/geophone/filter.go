package geophone

import (
	"github.com/mattetti/audio/dsp/filters"
	"github.com/mattetti/audio/dsp/windows"
)

func LowPassFilter(signal []float64, cutoff float64, taps int) ([]float64, error) {
	fir := &filters.FIR{Sinc: &filters.Sinc{
		Taps:         taps,
		CutOffFreq:   cutoff,
		SamplingFreq: PACKET_SIZE,
		Window:       windows.Hamming,
	}}

	filtered, err := fir.LowPass(signal)
	if err != nil {
		return nil, err
	}

	return filtered, nil
}
