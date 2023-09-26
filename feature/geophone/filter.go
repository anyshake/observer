package geophone

func (g *Geophone) getFilter(sensitivity, frequency, damping float64) *Filter {
	return nil
}

func (g *Geophone) Filter(xValues []int32, filter *Filter) []int32 {
	x := make([]float64, len(xValues))
	for i, v := range xValues {
		x[i] = float64(v)
	}

	y := []float64{x[0], x[1]}
	for i := 2; i < len(xValues); i++ {
		yNext := filter.a1*y[i-1] + filter.a2*y[i-2] + filter.b0*x[i] + filter.b1*x[i-1] + filter.b2*x[i-2]
		y = append(y, yNext)
	}

	yValues := make([]int32, len(y))
	for i, v := range y {
		yValues[i] = int32(v)
	}

	return yValues
}
