package seisevent

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const KRDAE_ID = "krdae"

type KRDAE struct {
	cache cache.AnyCache
}

func (c *KRDAE) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      KRDAE_ID,
		Country: "TR",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Kandilli Observatory and Earthquake Research Institute",
			"zh-TW": "坎迪利天文台和地震研究所",
			"zh-CN": "坎迪利天文台和地震研究所",
		},
	}
}

func (c *KRDAE) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"http://www.koeri.boun.edu.tr/scripts/lst4.asp",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(c.cache.Get().([]byte)))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find("pre").Each(func(i int, s *goquery.Selection) {
		table := strings.Split(s.Text(), "--------------")
		rows := strings.Split(table[len(table)-1], "\n")

		for idx, row := range rows {
			if len(row) > 0 {
				parsed := regexp.MustCompile(`\s+`).Split(row, 10)

				seisEvent := Event{
					Verfied:   true,
					Event:     fmt.Sprintf("#%d", idx),
					Latitude:  string2Float(parsed[2]),
					Longitude: string2Float(parsed[3]),
					Depth:     string2Float(parsed[4]),
					Timestamp: c.getTimestamp(parsed[0], parsed[1]),
					Magnitude: c.getMagnitude(parsed[5], parsed[6], parsed[7]),
					Region:    regexp.MustCompile(`\s+`).ReplaceAllString(parsed[8], " "),
				}

				seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
				seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

				resultArr = append(resultArr, seisEvent)
			}
		}
	})

	return sortSeismicEvents(resultArr), nil
}

func (c *KRDAE) getTimestamp(dateStr, timeStr string) int64 {
	loc, _ := time.LoadLocation("Europe/Istanbul")
	t, _ := time.ParseInLocation("2006.01.02 15:04:05", fmt.Sprintf("%s %s", dateStr, timeStr), loc)
	return t.UnixMilli()
}

func (c *KRDAE) getMagnitude(md, ml, mw string) []Magnitude {
	mdVal, mlVal, mwVal := string2Float(md), string2Float(ml), string2Float(mw)
	var result []Magnitude
	if mdVal != 0 {
		result = append(result, Magnitude{Type: ParseMagnitude("MD"), Value: mdVal})
	}
	if mlVal != 0 {
		result = append(result, Magnitude{Type: ParseMagnitude("ML"), Value: mlVal})
	}
	if mwVal != 0 {
		result = append(result, Magnitude{Type: ParseMagnitude("MW"), Value: mwVal})
	}
	return result
}
