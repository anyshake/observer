package trace

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/request"
)

type JMA struct {
	dataSourceCache
}

func (j *JMA) Property() string {
	return "気象庁地震情報"
}

func (j *JMA) Fetch() ([]byte, error) {
	if time.Since(j.Time) <= EXPIRATION {
		return j.Cache, nil
	}

	res, err := request.GET(
		"https://www.jma.go.jp/bosai/quake/data/list.json",
		10*time.Second, time.Second, 3, false, nil,
	)
	if err != nil {
		return nil, err
	}

	j.Time = time.Now()
	j.Cache = make([]byte, len(res))
	copy(j.Cache, res)

	return res, nil
}

func (j *JMA) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any, 0)
	result["data"] = make([]any, 0)

	arr := make([]map[string]any, 0)
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return nil, err
	}

	result["data"] = arr
	return result, nil
}

func (j *JMA) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	keys := []string{"anm", "mag", "cod", "at"}

	var list []seismicEvent
	for _, v := range data["data"].([]map[string]any) {
		if !isMapHasKeys(v, keys) || !isMapKeysEmpty(v, keys) {
			continue
		}

		ts, err := time.Parse("2006-01-02T15:04:05+09:00", v["at"].(string))
		if err != nil {
			continue
		}

		l := seismicEvent{
			Depth:     j.getDepth(v["cod"].(string)),
			Verfied:   true,
			Timestamp: ts.Add(-9 * time.Hour).UnixMilli(),
			Event:     v["anm"].(string),
			Region:    v["anm"].(string),
			Latitude:  j.getLatitude(v["cod"].(string)),
			Longitude: j.getLongitude(v["cod"].(string)),
			Magnitude: string2Float(v["mag"].(string)),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getSeismicEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (j *JMA) List(latitude, longitude float64) ([]seismicEvent, error) {
	res, err := j.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := j.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := j.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
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
