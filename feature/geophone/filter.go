package geophone

import (
	"github.com/bclswl0827/observer/config"
)

func (g *Geophone) getFilter(geophoneAttr config.Compensation, sampleRate int) filter {
	// Fallback filter coefficients for most geophones
	fallbackFilter := filter{
		a1: 1.99823115,
		a2: -0.99822469,
		b0: 1.03380975,
		b1: -1.99662644,
		b2: 0.96601161,
	}

	// // Get compensation parameters
	// var (
	// 	xi_1    = geophoneAttr.Damping
	// 	xi_c    = TARGET_DAMPING
	// 	omega_0 = 2 * math.Pi * geophoneAttr.Frequency
	// 	omega_c = 2 * math.Pi * TARGET_FREQUENCY
	// )

	// // Generate compensation network transfer function
	// compensation_transfer_func := func(s float64) float64 {
	// 	original := math.Pow(s, 2) + 2*omega_0*xi_1*s + math.Pow(omega_0, 2)
	// 	target := math.Pow(s, 2) + 2*omega_c*xi_c*s + math.Pow(omega_c, 2)
	// 	return original / target
	// }
	// // Generate discrete transfer function by bilinear transform
	// discrete_transfer_func := func(z float64, samplerate int) float64 {
	// 	reciprocal_z := 1 / z
	// 	t := 1 / float64(samplerate) / 2
	// 	s := (2 / t) * ((1 - reciprocal_z) / (1 + reciprocal_z))
	// 	return compensation_transfer_func(s)
	// }

	return fallbackFilter
}

func (g *Geophone) applyFilter(chValues []int32, filter filter) []int32 {
	x := make([]float64, len(chValues))
	for i, v := range chValues {
		x[i] = float64(v)
	}

	y := []float64{float64(chValues[0]), float64(chValues[1])}
	for i := 2; i < len(chValues); i++ {
		yNext := filter.a1*y[i-1] + filter.a2*y[i-2] + filter.b0*x[i] + filter.b1*x[i-1] + filter.b2*x[i-2]
		y = append(y, yNext)
	}

	yValues := make([]int32, len(y))
	for i, v := range y {
		yValues[i] = int32(v)
	}

	return yValues
}
