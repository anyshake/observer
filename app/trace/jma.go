package trace

import (
	"encoding/json"
	"strings"
	"time"

	"com.geophone.observer/common/request"
)

type JMA struct{}

func (j *JMA) Property() (string, string) {
	const (
		NAME  string = "日本気象庁地震情報"
		VALUE string = "JMA"
	)

	return NAME, VALUE
}

func (j *JMA) Fetch() ([]byte, error) {
	res, err := request.GETRequest(
		"https://www.jma.go.jp/bosai/quake/data/list.json",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (j *JMA) Parse(data []byte) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	result["data"] = make([]interface{}, 0)

	arr := make([]map[string]interface{}, 0)
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return nil, err
	}

	result["data"] = arr
	return result, nil
}

func (j *JMA) Format(latitude, longitude float64, data map[string]interface{}) ([]EarthquakeList, error) {
	var list []EarthquakeList
	for _, v := range data["data"].([]map[string]interface{}) {
		if !HasKey(v, []string{
			"anm", "mag", "cod", "at",
		}) {
			continue
		}

		ts, err := time.Parse("2006-01-02T15:04:05+09:00", v["at"].(string))
		if err != nil {
			continue
		}

		l := EarthquakeList{
			Depth:     j.GetDepth(v["cod"].(string)),
			Verfied:   true,
			Timestamp: ts.Add(-9 * time.Hour).UnixMilli(),
			Event:     v["anm"].(string),
			Region:    v["anm"].(string),
			Latitude:  j.GetLatitude(v["cod"].(string)),
			Longitude: j.GetLongitude(v["cod"].(string)),
			Magnitude: String2Float(v["mag"].(string)),
		}
		l.Distance = GetDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimated = GetEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (j *JMA) List(latitude, longitude float64) ([]EarthquakeList, error) {
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

func (j *JMA) GetDepth(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return String2Float(arr[2]) / 1000
}

func (j *JMA) GetLatitude(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return String2Float(arr[0])
}

func (j *JMA) GetLongitude(data string) float64 {
	arr := strings.FieldsFunc(data, func(c rune) bool {
		return c == '+' || c == '-' || c == '/'
	})
	if len(arr) < 3 {
		return 0
	}

	return String2Float(arr[1])
}
