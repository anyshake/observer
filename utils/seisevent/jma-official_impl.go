package seisevent

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const JMA_OFFICIAL_ID = "jma_official"

type JMA_OFFICIAL struct {
	cache cache.AnyCache
}

func (j *JMA_OFFICIAL) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      JMA_OFFICIAL_ID,
		Country: "JP",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Japan Meteorological Agency (Official)",
			"zh-TW": "氣象廳（官方）",
			"zh-CN": "气象厅（官方）",
		},
	}
}

func (j *JMA_OFFICIAL) GetEvents(latitude, longitude float64) ([]Event, error) {
	if j.cache.Valid() {
		return j.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://www.jma.go.jp/bosai/quake/data/list.json",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse JMA_OFFICIAL JSON response
	var dataMapEvents []map[string]any
	err = json.Unmarshal(res, &dataMapEvents)
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

		// Second is the last 2 digits of the event ID
		eventId := event["eid"].(string)
		timestamp, err := j.getTimestamp(event["at"].(string), eventId[len(eventId)-2:])
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Depth:     j.getDepth(event["cod"].(string)),
			Event:     eventId,
			Region:    event["anm"].(string),
			Latitude:  j.getLatitude(event["cod"].(string)),
			Longitude: j.getLongitude(event["cod"].(string)),
			Magnitude: j.getMagnitude(event["mag"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	j.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (j *JMA_OFFICIAL) getTimestamp(timeStr, secStr string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05+09:00", timeStr)
	if err != nil {
		return 0, err
	}

	sec := int(string2Float(secStr))
	return t.Add(time.Duration(sec) * time.Second).Add(-9 * time.Hour).UnixMilli(), nil
}

func (j *JMA_OFFICIAL) getDepth(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return string2Float(arr[2]) / 1000
}

func (j *JMA_OFFICIAL) getLatitude(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return string2Float(arr[0])
}

func (j *JMA_OFFICIAL) getLongitude(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return string2Float(arr[1])
}

func (j *JMA_OFFICIAL) getMagnitude(data string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: string2Float(data)}}
}
