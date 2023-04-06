package geophone

import (
	"github.com/shopspring/decimal"
)

func GetAcceleration(voltage int32, sensitivity float64) float64 {
	s := decimal.NewFromFloat(sensitivity)
	v := decimal.NewFromInt32(voltage)
	result := v.Div(s)

	return result.InexactFloat64()
}
