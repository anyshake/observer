package geophone

import "sort"

func (g *Geophone) getOffsetCounts(data []int32) []int32 {
	sortedData := make([]int32, len(data))
	copy(sortedData, data)
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i] < sortedData[j]
	})

	var median int32
	if len(sortedData)%2 == 1 {
		middleIndex := len(sortedData) / 2
		median = sortedData[middleIndex]
	} else {
		middleIndex1 := len(sortedData)/2 - 1
		middleIndex2 := len(sortedData) / 2
		median = (sortedData[middleIndex1] + sortedData[middleIndex2]) / 2
	}

	result := make([]int32, len(data))
	for i, value := range data {
		result[i] = value - median
	}

	return result
}
