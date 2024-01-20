package resample

import "math"

// Resample resamples the data from the original rate to the target rate using linear interpolation.
func Resample(data []int32, originalRate, targetRate int) []int32 {
	scale := float64(originalRate) / float64(targetRate)
	newDataSize := int(math.Floor(float64(len(data)) / scale))
	newData := make([]int32, newDataSize)

	for i := 0; i < newDataSize; i++ {
		// Calculate the position in the original data
		pos := scale * float64(i)
		base := int(math.Floor(pos))
		frac := pos - float64(base)

		// Handle the boundary
		if base >= len(data)-1 {
			newData[i] = data[len(data)-1]
		} else {
			// Linear interpolation
			interpolated := Interpolate(float64(data[base]), float64(data[base+1]), frac)
			newData[i] = int32(math.Round(interpolated))
		}
	}

	return newData
}
