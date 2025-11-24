package seisevent

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const BMKG_ID = "bmkg"

type BMKG struct {
	resolvers       dnsquery.Resolvers
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *BMKG) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      BMKG_ID,
		Country: "ID",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Meteorology, Climatology, and Geophysical Agency",
			"zh-TW": "印度尼西亞氣象、氣候和地球物理局",
			"zh-CN": "印度尼西亚气象、气候和地球物理局",
		},
	}
}

func (c *BMKG) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://bmkg-content-inatews.storage.googleapis.com/last30feltevent.xml",
		30*time.Second, time.Second, 3, false,
		// Set custom frontend SNI (bmkg) to bypass GFW in China
		createCustomTransport(c.resolvers, "bmkg"),
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	resultArr, err := c.parseXmlData(string(res), latitude, longitude)
	if err != nil {
		return nil, err
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *BMKG) parseXmlData(xmlData string, latitude, longitude float64) ([]Event, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlData))

	var (
		items          []map[string]string
		current        map[string]string
		currentElement string
	)
	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			currentElement = tok.Name.Local
			if currentElement == "info" {
				current = make(map[string]string)
			}
		case xml.EndElement:
			if tok.Name.Local == "info" && current != nil {
				items = append(items, current)
				current = nil
			}
			currentElement = ""
		case xml.CharData:
			if current != nil && currentElement != "" {
				text := strings.TrimSpace(string(tok))
				if text != "" {
					current[currentElement] = text
				}
			}
		}
	}

	var events []Event
	mapKeys := []string{"eventid", "date", "time", "magnitude", "depth", "area", "coordinates"}

	for _, item := range items {
		if !isMapHasKeys(item, mapKeys) {
			continue
		}

		lat, lng, err := c.getCoordinates(item["coordinates"])
		if err != nil {
			return nil, err
		}
		timestamp, err := c.getTimestamp(item["date"], item["time"])
		if err != nil {
			return nil, err
		}
		seisEvent := Event{
			Verfied:   true,
			Event:     item["eventid"],
			Region:    item["area"],
			Latitude:  lat,
			Longitude: lng,
			Depth:     c.getDepth(item["depth"]),
			Magnitude: c.getMagnitude(item["magnitude"]),
			Timestamp: timestamp,
		}

		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		events = append(events, seisEvent)
	}

	return events, nil
}

func (c *BMKG) getCoordinates(data string) (float64, float64, error) {
	split := strings.Split(data, ",")
	if len(split) != 2 {
		return 0, 0, errors.New("failed to parse coordinates")
	}
	return string2Float(split[1]), string2Float(split[0]), nil
}

func (c *BMKG) getTimestamp(dateStr, timeStr string) (int64, error) {
	t, err := time.Parse("02-01-06 15:04:05 WIB", fmt.Sprintf("%s %s", dateStr, timeStr))
	if err != nil {
		return 0, err
	}

	return t.Add(-7 * time.Hour).UnixMilli(), nil
}

func (c *BMKG) getDepth(data string) float64 {
	depthVal := strings.TrimSpace(strings.Replace(data, "Km", "", -1))
	return string2Float(depthVal)
}

func (c *BMKG) getMagnitude(data string) []Magnitude {
	return []Magnitude{
		{Type: ParseMagnitude("M"), Value: string2Float(data)},
	}
}
