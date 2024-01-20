package resample

// Interpolate performs linear interpolation between two points.
func Interpolate(y1, y2, mu float64) float64 {
	return y1*(1-mu) + y2*mu
}
