package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const ICL_ID = "icl"

type ICL struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *ICL) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      ICL_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Institute of Care-life (ICL)",
			"zh-TW": "成都高新減災研究所",
			"zh-CN": "成都高新减灾研究所",
		},
	}
}

func (c *ICL) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://mobile-new.chinaeew.cn/v1/earlywarnings?start_at=&updates=",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse ICL JSON response
	var dataMap map[string]any
	err = json.Unmarshal(res, &dataMap)
	if err != nil {
		return nil, err
	}

	// Check server response
	if dataMap["code"].(float64) != 0 {
		return nil, errors.New("remote server returned an error")
	}

	dataMapEvents, ok := dataMap["data"]
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	// Ensure the response has the expected keys and they are not empty
	expectedKeys := []string{"eventId", "latitude", "longitude", "depth", "epicenter", "startAt", "magnitude"}

	var resultArr []Event
	for _, v := range dataMapEvents.([]any) {
		event := v.(map[string]any)

		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Depth:     event["depth"].(float64),
			Region:    event["epicenter"].(string),
			Latitude:  event["latitude"].(float64),
			Longitude: event["longitude"].(float64),
			Event:     fmt.Sprintf("%v", event["eventId"]),
			Magnitude: c.getMagnitude(event["magnitude"].(float64)),
			Timestamp: int64(event["startAt"].(float64)),
		}

		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *ICL) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
