package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const ICL_ID = "icl"

type ICL struct {
	cache cache.BytesCache
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
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://mobile-new.chinaeew.cn/v1/earlywarnings?start_at=&updates=",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse ICL JSON response
	var dataMap map[string]any
	err := json.Unmarshal(c.cache.Get(), &dataMap)
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
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (c *ICL) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
