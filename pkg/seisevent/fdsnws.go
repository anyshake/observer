package seisevent

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"
)

func ParseFdsnwsEvent(dataText, timeLayout string) ([]Event, error) {
	reader := csv.NewReader(strings.NewReader(dataText))
	reader.Comma = '|'
	reader.LazyQuotes = true
	csvRecords, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(csvRecords) <= 1 {
		return nil, errors.New("no seismic event found")
	}

	var resultArr []Event
	for _, record := range csvRecords[1:] {
		if len(record) < 13 {
			return nil, errors.New("incomplete seismic event record")
		}
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
		if strings.TrimSpace(seisEvent.Region) == "" {
			seisEvent.Region = fmt.Sprintf("Latitude: %.4f°, Longitude: %.4f°", seisEvent.Latitude, seisEvent.Longitude)
		}
		resultArr = append(resultArr, seisEvent)
	}

	return resultArr, nil
}
