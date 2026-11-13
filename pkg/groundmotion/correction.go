package groundmotion

import (
	"math"
	"math/cmplx"
	"sort"
)

const (
	// CORRECTION_LIMIT selects the sensor poles that are inverted: only poles
	// with a natural frequency below this value (Hz) are treated as the
	// low-frequency roll-off of a geophone. Higher poles (anti-alias filters)
	// are left untouched because inverting them would amplify noise.
	CORRECTION_LIMIT = 10.0
	// MAX_BANDWIDTH_EXTENSION limits how far below the sensor's natural
	// frequency the response is extended, which bounds noise amplification.
	MAX_BANDWIDTH_EXTENSION = 10.0
)

// section is a second-order IIR section (a0 = 1) in transposed direct form II.
type section struct {
	b0, b1, b2, a1, a2 float64
	z1, z2             float64
}

func (s *section) apply(x float64) float64 {
	y := s.b0*x + s.z1
	s.z1 = s.b1*x - s.a1*y + s.z2
	s.z2 = s.b2*x - s.a2*y
	return y
}

// correction removes the low-frequency poles of an analog (Laplace) sensor
// response and replaces them with a Butterworth high-pass of the same order:
//
//	C(s) = Π(s - p) / Π(s - q)
//
// For a velocity transducer T(s) = s^n / Π(s - p), the corrected response
// T(s)·C(s) = s^n / Π(s - q) is a flat pass band down to the new corner.
// It is the streaming, O(1)-per-sample equivalent of a frequency-domain
// response removal restricted to the band that matters for PGV/PGD.
type correction struct {
	sections []section
}

// lowFrequencyPoles returns the poles (rad/s) below CORRECTION_LIMIT.
func lowFrequencyPoles(poles []complex128) []complex128 {
	var selected []complex128
	for _, pole := range poles {
		if cmplx.IsNaN(pole) || cmplx.IsInf(pole) || real(pole) >= 0 {
			continue // unstable or invalid poles cannot be inverted safely
		}
		if cmplx.Abs(pole)/(2*math.Pi) < CORRECTION_LIMIT {
			selected = append(selected, pole)
		}
	}
	return selected
}

// newCorrection designs the correction filter. It returns nil when no pole
// needs to be corrected or the design is not representable at this rate.
// When frequency > 0 the filter is normalized to unit magnitude there, so the
// overall sensitivity (defined at that frequency) remains valid.
func newCorrection(poles []complex128, frequency float64, rate int) *correction {
	low := lowFrequencyPoles(poles)
	if len(low) == 0 || rate < 1 {
		return nil
	}
	highest := 0.0
	for _, pole := range low {
		highest = max(highest, cmplx.Abs(pole)/(2*math.Pi))
	}
	corner := max(HIGHPASS_CORNER, highest/MAX_BANDWIDTH_EXTENSION)
	if corner >= float64(rate)/4 {
		return nil
	}
	replacement := butterworthPoles(len(low), 2*math.Pi*corner)

	// Bilinear transform: s - r  ->  (2/T - r)(z - zd)/(z + 1),
	// zd = (2/T + r)/(2/T - r). The (z + 1) terms cancel since the orders match.
	k := 2 * float64(rate)
	gain := complex(1, 0)
	zeros := make([]complex128, len(low))
	for i, pole := range low {
		gain *= complex(k, 0) - pole
		zeros[i] = (complex(k, 0) + pole) / (complex(k, 0) - pole)
	}
	digital := make([]complex128, len(replacement))
	for i, pole := range replacement {
		gain /= complex(k, 0) - pole
		digital[i] = (complex(k, 0) + pole) / (complex(k, 0) - pole)
	}

	numerators := quadraticFactors(zeros)
	denominators := quadraticFactors(digital)
	if len(numerators) != len(denominators) {
		return nil
	}
	result := &correction{sections: make([]section, len(numerators))}
	for i := range numerators {
		b, a := numerators[i], denominators[i]
		result.sections[i] = section{b0: b[0], b1: b[1], b2: b[2], a1: a[1], a2: a[2]}
	}
	g := real(gain)
	if frequency > 0 && frequency < float64(rate)/2 {
		s := complex(0, 2*math.Pi*frequency)
		magnitude := 1.0
		for i := range low {
			magnitude *= cmplx.Abs(s-low[i]) / cmplx.Abs(s-replacement[i])
		}
		if magnitude > 0 && !math.IsInf(magnitude, 0) {
			g /= magnitude
		}
	}
	result.sections[0].b0 *= g
	result.sections[0].b1 *= g
	result.sections[0].b2 *= g
	return result
}

func (c *correction) apply(value float64) float64 {
	for i := range c.sections {
		value = c.sections[i].apply(value)
	}
	return value
}

// butterworthPoles returns the n analog Butterworth poles at corner (rad/s).
func butterworthPoles(n int, corner float64) []complex128 {
	poles := make([]complex128, n)
	for k := range n {
		theta := math.Pi * float64(2*k+n+1) / float64(2*n)
		poles[k] = cmplx.Rect(corner, theta)
	}
	return poles
}

// quadraticFactors groups roots into real polynomials [1, c1, c2] in z^-1:
// complex roots with their conjugates, real roots in pairs, and a remaining
// real root as a first-order factor. Two root sets of equal size with the
// same number of real roots modulo two yield the same factor layout.
func quadraticFactors(roots []complex128) [][3]float64 {
	const tolerance = 1e-9
	var reals []float64
	var complexes []complex128
	for _, root := range roots {
		if math.Abs(imag(root)) <= tolerance*max(1, cmplx.Abs(root)) {
			reals = append(reals, real(root))
		} else if imag(root) > 0 {
			complexes = append(complexes, root)
		}
	}
	sort.Float64s(reals)
	var factors [][3]float64
	for _, root := range complexes {
		factors = append(factors, [3]float64{1, -2 * real(root), real(root)*real(root) + imag(root)*imag(root)})
	}
	for i := 0; i+1 < len(reals); i += 2 {
		factors = append(factors, [3]float64{1, -(reals[i] + reals[i+1]), reals[i] * reals[i+1]})
	}
	if len(reals)%2 == 1 {
		factors = append(factors, [3]float64{1, -reals[len(reals)-1], 0})
	}
	return factors
}
