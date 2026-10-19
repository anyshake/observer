package seisevent

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"github.com/sbabiv/xml2map"
	"golang.org/x/sync/singleflight"
)

const HKO_ID = "hko"

type HKO struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (h *HKO) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      HKO_ID,
		Country: "HK",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Hong Kong Observatory",
			"zh-TW": "天文台全球地震資訊網",
		},
	}
}

func (h *HKO) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if h.cache.Valid() {
		baseEvents = h.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &h.sf, HKO_ID, func(ctx context.Context) ([]Event, error) {
			if h.cache.Valid() {
				return h.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://www.hko.gov.hk/gts/QEM/eq_app-30d_uc.xml",
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse HKO XML response
			dataMap, err := xml2map.NewDecoder(strings.NewReader(string(res))).Decode()
			if err != nil {
				return nil, err
			}
			dataMapEvents, ok := dataMap["Earthquake"].(map[string]any)["EventGroup"].(map[string]any)["Event"]
			if !ok {
				return nil, errors.New("source data is not valid, missing Earthquake.EventGroup.Event")
			}

			// Ensure the response has the expected keys
			expectedKeys := []string{"EventId", "Verify", "HKTDate", "HKTTime", "City", "Region", "Lat", "Lon", "Mag", "Depth"}

			var resultArr []Event
			for _, v := range dataMapEvents.([]map[string]any) {
				if !isMapHasKeys(v, expectedKeys) {
					continue
				}

				timestamp, err := h.getTimestamp(fmt.Sprintf("%s %s00", v["HKTDate"].(string), v["HKTTime"].(string)))
				if err != nil {
					continue
				}

				resultArr = append(resultArr, Event{
					Depth:     string2Float(v["Depth"].(string)),
					Timestamp: timestamp,
					Event:     v["EventId"].(string),
					Verfied:   v["Verify"].(string) == "Y",
					Region:    fmt.Sprintf("%s - %s", v["City"].(string), v["Region"].(string)),
					Latitude:  string2Float(v["Lat"].(string)),
					Longitude: string2Float(v["Lon"].(string)),
					Magnitude: h.getMagnitude(v["Mag"].(string)),
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			h.cache.Set(sortedEvents)
			return sortedEvents, nil
		})
		if err != nil {
			return nil, err
		}

		baseEvents = v
	}

	baseEvents = slices.Clone(baseEvents)
	for i := range baseEvents {
		baseEvents[i].Distance = getDistance(latitude, baseEvents[i].Latitude, longitude, baseEvents[i].Longitude)
		baseEvents[i].Estimation = getSeismicEstimation(
			h.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (h *HKO) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("20060102 150405", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-8 * time.Hour).UnixMilli(), nil
}

func (h *HKO) getMagnitude(data string) []Magnitude {
	return []Magnitude{
		{Type: ParseMagnitude("M"), Value: string2Float(data)},
	}
}
