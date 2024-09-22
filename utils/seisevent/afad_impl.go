package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const AFAD_ID = "afad"

type AFAD struct {
	cache       cache.BytesCache
	currentTime time.Time
}

func (c *AFAD) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      AFAD_ID,
		Country: "TR",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Ministry of Disaster and Emergency Management",
			"zh-TW": "土耳其災害與應變管理署",
			"zh-CN": "土耳其灾害和应急管理署",
		},
	}
}

func (c *AFAD) getRequestParam() string {
	startTime := url.QueryEscape(c.currentTime.AddDate(0, 0, -5).UTC().Format("2006-01-02 15:04:05"))
	endTime := url.QueryEscape(c.currentTime.UTC().Format("2006-01-02 15:04:05"))
	return fmt.Sprintf("start=%s&end=%s&orderby=timedesc", startTime, endTime)
}

func (c *AFAD) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		// Get current time to construct the request parameter
		res, err := request.GET(
			"https://www.cloudflare.com/cdn-cgi/trace",
			10*time.Second, time.Second, 3, false,
			nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		ts := regexp.MustCompile(`ts=(\d+)`).FindStringSubmatch(string(res))
		if len(ts) == 0 {
			return nil, errors.New("failed to get current time from Cloudflare")
		}
		c.currentTime = time.Unix(int64(string2Float(ts[1])), 0)

		// Make AFAD API request
		res, err = request.GET(
			fmt.Sprintf("https://deprem.afad.gov.tr/apiv2/event/filter?%s", c.getRequestParam()),
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse AFAD JSON response
	var dataMapEvents []map[string]any
	err := json.Unmarshal(c.cache.Get(), &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"eventID", "location", "latitude", "longitude", "depth", "type", "magnitude", "date"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		timestamp, err := c.getTimestamp(event["date"].(string))
		if err != nil {
			return nil, err
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Event:     event["eventID"].(string),
			Region:    event["location"].(string),
			Depth:     string2Float(event["depth"].(string)),
			Latitude:  string2Float(event["latitude"].(string)),
			Longitude: string2Float(event["longitude"].(string)),
			Magnitude: c.getMagnitudeType(event["type"].(string), event["magnitude"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (c *AFAD) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.UnixMilli(), nil
}

func (c *AFAD) getMagnitudeType(magType, magText string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude(magType), Value: string2Float(magText)}}
}
