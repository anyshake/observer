package seisevent

import (
	"embed"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/bclswl0827/travel"
	"github.com/samber/lo"
)

//go:embed geojson/*.geojson
var geojsonData embed.FS

func getGeoJsonData(name string) ([]byte, error) {
	return geojsonData.ReadFile(fmt.Sprintf("geojson/%s.geojson", name))
}

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

func isMapHasKeys[T any](m map[string]T, keys []string) bool {
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
	const (
		radius = 6378.137
		rad    = math.Pi / 180.0
	)
	lat1 *= rad
	lng1 *= rad
	lat2 *= rad
	lng2 *= rad

	sinLat := math.Sin((lat1 - lat2) / 2)
	sinLng := math.Sin((lng1 - lng2) / 2)
	h := sinLat*sinLat + math.Cos(lat1)*math.Cos(lat2)*sinLng*sinLng
	cal := 2 * math.Asin(math.Sqrt(h)) * radius
	return math.Round(cal*10000) / 10000
}

func getSeismicEstimation(table *travel.AK135, lat1, lat2, lng1, lng2, depth float64) Estimation {
	result := table.Estimate(travel.GetDeltaByCoordinates(lat1, lng1, lat2, lng2), lo.Ternary(depth > 0, depth, 0), true)
	estObj := Estimation{P_Wave: -1, S_Wave: -1}

	if result.P != nil {
		estObj.P_Wave = result.P.Duration.Seconds()
	} else if result.PKIKP != nil {
		estObj.P_Wave = result.PKIKP.Duration.Seconds()
	}

	if result.S != nil {
		estObj.S_Wave = result.S.Duration.Seconds()
	} else if result.SKIKS != nil {
		estObj.S_Wave = result.SKIKS.Duration.Seconds()
	}

	return estObj
}
