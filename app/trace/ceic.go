package trace

import (
	"encoding/json"
	"time"

	"com.geophone.observer/utils/request"
)

type CEIC struct{}

func (c *CEIC) Property() (string, string) {
	const (
		NAME  string = "中国地震台网中心"
		VALUE string = "CEIC"
	)

	return NAME, VALUE
}

func (c *CEIC) Fetch() ([]byte, error) {
	res, err := request.GET(
		"https://news.ceic.ac.cn/ajax/google",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *CEIC) Parse(data []byte) (map[string]any, error) {
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

func (c *CEIC) Format(latitude, longitude float64, data map[string]any) ([]EarthquakeList, error) {
	var list []EarthquakeList
	for _, v := range data["data"].([]map[string]any) {
		if !hasKey(v, []string{
			"O_TIME", "EPI_LAT", "EPI_LON", "EPI_DEPTH", "M", "LOCATION_C",
		}) {
			continue
		}

		if !isEmpty(v, []string{
			"O_TIME", "EPI_LAT", "EPI_LON", "EPI_DEPTH", "M", "LOCATION_C",
		}) {
			continue
		}

		ts, err := time.Parse("2006-01-02 15:04:05", v["O_TIME"].(string))
		if err != nil {
			continue
		}

		l := EarthquakeList{
			Depth:     string2Float(v["EPI_DEPTH"].(string)),
			Verfied:   true,
			Timestamp: ts.Add(-8 * time.Hour).UnixMilli(),
			Event:     v["LOCATION_C"].(string),
			Region:    v["LOCATION_C"].(string),
			Latitude:  string2Float(v["EPI_LAT"].(string)),
			Longitude: string2Float(v["EPI_LON"].(string)),
			Magnitude: string2Float(v["M"].(string)),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimated = getEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (c *CEIC) List(latitude, longitude float64) ([]EarthquakeList, error) {
	res, err := c.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := c.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := c.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}
