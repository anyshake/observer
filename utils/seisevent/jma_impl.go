package seisevent

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const JMA_ID = "jma"

type JMA struct {
	cache cache.BytesCache
}

func (j *JMA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      JMA_ID,
		Country: "JP",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Japan Meteorological Agency",
			"zh-TW": "氣象廳",
			"zh-CN": "气象厅",
		},
	}
}

func (j *JMA) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !j.cache.Valid() {
		res, err := request.GET(
			"https://www.jma.go.jp/bosai/quake/data/list.json",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		j.cache.Set(res)
	}

	// Parse JMA JSON response
	var dataMapEvents []map[string]any
	err := json.Unmarshal(j.cache.Get(), &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"eid", "anm", "mag", "cod", "at"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		timestamp, err := j.getTimestamp(event["at"].(string))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Depth:     j.getDepth(event["cod"].(string)),
			Event:     event["eid"].(string),
			Region:    event["anm"].(string),
			Latitude:  j.getLatitude(event["cod"].(string)),
			Longitude: j.getLongitude(event["cod"].(string)),
			Magnitude: j.getMagnitude(event["mag"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (j *JMA) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05+09:00", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-9 * time.Hour).UnixMilli(), nil
}

func (j *JMA) getDepth(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return string2Float(arr[2]) / 1000
}

func (j *JMA) getLatitude(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return string2Float(arr[0])
}

func (j *JMA) getLongitude(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return string2Float(arr[1])
}

func (j *JMA) getMagnitude(data string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: string2Float(data)}}
}
