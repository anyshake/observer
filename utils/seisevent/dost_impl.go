package seisevent

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const DOST_ID = "dost"

type DOST struct {
	cache cache.AnyCache
}

func (c *DOST) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      DOST_ID,
		Country: "PH",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Philippine Institute of Volcanology and Seismology",
			"zh-TW": "菲律賓火山地震研究所",
			"zh-CN": "菲律宾火山地震研究所",
		},
	}
}

func (c *DOST) GetEvents(latitude, longitude float64) ([]Event, error) {
	// Get DOST HTML response
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://earthquake.phivolcs.dost.gov.ph/",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse DOST HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(c.cache.Get().([]byte)))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find("table").Each(func(i int, s *goquery.Selection) {
		htmlText := s.Text()
		if strings.Contains(htmlText, "Latitude") &&
			strings.Contains(htmlText, "Longitude") &&
			strings.Contains(htmlText, "Depth") &&
			strings.Contains(htmlText, "Mag") &&
			strings.Contains(htmlText, "Location") {
			var seisEvent Event
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch i % 6 {
				case 0:
					seisEvent.Timestamp = c.getTimestamp(s.Text())
				case 1:
					seisEvent.Latitude = c.getLatitude(s.Text())
				case 2:
					seisEvent.Longitude = c.getLongitude(s.Text())
				case 3:
					seisEvent.Depth = c.getDepth(s.Text())
				case 4:
					seisEvent.Magnitude = c.getMagnitude(s.Text())
				case 5:
					seisEvent.Verfied = true
					seisEvent.Event = fmt.Sprintf("%d", i/6+1)
					seisEvent.Region = c.getRegion(s.Text())
					seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
					seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)
					resultArr = append(resultArr, seisEvent)
				}
			})
		}
	})

	return sortSeismicEvents(resultArr), nil
}

func (c *DOST) getTimestamp(data string) int64 {
	data = strings.TrimSpace(data)
	t, err := time.Parse("2 January 2006 - 3:04 PM", data)
	if err != nil {
		return 0
	}

	return t.Add(-8 * time.Hour).UnixMilli()
}

func (c *DOST) getLatitude(data string) float64 {
	data = strings.TrimSpace(data)
	return string2Float(strings.TrimPrefix(data, "0"))
}

func (c *DOST) getLongitude(data string) float64 {
	data = strings.TrimSpace(data)
	return string2Float(strings.TrimPrefix(data, "0"))
}

func (c *DOST) getDepth(data string) float64 {
	data = strings.TrimSpace(data)
	return string2Float(strings.TrimPrefix(data, "0"))
}

func (c *DOST) getMagnitude(data string) []Magnitude {
	data = strings.TrimSpace(data)
	return []Magnitude{{Value: string2Float(data), Type: "MS"}}
}

func (c *DOST) getRegion(data string) string {
	data = strings.TrimSpace(data)
	data = strings.ReplaceAll(data, "\t", "")
	data = strings.ReplaceAll(data, "\n", "")
	return data
}
