package seisevent

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"github.com/sbabiv/xml2map"
)

type HKO struct {
	cache cache.BytesCache
}

func (h *HKO) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      HKO_ID,
		Country: "HK",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Hong Kong Observatory",
			"zh-TW": "天文台全球地震資訊網",
			"zh-CN": "天文台全球地震信息网",
		},
	}
}

func (h *HKO) GetEvents(latitude, longitude float64) ([]Event, error) {
	if h.cache.Valid() {
		res, err := request.GET(
			"https://www.hko.gov.hk/gts/QEM/eq_app-30d_uc.xml",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		h.cache.Set(res)
	}

	// Parse HKO XML response
	dataMap, err := xml2map.NewDecoder(strings.NewReader(string(h.cache.Get()))).Decode()
	if err != nil {
		return nil, err
	}
	dataMapEvents, ok := dataMap["Earthquake"].(map[string]any)["EventGroup"].(map[string]any)["Event"]
	if !ok {
		return nil, errors.New("source data is not valid, missing Earthquake.EventGroup.Event")
	}

	// Ensure the response has the expected keys
	expectedKeys := []string{"Verify", "HKTDate", "HKTTime", "City", "Region", "Lat", "Lon", "Mag"}

	var resultArr []Event
	for _, v := range dataMapEvents.([]map[string]any) {
		if !isMapHasKeys(v, expectedKeys) {
			continue
		}

		timestamp, err := h.getTimestamp(fmt.Sprintf("%s %s00", v["HKTDate"].(string), v["HKTTime"].(string)))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Depth:     -1,
			Timestamp: timestamp,
			Event:     v["City"].(string),
			Verfied:   v["Verify"].(string) == "Y",
			Region:    v["Region"].(string),
			Latitude:  string2Float(v["Lat"].(string)),
			Longitude: string2Float(v["Lon"].(string)),
			Magnitude: string2Float(v["Mag"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (h *HKO) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("20060102 150405", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-8 * time.Hour).UnixMilli(), nil
}
