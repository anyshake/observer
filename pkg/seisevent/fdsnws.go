package seisevent

import (
	"encoding/csv"
	"errors"
	"fmt"
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
				if len(val) > len(timeLayout) {
					val = val[:len(timeLayout)]
				}
				t, err := time.Parse(timeLayout, val)
				if err != nil {
					return nil, err
				}
				seisEvent.Timestamp = t.UnixMilli()
			case 2:
				seisEvent.Latitude = string2Float(val)
			case 3:
				seisEvent.Longitude = string2Float(val)
			case 4:
				seisEvent.Depth = string2Float(val)
			case 9:
				magType = val
			case 10:
				seisEvent.Magnitude = []Magnitude{
					{Type: ParseMagnitude(magType), Value: string2Float(val)},
				}
			case 12:
				seisEvent.Region = val
			}
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		fmt.Printf("%+v\n", seisEvent)
		resultArr = append(resultArr, seisEvent)
	}

	return resultArr, nil
}
