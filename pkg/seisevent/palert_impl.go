package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const PALERT_ID = "p-alert"

type PALERT struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *PALERT) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      PALERT_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "P-Alert Strong Motion Network",
			"zh-TW": "P-Alert 觀測網",
			"zh-CN": "P-Alert 观测网",
		},
	}
}

func (c *PALERT) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	// To build the query for P-Alert, we need to get the current timestamp
	res, err := request.GET(
		"https://www.cloudflare.com/cdn-cgi/trace",
		10*time.Second, time.Second, 3, false,
		nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}
	ts := regexp.MustCompile(`ts=(\d+)`).FindStringSubmatch(string(res))
	if len(ts) == 0 {
		return nil, errors.New("failed to get current time from Cloudflare")
	}
	currentTime := time.Unix(int64(string2Float(ts[1])), 0)

	res, err = request.POST(
		"https://palert.earth.sinica.edu.tw/graphql/",
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

		seisEvent.Region = c.getRegion(seisEvent.Latitude, seisEvent.Longitude)
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *PALERT) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("ML"), Value: data}}
}

func (c *PALERT) getRegion(latitude, longitude float64) string {
	return fmt.Sprintf("Latitude: %.3f°, Longitude: %.3f°", latitude, longitude)
}

func (c *PALERT) getTimestamp(textValue string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05", textValue)
	return t.UnixMilli(), err
}
