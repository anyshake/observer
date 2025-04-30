package seisevent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const KNDC_ID = "kndc"

type KNDC struct {
	cache cache.AnyCache
}

func (c *KNDC) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      KNDC_ID,
		Country: "KZ",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Kazakhstan National Data Center",
			"zh-TW": "哈薩克國家資料中心",
			"zh-CN": "哈萨克国家数据中心",
		},
	}
}

func (c *KNDC) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://kndc.kz/kndc/pagecontent/alarm-bulletin/getOriginList.php?orderby=epochtime&desc=yes&limit=100",
		30*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	var dataMapEvents []map[string]any
	err = json.Unmarshal(res, &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"id", "epochtime", "lat", "lon", "depth", "mb", "gregion", "sregion"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		timestamp := int64(string2Float(event["epochtime"].(string))) * 1000
		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Event:     event["id"].(string),
			Depth:     event["depth"].(float64),
			Latitude:  string2Float(event["lat"].(string)),
			Longitude: string2Float(event["lon"].(string)),
			Magnitude: c.getMagnitude(event["mb"].(float64)),
			Region:    fmt.Sprintf("%s (%s)", event["gregion"].(string), event["sregion"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *KNDC) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("Mb"), Value: data}}
}
