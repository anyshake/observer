package seisevent

import (
	"math"
	"sort"
	"strconv"

	"github.com/bclswl0827/travel"
)

func string2Float(num string) float64 {
	r, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0.0
	}

	return r
}

func isMapKeysEmpty(m map[string]any, keys []string) bool {
	for _, key := range keys {
		switch m[key].(type) {
		case string:
			if len(m[key].(string)) == 0 {
				return false
			}
		default:
			continue
		}
	}

	return true
}

func isMapHasKeys(m map[string]any, keys []string) bool {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}

	return true
}

func sortSeismicEvents(events []Event) []Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp > events[j].Timestamp
	})

	return events
}

func getDistance(lat1, lat2, lng1, lng2 float64) float64 {
	var (
		radius = 6378.137
		rad    = math.Pi / 180.0
	)
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	var (
		a      = lat1 - lat2
		b      = lng1 - lng2
		cal    = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(b/2), 2))) * radius
		result = math.Round(cal*10000) / 10000
	)
	return result
}

func getSeismicEstimation(table *travel.AK135, lat1, lat2, lng1, lng2, depth float64) Estimation {
	result := table.Estimate(travel.GetDeltaByCoordinates(lat1, lng1, lat2, lng2), depth, true)
	estObj := Estimation{P_Wave: -1, S_Wave: -1}

	if result.P != nil {
		estObj.P_Wave = result.P.Duration.Seconds()
	} else if result.PKPdf != nil {
		estObj.P_Wave = result.PKPdf.Duration.Seconds()
	}

	if result.S != nil {
		estObj.S_Wave = result.S.Duration.Seconds()
	} else if result.SKSdf != nil {
		estObj.S_Wave = result.SKSdf.Duration.Seconds()
	}

	return estObj
}
