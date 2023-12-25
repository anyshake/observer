package trace

import "math"

func getDistance(lat1, lat2, lng1, lng2 float64) float64 {
	var (
		radius = 6378.137
		rad    = math.Pi / 180.0
	)
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	var (
		a      = lat1 - lat2
		b      = lng1 - lng2
		cal    = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(b/2), 2))) * radius
		result = math.Round(cal*10000) / 10000
	)
	return result
}
