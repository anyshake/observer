package seisevent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"golang.org/x/sync/singleflight"
)

const PALERT_ID = "p-alert"

type PALERT struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
	timeSource      *timesource.Source
}

func (c *PALERT) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      PALERT_ID,
		Country: "TW",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "P-Alert Strong Motion Network",
			"zh-TW": "P-Alert 觀測網",
		},
	}
}

func (c *PALERT) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &c.sf, PALERT_ID, func(ctx context.Context) ([]Event, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			// source: https://raw.githubusercontent.com/g0v/twgeojson/master/json/twCounty2010merge.topo.json
			twCounty2010, err := getGeoJsonData("twCounty2010")
			if err != nil {
				return nil, err
			}
			featureCollection, err := geojson.UnmarshalFeatureCollection(twCounty2010)
			if err != nil {
				return nil, err
			}

			currentTime := c.timeSource.Now()
			res, err := request.POST(
				ctx, "https://palert.earth.sinica.edu.tw/graphql/",
				fmt.Sprintf(
					`{"query":"query ($date: [Date!], $depth: [Float!], $ml: [Float!], $dateTime: DateTime, $needHaspga: Boolean!) {\n  eventList(\n    QueryEvent: {depth: $depth, date: $date, ml: $ml, dateTime: $dateTime}\n    needHaspga: $needHaspga\n  ) {\n    DateUTC\n    Depth\n    Latitude\n    Longitude\n    ML\n    hasPGA @include(if: $needHaspga)\n  }\n}","variables":{"date":["%s","%s"],"ml":[0,10],"depth":[0,1000],"needHaspga":false}}`,
					currentTime.AddDate(-1, 0, 0).Format("2006-01-02"),
					currentTime.Format("2006-01-02"),
				),
				"application/json", 10*time.Second, time.Second, 3, false, nil,
				map[string]string{
					"User-Agent": uarand.GetRandom(),
					"Referer":    "https://palert.earth.sinica.edu.tw/database",
				},
			)
			if err != nil {
				return nil, err
			}

			// Parse P-Alert JSON response
			var dataMap map[string]any
			err = json.Unmarshal(res, &dataMap)
			if err != nil {
				return nil, err
			}

			dataMapObj, ok := dataMap["data"].(map[string]any)
			if !ok {
				return nil, errors.New("seismic event data object is not available")
			}

			dataMapEvents, ok := dataMapObj["eventList"].([]any)
			if !ok {
				return nil, errors.New("seismic event data is not available")
			}

			// Ensure the response has the expected keys and they are not empty
			expectedKeys := []string{"DateUTC", "Depth", "Latitude", "Longitude", "ML"}

			var resultArr []Event
			for idx, v := range dataMapEvents {
				event := v.(map[string]any)

				if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
					continue
				}

				timestamp, err := c.getTimestamp(event["DateUTC"].(string))
				if err != nil {
					continue
				}

				seisEvent := Event{
					Verfied:   true,
					Timestamp: timestamp,
					Depth:     event["Depth"].(float64),
					Event:     fmt.Sprintf("P-Alert#%d", idx),
					Latitude:  event["Latitude"].(float64),
					Longitude: event["Longitude"].(float64),
					Magnitude: c.getMagnitude(event["ML"].(float64)),
				}
				seisEvent.Region = c.getRegion(featureCollection, seisEvent.Latitude, seisEvent.Longitude)
				resultArr = append(resultArr, seisEvent)
			}

			sortedEvents := sortSeismicEvents(resultArr)
			c.cache.Set(sortedEvents)
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
			c.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (c *PALERT) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("ML"), Value: data}}
}

func (c *PALERT) pointToSegmentDistance(p, a, b orb.Point) float64 {
	px, py := p[0], p[1]
	ax, ay := a[0], a[1]
	bx, by := b[0], b[1]

	dx := bx - ax
	dy := by - ay

	if dx == 0 && dy == 0 {
		return math.Hypot(px-ax, py-ay)
	}

	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)

	if t < 0 {
		return math.Hypot(px-ax, py-ay)
	} else if t > 1 {
		return math.Hypot(px-bx, py-by)
	}

	projX := ax + t*dx
	projY := ay + t*dy

	return math.Hypot(px-projX, py-projY)
}

func (c *PALERT) distanceToPolygon(p orb.Point, poly orb.Polygon) float64 {
	min := math.MaxFloat64

	for _, ring := range poly {
		for i := 0; i < len(ring)-1; i++ {
			d := c.pointToSegmentDistance(p, ring[i], ring[i+1])
			if d < min {
				min = d
			}
		}
	}

	return min
}

func (c *PALERT) getRegion(fc *geojson.FeatureCollection, latitude, longitude float64) string {
	point := orb.Point{longitude, latitude}

	var (
		minDist    = math.MaxFloat64
		nearestReg string
	)

	for _, feature := range fc.Features {
		name := feature.Properties["COUNTYNAME"].(string)

		switch geom := feature.Geometry.(type) {
		case orb.Polygon:
			if planar.PolygonContains(geom, point) {
				return name
			}
			d := c.distanceToPolygon(point, geom)
			if d < minDist {
				minDist = d
				nearestReg = name
			}

		case orb.MultiPolygon:
			if planar.MultiPolygonContains(geom, point) {
				return name
			}
			for _, poly := range geom {
				d := c.distanceToPolygon(point, poly)
				if d < minDist {
					minDist = d
					nearestReg = name
				}
			}
		}
	}

	return fmt.Sprintf("%s（附近海域）", nearestReg)
}

func (c *PALERT) getTimestamp(textValue string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05", textValue)
	return t.UnixMilli(), err
}
