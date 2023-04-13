package geophone

import (
	"fmt"
	"math"
	"strconv"
)

func GetSynthesis(vertical, eastWest, northSouth float64) float64 {
	value := math.Pow(vertical, 2) + math.Pow(eastWest, 2) + math.Pow(northSouth, 2)
	synthesis, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", math.Sqrt(value)), 64)

	return synthesis
}
