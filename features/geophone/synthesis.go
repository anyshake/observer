package geophone

import "math"

func GetSynthesis(vertical, eastWest, northSouth [100]float64) [100]float64 {
	var synthesis [100]float64
	for i := range synthesis {
		value := math.Pow(vertical[i], 2) + math.Pow(eastWest[i], 2) + math.Pow(northSouth[i], 2)
		synthesis[i] = math.Sqrt(value)
	}

	return synthesis
}
