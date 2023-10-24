package trace

func getEstimation(distance float64) estimation {
	pWave := distance / 6.0
	sWave := distance / 3.5
	return estimation{
		P: pWave,
		S: sWave,
	}
}
