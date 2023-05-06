package trace

import (
	"fmt"
	"strings"
	"time"

	"com.geophone.observer/common/request"
	"github.com/sbabiv/xml2map"
)

type HKO struct{}

func (h *HKO) Property() (string, string) {
	const (
		NAME  string = "香港天文台全球地震資訊網"
		VALUE string = "HKO"
	)

	return NAME, VALUE
}

func (h *HKO) Fetch() ([]byte, error) {
	res, err := request.GETRequest(
		"https://www.hko.gov.hk/gts/QEM/eq_app-30d_uc.xml",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *HKO) Parse(data []byte) (map[string]interface{}, error) {
	decoder := xml2map.NewDecoder(strings.NewReader(string(data)))

	result, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *HKO) Format(latitude, longitude float64, data map[string]interface{}) ([]EarthquakeList, error) {
	events, ok := data["Earthquake"].(map[string]interface{})["EventGroup"].(map[string]interface{})["Event"]
	if !ok {
		return nil, fmt.Errorf("source data is not valid")
	}

	var list []EarthquakeList
	for _, v := range events.([]map[string]interface{}) {
		if !HasKey(v, []string{
			"Verify", "HKTDate", "HKTTime", "City",
			"Region", "Lat", "Lon", "Mag",
		}) {
			continue
		}

		ts, err := time.Parse("20060102 150400", fmt.Sprintf(
			"%s %s00", v["HKTDate"].(string), v["HKTTime"].(string)),
		)
		if err != nil {
			continue
		}

		l := EarthquakeList{
			Depth:     -1,
			Verfied:   v["Verify"].(string) == "Y",
			Timestamp: ts.Add(-8 * time.Hour).UnixMilli(),
			Event:     v["City"].(string),
			Region:    v["Region"].(string),
			Latitude:  String2Float(v["Lat"].(string)),
			Longitude: String2Float(v["Lon"].(string)),
			Magnitude: String2Float(v["Mag"].(string)),
		}
		l.Distance = GetDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimated = GetEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (h *HKO) List(latitude, longitude float64) ([]EarthquakeList, error) {
	res, err := h.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := h.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := h.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}
