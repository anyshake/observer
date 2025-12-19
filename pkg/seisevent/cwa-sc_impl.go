package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const CWA_SC_ID = "cwa_sc"

type CWA_SC struct {
	resolvers       dnsquery.Resolvers
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *CWA_SC) getRequestBody(limit int) string {
	return fmt.Sprintf(`draw=1&columns[0][data]=0&columns[0][name]=EventNo&columns[0][searchable]=false&columns[0][orderable]=true&columns[0][search][value]=&columns[0][search][regex]=false&columns[1][data]=1&columns[1][name]=MaxIntensity&columns[1][searchable]=true&columns[1][orderable]=true&columns[1][search][value]=&columns[1][search][regex]=false&columns[2][data]=2&columns[2][name]=OriginTime&columns[2][searchable]=true&columns[2][orderable]=true&columns[2][search][value]=&columns[2][search][regex]=false&columns[3][data]=3&columns[3][name]=MagnitudeValue&columns[3][searchable]=true&columns[3][orderable]=true&columns[3][search][value]=&columns[3][search][regex]=false&columns[4][data]=4&columns[4][name]=Depth&columns[4][searchable]=true&columns[4][orderable]=true&columns[4][search][value]=&columns[4][search][regex]=false&columns[5][data]=5&columns[5][name]=Description&columns[5][searchable]=true&columns[5][orderable]=true&columns[5][search][value]=&columns[5][search][regex]=false&columns[6][data]=6&columns[6][name]=Description&columns[6][searchable]=true&columns[6][orderable]=true&columns[6][search][value]=&columns[6][search][regex]=false&order[0][column]=2&order[0][dir]=desc&start=0&length=%d&search[value]=&search[regex]=false&Search=&txtSDate=&txtEDate=&txtSscale=&txtEscale=&txtSdepth=&txtEdepth=&txtLonS=&txtLonE=&txtLatS=&txtLatE=&ddlCity=&ddlCitySta=&txtIntensityB=&txtIntensityE=&txtLon=&txtLat=&txtKM=&ddlStationName=&cblEventNo=&txtSDatePWS=&txtEDatePWS=&txtSscalePWS=&txtEscalePWS=&ddlMark=`, limit)
}

func (c *CWA_SC) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CWA_SC_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Central Weather Administration Seismological Center",
			"zh-TW": "交通部中央氣象署地震測報中心",
			"zh-CN": "交通部中央气象署地震测报中心",
		},
	}
}

func (c *CWA_SC) GetEvents(latitude, longitude float64) ([]Event, error) {
	// Get CWA HTML response
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.POST(
		"https://scweb.cwa.gov.tw/zh-tw/earthquake/ajaxhandler",
		c.getRequestBody(100),
		"application/x-www-form-urlencoded",
		10*time.Second, time.Second, 3, false,
		// Query CWA IP from custom encrypted DNS servers
		createCustomTransport(c.resolvers, ""),
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse CWA JSON response
	var dataMap map[string]any
	err = json.Unmarshal(res, &dataMap)
	if err != nil {
		return nil, err
	}

	dataMapEvents, ok := dataMap["data"].([]any)
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	var resultArr []Event
	for _, event := range dataMapEvents {
		eventData := event.([]any)
		if len(eventData) < 10 {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Depth:     string2Float(eventData[4].(string)),
			Event:     eventData[0].(string),
			Region:    eventData[5].(string),
			Latitude:  string2Float(eventData[8].(string)),
			Longitude: string2Float(eventData[7].(string)),
			Magnitude: c.getMagnitude(eventData[3].(string)),
			Timestamp: c.getTimestamp(eventData[2].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *CWA_SC) getTimestamp(textValue string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", textValue)
	return t.Add(-8 * time.Hour).UnixMilli()
}

func (c *CWA_SC) getMagnitude(textValue string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("ML"), Value: string2Float(textValue)}}
}
