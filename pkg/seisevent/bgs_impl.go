package seisevent

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const BGS_ID = "bgs"

type BGS struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *BGS) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      BGS_ID,
		Country: "GB",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "British Geological Survey",
			"zh-TW": "英國地質調查局",
			"zh-CN": "英国地质调查局",
		},
	}
}

func (c *BGS) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	worldEventRes, err := request.GET(
		"https://quakes.bgs.ac.uk/feeds/WorldSeismology.xml",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get WorldSeismology.xml: %w", err)
	}
	ukEventRes, err := request.GET(
		"https://quakes.bgs.ac.uk/feeds/MhSeismology.xml",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get MhSeismology.xml: %w", err)
	}

	worldEvent, err := c.parseXmlData("WorldSeismology", string(worldEventRes), latitude, longitude)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WorldSeismology.xml: %w", err)
	}
	ukEvent, err := c.parseXmlData("MhSeismology", string(ukEventRes), latitude, longitude)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MhSeismology.xml: %w", err)
	}

	sortedEvents := sortSeismicEvents(append(worldEvent, ukEvent...))
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *BGS) parseXmlData(tag, xmlData string, latitude, longitude float64) ([]Event, error) {
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
			if currentElement == "item" {
				current = make(map[string]string)
			}
		case xml.EndElement:
			if tok.Name.Local == "item" && current != nil {
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
	for idx, item := range items {
		seisEvent := Event{
			Event:   fmt.Sprintf("%s-%d", tag, idx),
			Verfied: true,
		}
		for key, val := range item {
			switch key {
			case "description":
				// Origin date/time: Sun, 10 Aug 2025 16:53:46 ; Location: WESTERN TURKEY ; Lat/long: 39.312,28.069 ; Depth: 10 km ; Magnitude: 6.1
				splitedFields := strings.Split(val, ";")
				for fieldIdx, field := range splitedFields {
					field = strings.TrimSpace(field)
					switch fieldIdx {
					case 0:
						timestamp, err := c.getTimestamp(field)
						if err != nil {
							continue
						}
						seisEvent.Timestamp = timestamp
					case 1:
						seisEvent.Region = c.getRegion(field)
					case 3:
						depth, err := c.getDepth(field)
						if err != nil {
							continue
						}
						seisEvent.Depth = depth
					case 4:
						depth, err := c.getMagnitude(field)
						if err != nil {
							continue
						}
						seisEvent.Magnitude = depth
					}
				}
			case "lat":
				seisEvent.Latitude = string2Float(val)
			case "lon":
				seisEvent.Longitude = string2Float(val)
			}
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		events = append(events, seisEvent)
	}

	return events, nil
}

func (c *BGS) getTimestamp(data string) (int64, error) {
	fields := strings.Split(data, ",")
	if len(fields) != 2 {
		return 0, errors.New("failed to parse timestamp")
	}

	timeStr := strings.TrimSpace(fields[1])
	tm, err := time.Parse("02 Jan 2006 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return tm.UnixMilli(), nil
}

func (c *BGS) getRegion(data string) string {
	return strings.TrimSpace(strings.Replace(data, "Location:", "", -1))
}

func (c *BGS) getDepth(data string) (float64, error) {
	fields := strings.Split(data, ":")
	if len(fields) != 2 {
		return 0, errors.New("failed to parse depth")
	}

	depthStr := strings.TrimSpace(strings.ReplaceAll(fields[1], "km", ""))
	return string2Float(depthStr), nil
}

func (c *BGS) getMagnitude(data string) ([]Magnitude, error) {
	fields := strings.Split(data, ":")
	if len(fields) != 2 {
		return nil, errors.New("failed to parse magnitude")
	}

	depthStr := strings.TrimSpace(fields[1])
	return []Magnitude{{Type: ParseMagnitude("M"), Value: string2Float(depthStr)}}, nil
}
