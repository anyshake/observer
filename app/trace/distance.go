package trace

import "math"

func GetDistance(lat1, lat2, lng1, lng2 float64) float64 {
	var (
		radius = 6371000.0
		rad    = math.Pi / 180.0
	)

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(
		math.Sin(lat1)*math.Sin(lat2) +
			math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta),
	)

	return math.Abs(dist * radius / 1000)
}
