package geophone

import (
	"github.com/shopspring/decimal"
)

func GetAcceleration(voltage int32, sensitivity float64) float64 {
	s := decimal.NewFromFloat(sensitivity)
	v := decimal.NewFromInt32(voltage)
	r, _ := v.Div(s).Float64()

	result, _ := decimal.NewFromFloat(r).Round(5).Float64()
	return result
}
