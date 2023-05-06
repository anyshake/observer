package trace

import "strconv"

func String2Float(num string) float64 {
	r, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0.0
	}

	return r
}
