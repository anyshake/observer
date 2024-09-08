package seisevent

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const NCS_ID = "ncs"

type NCS struct {
	cache cache.BytesCache
}

func (c *NCS) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      NCS_ID,
		Country: "IN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "National Center for Seismology",
			"zh-TW": "印度國家地震中心",
			"zh-CN": "印度国家地震中心",
		},
	}
}

func (c *NCS) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://riseq.seismo.gov.in/riseq/earthquake",
			30*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse NCS HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(c.cache.Get()))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find("#sidebar-wrapper").Each(func(i int, s *goquery.Selection) {
		s.Find(".event_list").Each(func(i int, s *goquery.Selection) {
			dataJson, ok := s.Attr("data-json")
			if !ok {
				return
			}

			var dataMap map[string]any
			err := json.Unmarshal([]byte(dataJson), &dataMap)
			if err != nil {
				return
			}

			eventId, ok := dataMap["event_id"].(string)
			if !ok {
				return
			}
			eventName, ok := dataMap["event_name"].(string)
			if !ok {
				return
			}
			eventTime, ok := dataMap["origin_time"].(string)
			if !ok {
				return
			}
			eventLatLon, ok := dataMap["lat_long"].(string)
			if !ok {
				return
			}
			eventMagDepth, ok := dataMap["magnitude_depth"].(string)
			if !ok {
				return
			}
			eventType, ok := dataMap["event_type"].(string)
			if !ok {
				return
			}

			seisEvent := Event{
				Event:     eventId,
				Region:    eventName,
				Verfied:   eventType == "Reviewed",
				Latitude:  c.getLatitude(eventLatLon),
				Longitude: c.getLongitude(eventLatLon),
				Magnitude: c.getMagnitude(eventMagDepth),
				Depth:     c.getDepth(eventMagDepth),
				Timestamp: c.getTimestamp(eventTime),
			}
			seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
			seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

			resultArr = append(resultArr, seisEvent)
		})
	})

	return sortSeismicEvents(resultArr), nil
}

func (c *NCS) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05 IST", data)
	return t.Add(-(5*time.Hour + 30*time.Minute)).UnixMilli()
}

func (c *NCS) getLatitude(data string) float64 {
	coords := strings.Split(data, ",")
	if len(coords) == 2 {
		return string2Float(strings.TrimSpace(coords[0]))
	}

	return -1
}

func (c *NCS) getLongitude(data string) float64 {
	coords := strings.Split(data, ",")
	if len(coords) == 2 {
		return string2Float(strings.TrimSpace(coords[1]))
	}

	return -1
}

func (c *NCS) getDepth(data string) float64 {
	depthStr := strings.Split(data, ",")
	if len(depthStr) == 2 {
		depth := strings.Split(depthStr[1], ":")
		if len(depth) == 2 {
			depth[1] = strings.ReplaceAll(depth[1], "km", "")
			return string2Float(strings.TrimSpace(depth[1]))
		}
	}

	return -1
}

func (c *NCS) getMagnitude(data string) []Magnitude {
	magStr := strings.Split(data, ",")
	if len(magStr) == 2 {
		mag := strings.Split(magStr[0], ":")
		if len(mag) == 2 {
			return []Magnitude{{
				Type:  ParseMagnitude("M"),
				Value: string2Float(strings.TrimSpace(mag[1])),
			}}
		}
	}

	return []Magnitude{}
}
