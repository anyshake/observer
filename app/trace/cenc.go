package trace

import (
	"encoding/json"
	"fmt"
	"time"

	"com.geophone.observer/common/request"
)

type CENC struct{}

func (c *CENC) Property() (string, string) {
	const (
		NAME  string = "中国地震台网地震信息"
		VALUE string = "CENC"
	)

	return NAME, VALUE
}

func (c *CENC) Fetch() ([]byte, error) {
	res, err := request.GETRequest(
		"https://api.wolfx.jp/cenc_eqlist.json",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *CENC) Parse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CENC) Format(latitude, longitude float64, data map[string]interface{}) ([]EarthquakeList, error) {
	var list []EarthquakeList
	for k, v := range data {
		if k == "No0" {
			continue
		}

		if !HasKey(v.(map[string]interface{}), []string{
			"depth", "latitude", "location",
			"longitude", "magnitude", "time",
		}) {
			continue
		}

		ts, err := time.Parse("2006-01-02 15:04:05", v.(map[string]interface{})["time"].(string))
		if err != nil {
			fmt.Println(err)
			continue
		}

		l := EarthquakeList{
			Depth:     String2Float(v.(map[string]interface{})["depth"].(string)),
			Verfied:   true,
			Timestamp: ts.Add(-8 * time.Hour).UnixMilli(),
			Event:     v.(map[string]interface{})["location"].(string),
			Region:    v.(map[string]interface{})["location"].(string),
			Latitude:  String2Float(v.(map[string]interface{})["latitude"].(string)),
			Longitude: String2Float(v.(map[string]interface{})["longitude"].(string)),
			Magnitude: String2Float(v.(map[string]interface{})["magnitude"].(string)),
		}
		l.Distance = GetDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimated = GetEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (c *CENC) List(latitude, longitude float64) ([]EarthquakeList, error) {
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
