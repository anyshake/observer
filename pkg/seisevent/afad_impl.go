package seisevent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const AFAD_ID = "afad"

type AFAD struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
	timeSource      *timesource.Source
}

func (c *AFAD) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      AFAD_ID,
		Country: "TR",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Ministry of Disaster and Emergency Management",
			"zh-TW": "土耳其災害與應變管理署",
		},
	}
}

func (c *AFAD) getRequestParam(currentTime time.Time) string {
	startTime := url.QueryEscape(currentTime.AddDate(0, 0, -5).UTC().Format("2006-01-02 15:04:05"))
	endTime := url.QueryEscape(currentTime.UTC().Format("2006-01-02 15:04:05"))
	return fmt.Sprintf("start=%s&end=%s&orderby=timedesc", startTime, endTime)
}

func (c *AFAD) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &c.sf, AFAD_ID, func(ctx context.Context) ([]Event, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			// Make AFAD API request
			res, err := request.GET(
				ctx, fmt.Sprintf("https://deprem.afad.gov.tr/apiv2/event/filter?%s", c.getRequestParam(c.timeSource.Now())),
				30*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse AFAD JSON response
			var dataMapEvents []map[string]any
			err = json.Unmarshal(res, &dataMapEvents)
			if err != nil {
				return nil, err
			}

			// Ensure the response has the expected keys and values
			expectedKeys := []string{"eventID", "location", "latitude", "longitude", "depth", "type", "magnitude", "date"}

			var resultArr []Event
			for _, event := range dataMapEvents {
				if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
					continue
				}

				timestamp, err := c.getTimestamp(event["date"].(string))
				if err != nil {
					return nil, err
				}

				resultArr = append(resultArr, Event{
					Verfied:   true,
					Timestamp: timestamp,
					Event:     event["eventID"].(string),
					Region:    event["location"].(string),
					Depth:     string2Float(event["depth"].(string)),
					Latitude:  string2Float(event["latitude"].(string)),
					Longitude: string2Float(event["longitude"].(string)),
					Magnitude: c.getMagnitudeType(event["type"].(string), event["magnitude"].(string)),
				})
			}

			sorted := sortSeismicEvents(resultArr)
			c.cache.Set(sorted)
			return sorted, nil
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

func (c *AFAD) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.UnixMilli(), nil
}

func (c *AFAD) getMagnitudeType(magType, magText string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude(magType), Value: string2Float(magText)}}
}
