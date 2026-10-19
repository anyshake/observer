package seisevent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const SCEA_ID = "scea"
const SCEAListURL = "https://api.wolfx.jp/sc_eew_list.json"

type SCEA struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
	timeSource      *timesource.Source
}

type sceaListResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []sceaEvent `json:"data"`
}

type sceaEvent struct {
	EventID      string   `json:"eventId"`
	ShockTime    int64    `json:"shockTime"`
	Longitude    float64  `json:"longitude"`
	Latitude     float64  `json:"latitude"`
	PlaceName    string   `json:"placeName"`
	Magnitude    float64  `json:"magnitude"`
	Depth        *float64 `json:"depth"`
	InfoTypeName string   `json:"infoTypeName"`
}

func (s *SCEA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      SCEA_ID,
		Country: "CN",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Sichuan Earthquake Administration Early Warning (Wolfx)",
			"zh-TW": "四川地震局預警（Wolfx）",
		},
	}
}

func (s *SCEA) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if s.cache.Valid() {
		baseEvents = s.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &s.sf, SCEA_ID, func(ctx context.Context) ([]Event, error) {
			if s.cache.Valid() {
				return s.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, SCEAListURL,
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			var data sceaListResponse
			if err := json.Unmarshal(res, &data); err != nil {
				return nil, err
			}

			if data.Code != 0 {
				return nil, fmt.Errorf("server error: %s", data.Msg)
			}

			if len(data.Data) == 0 {
				return nil, errors.New("seismic event data is not available")
			}

			resultArr := make([]Event, 0, len(data.Data))
			for _, event := range data.Data {
				if event.EventID == "" || event.ShockTime == 0 || event.PlaceName == "" {
					continue
				}

				depth := -1.0
				if event.Depth != nil {
					depth = *event.Depth
				}

				resultArr = append(resultArr, Event{
					Depth:     depth,
					Verfied:   event.InfoTypeName == "[正式]",
					Event:     event.EventID,
					Region:    event.PlaceName,
					Latitude:  event.Latitude,
					Longitude: event.Longitude,
					Magnitude: s.getMagnitude(event.Magnitude),
					Timestamp: time.UnixMilli(event.ShockTime).UnixMilli(),
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			s.cache.Set(sortedEvents)
			return sortedEvents, nil
		})
		if err != nil {
			return nil, err
		}

		baseEvents = v
	}

	baseEvents = slices.Clone(baseEvents)
	for i := range baseEvents {
		baseEvents[i].Distance = getDistance(latitude, baseEvents[i].Latitude, longitude, baseEvents[i].Longitude)
		baseEvents[i].Estimation = getSeismicEstimation(
			s.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (s *SCEA) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
