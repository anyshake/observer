package seisevent

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const CENC_APP_ID = "cenc_app"

type CENC_APP struct {
	cache cache.BytesCache
}

func (c *CENC_APP) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CENC_APP_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "China Earthquake Networks Center (APP)",
			"zh-TW": "中國地震台網中心（APP 端）",
			"zh-CN": "中国地震台网中心（APP 端）",
		},
	}
}

func (c *CENC_APP) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.POST(
			"http://api.dizhensubao.igexin.com/api.htm",
			`{"action":"requestMonitorDataAction","startTime":"0","dataSource":"CEIC"}`,
			"application/json", 10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse CENC JSON response
	var dataMap map[string]any
	err := json.Unmarshal(c.cache.Get(), &dataMap)
	if err != nil {
		return nil, err
	}

	// Check server response
	if dataMap["result"].(string) != "OK" {
		return nil, errors.New("remote server returned an error")
	}

	dataMapEvents, ok := dataMap["values"]
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	// Ensure the response has the expected keys and they are not empty
	expectedKeys := []string{"time", "longitude", "latitude", "depth", "eqid", "loc_name", "mag", "eq_type"}

	var resultArr []Event
	for _, v := range dataMapEvents.([]any) {
		event := v.(map[string]any)

		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		region := event["loc_name"].(string)
		if strings.HasPrefix(region, "中国") {
			region = strings.ReplaceAll(region, "中国", "")
		}
		if strings.HasPrefix(region, "台湾省") {
			region = strings.ReplaceAll(region, "台湾", "")
		}

		seisEvent := Event{
			Region:    region,
			Verfied:   event["eq_type"].(string) == "M",
			Depth:     event["depth"].(float64) / 1000,
			Event:     event["eqid"].(string),
			Latitude:  event["latitude"].(float64),
			Longitude: event["longitude"].(float64),
			Magnitude: c.getMagnitude(event["mag"].(float64)),
			Timestamp: time.UnixMilli(int64(event["time"].(float64))).UnixMilli(),
		}

		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (c *CENC_APP) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}

}
