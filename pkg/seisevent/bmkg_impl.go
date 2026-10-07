package seisevent

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const BMKG_ID = "bmkg"

type BMKG struct {
	resolvers       dnsquery.Resolvers
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
	transport       http.RoundTripper
}

func (c *BMKG) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      BMKG_ID,
		Country: "ID",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Meteorology, Climatology, and Geophysical Agency",
			"zh-TW": "印度尼西亞氣象、氣候和地球物理局",
		},
	}
}

func (c *BMKG) GetEvents(latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err, _ := c.sf.Do(BMKG_ID, func() (any, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			transport := c.transport
			if transport == nil {
				// Set custom frontend SNI (bmkg) to bypass GFW in China
				transport = createCustomTransport(c.resolvers, "bmkg")
			}
			res, err := request.GET(
				"https://bmkg-content-inatews.storage.googleapis.com/last30feltevent.xml",
				30*time.Second, time.Second, 3, false,
				transport,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			resultArr, err := c.parseXmlData(string(res))
			if err != nil {
				return nil, err
			}

			sorted := sortSeismicEvents(resultArr)
			c.cache.Set(sorted)
			return sorted, nil
		})
		if err != nil {
			return nil, err
		}

		baseEvents = v.([]Event)
	}

	for i := range baseEvents {
		baseEvents[i].Distance = getDistance(latitude, baseEvents[i].Latitude, longitude, baseEvents[i].Longitude)
		baseEvents[i].Estimation = getSeismicEstimation(
			c.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (c *BMKG) parseXmlData(xmlData string) ([]Event, error) {
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

		events = append(events, Event{
			Verfied:   true,
			Event:     item["eventid"],
			Region:    item["area"],
			Latitude:  lat,
			Longitude: lng,
			Depth:     c.getDepth(item["depth"]),
			Magnitude: c.getMagnitude(item["magnitude"]),
			Timestamp: timestamp,
		})
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
