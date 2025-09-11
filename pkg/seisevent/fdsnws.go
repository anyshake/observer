package seisevent

import (
	"encoding/csv"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/bclswl0827/travel"
)

func ParseFdsnwsEvent(travelTimeTable *travel.AK135, dataText, timeLayout string, latitude, longitude float64) ([]Event, error) {
	// Convert to CSV format
	csvDataStr := strings.ReplaceAll(dataText, ",", " - ")
	csvDataStr = strings.ReplaceAll(csvDataStr, "|", ",")
	csvRecords, err := csv.NewReader(strings.NewReader(csvDataStr)).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(csvRecords) <= 1 {
		return nil, errors.New("no seismic event found")
	}

	var resultArr []Event
	for _, record := range csvRecords[1:] {
		var (
			seisEvent Event
			magType   string
		)
		for idx, val := range record {
			switch idx {
			case 0:
				seisEvent.Event = val
			case 1:
				seisEvent.Verfied = true
				t, _ := time.Parse(timeLayout, val)
				seisEvent.Timestamp = t.UnixMilli()
			case 2:
				lat, _ := strconv.ParseFloat(val, 64)
				seisEvent.Latitude = lat
			case 3:
				lon, _ := strconv.ParseFloat(val, 64)
				seisEvent.Longitude = lon
			case 4:
				depth, _ := strconv.ParseFloat(val, 64)
				seisEvent.Depth = depth
			case 9:
				magType = val
			case 10:
				m, _ := strconv.ParseFloat(val, 64)
				seisEvent.Magnitude = []Magnitude{
					{Type: ParseMagnitude(magType), Value: m},
				}
			case 12:
				seisEvent.Region = val
			}
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	return resultArr, nil
}
