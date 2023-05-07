package trace

import (
	"encoding/json"
	"fmt"
	"time"

	"com.geophone.observer/common/request"
)

type USGS struct{}

func (u *USGS) Property() (string, string) {
	const (
		NAME  string = "United States Geological Survey"
		VALUE string = "USGS"
	)

	return NAME, VALUE
}

func (u *USGS) Fetch() ([]byte, error) {
	res, err := request.GETRequest(
		"https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_day.geojson",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *USGS) Parse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *USGS) Format(latitude, longitude float64, data map[string]interface{}) ([]EarthquakeList, error) {
	events, ok := data["features"]
	if !ok {
		return nil, fmt.Errorf("source data is not valid")
	}

	var list []EarthquakeList
	for _, v := range events.([]interface{}) {
		if !HasKey(v.(map[string]interface{}), []string{"properties"}) {
			continue
		}

		properties := v.(map[string]interface{})["properties"]
		if !HasKey(properties.(map[string]interface{}), []string{
			"mag", "place", "time", "type", "title",
		}) {
			continue
		}

		geometry := v.(map[string]interface{})["geometry"]
		if !HasKey(geometry.(map[string]interface{}), []string{"coordinates"}) {
			continue
		}

		coordinates := geometry.(map[string]interface{})["coordinates"]
		if len(coordinates.([]interface{})) != 3 {
			continue
		}

		if properties.(map[string]interface{})["type"].(string) != "earthquake" {
			continue
		}

		l := EarthquakeList{
			Depth:     coordinates.([]interface{})[2].(float64),
			Verfied:   true,
			Timestamp: int64(properties.(map[string]interface{})["time"].(float64)),
			Event:     properties.(map[string]interface{})["title"].(string),
			Region:    properties.(map[string]interface{})["place"].(string),
			Latitude:  coordinates.([]interface{})[1].(float64),
			Longitude: coordinates.([]interface{})[0].(float64),
			Magnitude: properties.(map[string]interface{})["mag"].(float64),
		}
		l.Distance = GetDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimated = GetEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (u *USGS) List(latitude, longitude float64) ([]EarthquakeList, error) {
	res, err := u.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := u.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := u.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}
