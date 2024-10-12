package seisevent

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const JMA_WOLFX_ID = "jma_wolfx"

type JMA_WOLFX struct {
	cache cache.AnyCache
}

func (j *JMA_WOLFX) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      JMA_WOLFX_ID,
		Country: "JP",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Japan Meteorological Agency (Wolfx)",
			"zh-TW": "氣象廳（Wolfx）",
			"zh-CN": "气象厅（Wolfx）",
		},
	}
}

func (j *JMA_WOLFX) GetEvents(latitude, longitude float64) ([]Event, error) {
	if j.cache.Valid() {
		return j.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://api.wolfx.jp/jma_eqlist.json",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse JMA_WOLFX JSON response
	var dataMapEvents map[string]any
	err = json.Unmarshal(res, &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"EventID", "time_full", "location", "magnitude", "depth", "latitude", "longitude"}

	var resultArr []Event
	for _, v := range dataMapEvents {
		event, ok := v.(map[string]any)
		if !ok {
			continue
		}

		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		timestamp, err := j.getTimestamp(event["time_full"].(string))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Depth:     j.getDepth(event["depth"].(string)),
			Event:     event["EventID"].(string),
			Region:    event["location"].(string),
			Latitude:  string2Float(event["latitude"].(string)),
			Longitude: string2Float(event["longitude"].(string)),
			Magnitude: j.getMagnitude(event["magnitude"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	j.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (j *JMA_WOLFX) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-9 * time.Hour).UnixMilli(), nil
}

func (j *JMA_WOLFX) getDepth(data string) float64 {
	depthStr := strings.ReplaceAll(data, "km", "")
	return string2Float(depthStr)
}

func (j *JMA_WOLFX) getMagnitude(data string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: string2Float(data)}}
}
