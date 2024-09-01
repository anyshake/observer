package trace

import (
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"github.com/sbabiv/xml2map"
)

type HKO struct {
	dataSourceCache
}

func (h *HKO) Property() string {
	return "天文台全球地震資訊網"
}

func (h *HKO) Fetch() ([]byte, error) {
	if time.Since(h.Time) <= EXPIRATION {
		return h.Cache, nil
	}

	res, err := request.GET(
		"https://www.hko.gov.hk/gts/QEM/eq_app-30d_uc.xml",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	h.Time = time.Now()
	h.Cache = make([]byte, len(res))
	copy(h.Cache, res)

	return res, nil
}

func (h *HKO) Parse(data []byte) (map[string]any, error) {
	decoder := xml2map.NewDecoder(strings.NewReader(string(data)))

	result, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *HKO) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	events, ok := data["Earthquake"].(map[string]any)["EventGroup"].(map[string]any)["Event"]
	if !ok {
		return nil, fmt.Errorf("source data is not valid")
	}

	var list []seismicEvent
	for _, v := range events.([]map[string]any) {
		if !isMapHasKeys(v, []string{
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

		l := seismicEvent{
			Depth:     -1,
			Verfied:   v["Verify"].(string) == "Y",
			Timestamp: ts.Add(-8 * time.Hour).UnixMilli(),
			Event:     v["City"].(string),
			Region:    v["Region"].(string),
			Latitude:  string2Float(v["Lat"].(string)),
			Longitude: string2Float(v["Lon"].(string)),
			Magnitude: string2Float(v["Mag"].(string)),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getSeismicEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (h *HKO) List(latitude, longitude float64) ([]seismicEvent, error) {
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
