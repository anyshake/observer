package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const CENC_WEB_ID = "cenc_web"

type CENC_WEB struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
	timeSource      *timesource.Source
}

func (c *CENC_WEB) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CENC_WEB_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "China Earthquake Networks Center (Web)",
			"zh-TW": "中國地震台網中心（網頁端）",
			"zh-CN": "中国地震台网中心（网页端）",
		},
	}
}

func (c *CENC_WEB) getRequestParam(currentTime time.Time) string {
	startTime := url.QueryEscape(currentTime.AddDate(0, 0, -30).UTC().Add(8 * time.Hour).Format("2006-01-02 15:04:05"))
	endTime := url.QueryEscape(currentTime.UTC().Add(8 * time.Hour).Format("2006-01-02 15:04:05"))
	return fmt.Sprintf("orderBy=id&isAsc=false&startMg=1&startTime=%s&endTime=%s", startTime, endTime)
}

func (c *CENC_WEB) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		fmt.Sprintf("https://www.cenc.ac.cn/prodlaunch-web-backend/open/data/geojson/catalogs?%s", c.getRequestParam(c.timeSource.Now())),
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse CENC GeoJSON response
	var dataMap map[string]any
	err = json.Unmarshal(res, &dataMap)
	if err != nil {
		return nil, err
	}

	dataMapEvents, ok := dataMap["features"].([]any)
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	// Ensure the response has the expected keys and values
	expectedKeysInDataMap := []string{"properties", "geometry", "id"}
	expectedKeysInProperties := []string{"震级", "参考位置", "发震时刻", "深度（千米）"}
	expectedKeysInGeometry := []string{"coordinates"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event.(map[string]any), expectedKeysInDataMap) {
			continue
		}

		var (
			properties = event.(map[string]any)["properties"]
			geometry   = event.(map[string]any)["geometry"]
		)
		if !isMapHasKeys(properties.(map[string]any), expectedKeysInProperties) || !isMapHasKeys(geometry.(map[string]any), expectedKeysInGeometry) {
			continue
		}
		coordinates := geometry.(map[string]any)["coordinates"]

		ts, err := c.getTimestamp(properties.(map[string]any)["发震时刻"].(string))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: ts,
			Depth:     c.getDepth(properties.(map[string]any)["深度（千米）"]),
			Event:     event.(map[string]any)["id"].(string),
			Region:    properties.(map[string]any)["参考位置"].(string),
			Latitude:  coordinates.([]any)[1].(float64),
			Longitude: coordinates.([]any)[0].(float64),
			Magnitude: []Magnitude{c.getMagnitude("M", properties.(map[string]any)["震级"].(string))},
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *CENC_WEB) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-8 * time.Hour).UnixMilli(), nil
}

func (c *CENC_WEB) getDepth(depth any) float64 {
	switch d := depth.(type) {
	case string:
		return string2Float(d)
	case float64:
		return d
	}

	return -1
}

func (c *CENC_WEB) getMagnitude(magType, data string) Magnitude {
	return Magnitude{Type: ParseMagnitude(magType), Value: string2Float(data)}
}
