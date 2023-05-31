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

func (u *USGS) Parse(data []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *USGS) Format(latitude, longitude float64, data map[string]any) ([]EarthquakeList, error) {
	events, ok := data["features"]
	if !ok {
		return nil, fmt.Errorf("source data is not valid")
	}

	var list []EarthquakeList
	for _, v := range events.([]any) {
		if !HasKey(v.(map[string]any), []string{"properties"}) {
			continue
		}

		properties := v.(map[string]any)["properties"]
		if !HasKey(properties.(map[string]any), []string{
			"mag", "place", "time", "type", "title",
		}) {
			continue
		}

		geometry := v.(map[string]any)["geometry"]
		if !HasKey(geometry.(map[string]any), []string{"coordinates"}) {
			continue
		}

		coordinates := geometry.(map[string]any)["coordinates"]
		if len(coordinates.([]any)) != 3 {
			continue
		}

		if properties.(map[string]any)["type"].(string) != "earthquake" {
			continue
		}

		l := EarthquakeList{
			Depth:     coordinates.([]any)[2].(float64),
			Verfied:   true,
			Timestamp: int64(properties.(map[string]any)["time"].(float64)),
			Event:     properties.(map[string]any)["title"].(string),
			Region:    properties.(map[string]any)["place"].(string),
			Latitude:  coordinates.([]any)[1].(float64),
			Longitude: coordinates.([]any)[0].(float64),
			Magnitude: properties.(map[string]any)["mag"].(float64),
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
