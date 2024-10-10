package seisevent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"golang.org/x/exp/rand"
)

const CWA_API_ID = "cwa_api"

type CWA_API struct {
	cache cache.AnyCache
}

// Magic function that bypasses the Great Firewall of China
func (c *CWA_API) createGfwBypasser(customAddrs []string) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Choose a random address from the list
			customAddr := customAddrs[rand.Intn(len(customAddrs))]
			return (&net.Dialer{}).DialContext(ctx, network, customAddr)
		},
	}
}

func (c *CWA_API) getRequestBody(limit int) string {
	return fmt.Sprintf(`draw=1&columns[0][data]=0&columns[0][name]=EventNo&columns[0][searchable]=false&columns[0][orderable]=true&columns[0][search][value]=&columns[0][search][regex]=false&columns[1][data]=1&columns[1][name]=MaxIntensity&columns[1][searchable]=true&columns[1][orderable]=true&columns[1][search][value]=&columns[1][search][regex]=false&columns[2][data]=2&columns[2][name]=OriginTime&columns[2][searchable]=true&columns[2][orderable]=true&columns[2][search][value]=&columns[2][search][regex]=false&columns[3][data]=3&columns[3][name]=MagnitudeValue&columns[3][searchable]=true&columns[3][orderable]=true&columns[3][search][value]=&columns[3][search][regex]=false&columns[4][data]=4&columns[4][name]=Depth&columns[4][searchable]=true&columns[4][orderable]=true&columns[4][search][value]=&columns[4][search][regex]=false&columns[5][data]=5&columns[5][name]=Description&columns[5][searchable]=true&columns[5][orderable]=true&columns[5][search][value]=&columns[5][search][regex]=false&columns[6][data]=6&columns[6][name]=Description&columns[6][searchable]=true&columns[6][orderable]=true&columns[6][search][value]=&columns[6][search][regex]=false&order[0][column]=2&order[0][dir]=desc&start=0&length=%d&search[value]=&search[regex]=false&Search=&txtSDate=&txtEDate=&txtSscale=&txtEscale=&txtSdepth=&txtEdepth=&txtLonS=&txtLonE=&txtLatS=&txtLatE=&ddlCity=&ddlCitySta=&txtIntensityB=&txtIntensityE=&txtLon=&txtLat=&txtKM=&ddlStationName=&cblEventNo=&txtSDatePWS=&txtEDatePWS=&txtSscalePWS=&txtEscalePWS=&ddlMark=`, limit)
}

func (c *CWA_API) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CWA_API_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Central Weather Administration (API)",
			"zh-TW": "交通部中央氣象署（API）",
			"zh-CN": "交通部中央气象署（API）",
		},
	}
}

func (c *CWA_API) GetEvents(latitude, longitude float64) ([]Event, error) {
	// Get CWA HTML response
	if !c.cache.Valid() {
		res, err := request.POST(
			"https://scweb.cwa.gov.tw/zh-tw/earthquake/ajaxhandler",
			c.getRequestBody(100),
			"application/x-www-form-urlencoded; charset=UTF-8",
			10*time.Second, time.Second, 3, false,
			// HiNet CDN IP addresses
			c.createGfwBypasser([]string{
				"168.95.245.1:443", "168.95.245.2:443", "168.95.245.3:443", "168.95.245.4:443",
				"168.95.246.1:443", "168.95.246.2:443", "168.95.246.3:443", "168.95.246.4:443",
			}),
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse CWA JSON response
	var dataMap map[string]any
	err := json.Unmarshal(c.cache.Get().([]byte), &dataMap)
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
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (c *CWA_API) getTimestamp(textValue string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", textValue)
	return t.Add(-8 * time.Hour).UnixMilli()
}

func (c *CWA_API) getMagnitude(textValue string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("ML"), Value: string2Float(textValue)}}
}
