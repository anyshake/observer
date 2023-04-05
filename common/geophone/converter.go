package geophone

import (
	"strconv"
)

func GetAcceleration(voltage int32, sensitivity float64) float64 {
	result := float64(voltage) / sensitivity
	formatted := strconv.FormatFloat(result, 'f', 5, 64)
	f, err := strconv.ParseFloat(formatted, 64)
	if err != nil {
		return 0.0
	}

	return f
}
