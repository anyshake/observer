package seisevent

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type stubRT struct {
	calls atomic.Int32
	mu    sync.Mutex
	code  int
	body  string
	byURL func(*http.Request) string
}

func (s *stubRT) RoundTrip(r *http.Request) (*http.Response, error) {
	s.calls.Add(1)
	s.mu.Lock()
	code, body, byURL := s.code, s.body, s.byURL
	s.mu.Unlock()
	if byURL != nil {
		body = byURL(r)
	}
	if code == 0 {
		code = http.StatusOK
	}
	return &http.Response{
		StatusCode:    code,
		Status:        fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Body:          io.NopCloser(strings.NewReader(body)),
		Header:        make(http.Header),
		ContentLength: int64(len(body)),
		Request:       r,
	}, nil
}

func (s *stubRT) use(code int, body string, byURL func(*http.Request) string) {
	s.mu.Lock()
	s.code, s.body, s.byURL = code, body, byURL
	s.mu.Unlock()
}

type sourceFixture struct {
	id           string
	valid        string
	validFor     func(*http.Request) string
	malformed    string
	malformedErr bool
	lat          float64
	lon          float64
	depth        float64
	region       string
	verified     bool
	magType      MagnitudeType
	mag          float64
}

func fdsnValidBody() string {
	return fdsnHeader + "event-1|2026-01-02T03:04:05|-12.5|123.75|10|agency|catalog|contributor|id|mww|5.6|agency|Near Coast\n"
}

func resetCache(src IDataSource) {
	fresh := cache.NewGeneric[[]Event](time.Minute)
	switch s := src.(type) {
	case *AFAD:
		s.cache = fresh
	case *BCSF:
		s.cache = fresh
	case *BGS:
		s.cache = fresh
	case *BMKG:
		s.cache = fresh
	case *CEA:
		s.cache = fresh
	case *CENC_APP:
		s.cache = fresh
	case *CENC_WEB:
		s.cache = fresh
	case *CENC_WOLFX:
		s.cache = fresh
	case *CWA_SC:
		s.cache = fresh
	case *DOST:
		s.cache = fresh
	case *EMSC:
		s.cache = fresh
	case *GA:
		s.cache = fresh
	case *GEONET:
		s.cache = fresh
	case *GFZ:
		s.cache = fresh
	case *HKO:
		s.cache = fresh
	case *ICL:
		s.cache = fresh
	case *INFP:
		s.cache = fresh
	case *INGV:
		s.cache = fresh
	case *JMA_OFFICIAL:
		s.cache = fresh
	case *JMA_P2PQUAKE:
		s.cache = fresh
	case *JMA_WOLFX:
		s.cache = fresh
	case *KMA:
		s.cache = fresh
	case *KNDC:
		s.cache = fresh
	case *KNMI:
		s.cache = fresh
	case *KRDAE:
		s.cache = fresh
	case *NCS:
		s.cache = fresh
	case *NRCAN:
		s.cache = fresh
	case *PALERT:
		s.cache = fresh
	case *SCEA:
		s.cache = fresh
	case *SED:
		s.cache = fresh
	case *SSN:
		s.cache = fresh
	case *TMD:
		s.cache = fresh
	case *USGS:
		s.cache = fresh
	case *USP:
		s.cache = fresh
	default:
		panic(fmt.Sprintf("resetCache: %T", src))
	}
}

func injectTransport(src IDataSource, rt http.RoundTripper) {
	switch s := src.(type) {
	case *BMKG:
		s.transport = rt
	case *NCS:
		s.transport = rt
	case *CWA_SC:
		s.transport = rt
	}
}

func findEvent(events []Event, lat, lon float64) (Event, bool) {
	for _, ev := range events {
		if math.Abs(ev.Latitude-lat) < 1e-6 && math.Abs(ev.Longitude-lon) < 1e-6 {
			return ev, true
		}
	}
	return Event{}, false
}

func assertMag(t *testing.T, ev Event, typ MagnitudeType, val float64) {
	t.Helper()
	for _, m := range ev.Magnitude {
		if m.Type == typ && math.Abs(m.Value-val) < 1e-6 {
			return
		}
	}
	t.Fatalf("magnitude = %+v, want %s %.1f", ev.Magnitude, typ, val)
}

func TestNewSources(t *testing.T) {
	sources, err := New(nil, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	wantIDs := []string{
		AFAD_ID, BCSF_ID, BGS_ID, BMKG_ID, CEA_ID, CENC_APP_ID, CENC_WEB_ID, CENC_WOLFX_ID,
		CWA_SC_ID, DOST_ID, EMSC_ID, GA_ID, GEONET_ID, GFZ_ID, HKO_ID, ICL_ID, INFP_ID, INGV_ID,
		JMA_OFFICIAL_ID, JMA_P2PQUAKE_ID, JMA_WOLFX_ID, KMA_ID, KNDC_ID, KNMI_ID, KRDAE_ID, NCS_ID,
		NRCAN_ID, PALERT_ID, SCEA_ID, SED_ID, SSN_ID, TMD_ID, USGS_ID, USP_ID,
	}
	if len(sources) != len(wantIDs) {
		t.Fatalf("New returned %d sources, want %d", len(sources), len(wantIDs))
	}
	for _, id := range wantIDs {
		src, ok := sources[id]
		if !ok {
			t.Errorf("missing source %s", id)
			continue
		}
		prop := src.GetProperty()
		if prop.ID != id || prop.Country == "" || prop.Default == "" || prop.Locales[prop.Default] == "" {
			t.Errorf("%s property = %+v", id, prop)
		}
	}
	if sources[BMKG_ID].(*BMKG).transport != nil || sources[NCS_ID].(*NCS).transport != nil || sources[CWA_SC_ID].(*CWA_SC).transport != nil {
		t.Fatal("injected transport must stay nil in production")
	}

	again, err := New(timesource.New(time.Now), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != len(sources) {
		t.Fatalf("second New returned %d sources", len(again))
	}
	if _, err := New(nil, time.Minute); err != nil {
		t.Fatal(err)
	}
}

func TestSeismicHelpers(t *testing.T) {
	raw, err := getGeoJsonData("twCounty2010")
	if err != nil || len(raw) == 0 {
		t.Fatalf("getGeoJsonData hit: %v len %d", err, len(raw))
	}
	if _, err := getGeoJsonData("missing"); err == nil {
		t.Fatal("getGeoJsonData miss returned nil error")
	}
	if string2Float("1.5") != 1.5 || string2Float("nope") != 0 {
		t.Fatal("string2Float")
	}
	if isMapKeysEmpty(map[string]any{"a": ""}, []string{"a"}) {
		t.Fatal("empty string key should fail isMapKeysEmpty")
	}
	if !isMapKeysEmpty(map[string]any{"a": "x", "n": 1.0}, []string{"a", "n"}) {
		t.Fatal("non-empty and non-string keys should pass isMapKeysEmpty")
	}
	if !isMapHasKeys(map[string]int{"a": 1}, []string{"a"}) || isMapHasKeys(map[string]int{"a": 1}, []string{"b"}) {
		t.Fatal("isMapHasKeys")
	}

	table, err := travel.NewAK135()
	if err != nil {
		t.Fatal(err)
	}
	near := getSeismicEstimation(table, 0, 0, 0, 1, 10)
	if near.P_Wave <= 0 || near.S_Wave <= 0 {
		t.Fatalf("short-distance estimation = %+v", near)
	}
	far := getSeismicEstimation(table, 0, 0, 0, 140, -5)
	if far.P_Wave <= 0 || far.S_Wave <= 0 {
		t.Fatalf("far estimation = %+v", far)
	}
	if far.P_Wave <= near.P_Wave {
		t.Fatalf("far P %.3f should exceed near P %.3f", far.P_Wave, near.P_Wave)
	}
}

func TestDataSourceGetEvents(t *testing.T) {
	rt := &stubRT{}
	restore := request.SetRoundTripperOverrideForTest(rt)
	t.Cleanup(restore)

	sources, err := New(nil, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	for _, src := range sources {
		injectTransport(src, rt)
	}

	for _, fx := range sourceFixtures {
		t.Run(fx.id, func(t *testing.T) {
			src, ok := sources[fx.id]
			if !ok {
				t.Fatalf("missing source %s", fx.id)
			}

			resetCache(src)
			rt.use(http.StatusOK, fx.malformed, nil)
			events, err := src.GetEvents(25, 121)
			if fx.malformedErr {
				if err == nil {
					t.Fatal("malformed body accepted")
				}
			} else if err != nil {
				t.Fatalf("malformed body: %v", err)
			} else if len(events) != 0 {
				t.Fatalf("malformed body produced %#v", events)
			}

			resetCache(src)
			rt.use(http.StatusServiceUnavailable, "down", nil)
			if _, err := src.GetEvents(25, 121); err == nil {
				t.Fatal("non-200 response accepted")
			}

			resetCache(src)
			before := rt.calls.Load()
			if fx.validFor != nil {
				rt.use(http.StatusOK, "", fx.validFor)
			} else {
				rt.use(http.StatusOK, fx.valid, nil)
			}
			events, err = src.GetEvents(25, 121)
			if err != nil {
				t.Fatalf("valid fixture: %v", err)
			}
			if rt.calls.Load() <= before {
				t.Fatal("transport was not called for a cache miss")
			}
			ev, ok := findEvent(events, fx.lat, fx.lon)
			if !ok {
				t.Fatalf("no event at %.4f, %.4f in %#v", fx.lat, fx.lon, events)
			}
			if ev.Region != fx.region {
				t.Fatalf("region = %q, want %q", ev.Region, fx.region)
			}
			if ev.Verfied != fx.verified {
				t.Fatalf("verified = %v, want %v", ev.Verfied, fx.verified)
			}
			if math.Abs(ev.Depth-fx.depth) > 1e-6 {
				t.Fatalf("depth = %v, want %v", ev.Depth, fx.depth)
			}
			assertMag(t, ev, fx.magType, fx.mag)
			if ev.Distance <= 0 {
				t.Fatalf("distance = %v", ev.Distance)
			}
			if ev.Estimation.P_Wave < 0 || ev.Estimation.S_Wave < 0 {
				t.Fatalf("estimation = %+v", ev.Estimation)
			}

			calls := rt.calls.Load()
			again, err := src.GetEvents(25, 121)
			if err != nil {
				t.Fatal(err)
			}
			if rt.calls.Load() != calls {
				t.Fatalf("cache miss: calls %d -> %d", calls, rt.calls.Load())
			}
			if len(again) != len(events) {
				t.Fatalf("cached len %d, want %d", len(again), len(events))
			}
			cached, ok := findEvent(again, fx.lat, fx.lon)
			if !ok || cached.Region != ev.Region {
				t.Fatalf("cached event = %+v", cached)
			}
		})
	}
}

func bgsItem(string) string {
	return `<rss><channel><item>
<description>Origin date/time: Sun, 10 Aug 2025 16:53:46 ; Location: WESTERN TURKEY ; Lat/long: 39.312,28.069 ; Depth: 10 km ; Magnitude: 6.1</description>
<lat>39.312</lat><lon>28.069</lon>
</item><item>
<description>not a timestamp ; nope ; x ; depthonly ; magonly</description>
<lat>1</lat><lon>2</lon>
</item></channel></rss>`
}

var sourceFixtures = []sourceFixture{
	{
		id: AFAD_ID, malformed: `{`, malformedErr: true,
		valid: `[{"eventID":"E1","location":"VAN","latitude":"38.5","longitude":"43.4","depth":"10","type":"ML","magnitude":"4.5","date":"2026-01-02T03:04:05"},{"eventID":""}]`,
		lat:   38.5, lon: 43.4, depth: 10, region: "VAN", verified: true, magType: "Ml", mag: 4.5,
	},
	{
		id: BCSF_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: BGS_ID, malformed: `<`, malformedErr: false,
		validFor: func(r *http.Request) string { return bgsItem(r.URL.Path) },
		lat:      39.312, lon: 28.069, depth: 10, region: "WESTERN TURKEY", verified: true, magType: "M", mag: 6.1,
	},
	{
		id: BMKG_ID, malformed: `<`, malformedErr: false,
		valid: `<gempa><info><eventid>E1</eventid><date>02-01-26</date><time>03:04:05 WIB</time><magnitude>4.5</magnitude><depth>10 Km</depth><area>Jawa</area><coordinates>106.8,-6.2</coordinates></info><info><eventid>skip</eventid></info></gempa>`,
		lat:   -6.2, lon: 106.8, depth: 10, region: "Jawa", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: CEA_ID, malformed: `<html></html>`, malformedErr: false,
		valid: `<table><tr class="arial11bleu"><td>02/01/2026</td><td>03:04:05</td><td>25.1,121.5</td><td>Taiwan</td><td>ML=4.2</td></tr><tr class="arial11bleu"><td>bad</td><td>bad</td><td>nocoord</td><td>x</td><td>nomag</td></tr></table>`,
		lat:   25.1, lon: 121.5, depth: -1, region: "Taiwan", verified: false, magType: "Ml", mag: 4.2,
	},
	{
		id: CENC_APP_ID, malformed: `{`, malformedErr: true,
		valid: `{"result":"OK","values":[{"time":1735689845000,"longitude":121.5,"latitude":25.1,"depth":10000,"eqid":"E1","loc_name":"四川","mag":4.5,"eq_type":"M"},{"time":1,"longitude":1,"latitude":2,"depth":1000,"eqid":"E2","loc_name":"中国四川","mag":3,"eq_type":"A"},{"time":2,"longitude":3,"latitude":4,"depth":1000,"eqid":"E3","loc_name":"台湾省花莲","mag":3,"eq_type":"M"},{"eqid":""}]}`,
		lat:   25.1, lon: 121.5, depth: 10, region: "四川", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: CENC_WEB_ID, malformed: `{`, malformedErr: true,
		valid: `{"features":[{"id":"E1","properties":{"震级":"4.5","参考位置":"台湾","发震时刻":"2026-01-02 03:04:05","深度（千米）":"10"},"geometry":{"coordinates":[121.5,25.1]}},{"id":"E2","properties":{"震级":"3","参考位置":"深度数字","发震时刻":"2026-01-02 04:05:06","深度（千米）":8},"geometry":{"coordinates":[100,20]}},{"id":"E3","properties":{"震级":"3","参考位置":"深度其它","发震时刻":"2026-01-02 05:06:07","深度（千米）":true},"geometry":{"coordinates":[101,21]}},{"id":"E4","properties":{"震级":"3","参考位置":"坏时间","发震时刻":"not-a-time","深度（千米）":"1"},"geometry":{"coordinates":[102,22]}},{"id":"nog","properties":{"震级":"1","参考位置":"x","发震时刻":"2026-01-02 03:04:05","深度（千米）":"1"},"geometry":{}},{"id":"skip"}]}`,
		lat:   25.1, lon: 121.5, depth: 10, region: "台湾", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: CENC_WOLFX_ID, malformed: `{`, malformedErr: true,
		valid: `{"No1":{"EventID":"E1","type":"reviewed","time":"2026-01-02 03:04:05","location":"台湾","magnitude":"4.5","depth":"10","latitude":"25.1","longitude":"121.5"},"count":1,"skip":{"EventID":""},"bad":{"EventID":"E2","type":"automatic","time":"bad","location":"x","magnitude":"1","depth":"1","latitude":"1","longitude":"2"}}`,
		lat:   25.1, lon: 121.5, depth: 10, region: "台湾", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: CWA_SC_ID, malformed: `{`, malformedErr: true,
		valid: `{"data":[["E1","3","2026-01-02 03:04:05","4.5","10","花蓮","x","121.5","25.1","z"],["short"]]}`,
		lat:   25.1, lon: 121.5, depth: 10, region: "花蓮", verified: true, magType: "Ml", mag: 4.5,
	},
	{
		id: DOST_ID, malformed: `<html></html>`, malformedErr: false,
		valid: `<table><caption>Latitude Longitude Depth Mag Location</caption><tr><td>2 January 2026 - 3:04 PM</td><td>14.5</td><td>121.0</td><td>10</td><td>4.5</td><td>Luzon</td></tr></table>`,
		lat:   14.5, lon: 121, depth: 10, region: "Luzon", verified: true, magType: "MS", mag: 4.5,
	},
	{
		id: EMSC_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: GA_ID, malformed: `{`, malformedErr: true,
		valid: `{"features":[{"properties":{"depth":10,"preferred_magnitude_type":"ML","preferred_magnitude":4.5,"event_id":"E1","latitude":-25.1,"longitude":133.5,"description":"Australia","evaluation_status":"confirmed","origin_time":"2026-01-02T03:04:05.000Z"},"geometry":{}},{"properties":{"depth":5,"preferred_magnitude_type":"ML","preferred_magnitude":2,"event_id":"E2","latitude":1,"longitude":2,"description":"prelim","evaluation_status":"preliminary","origin_time":"bad"},"geometry":{}},{"properties":{"depth":1},"geometry":{}},{"id":"skip"}]}`,
		lat:   -25.1, lon: 133.5, depth: 10, region: "Australia", verified: true, magType: "Ml", mag: 4.5,
	},
	{
		id: GEONET_ID, malformed: `{`, malformedErr: true,
		valid: `{"features":[{"properties":{"publicID":"E1","time":"2026-01-02T03:04:05.000Z","depth":10,"magnitude":4.5,"locality":"Canterbury"},"geometry":{"coordinates":[172.6,-43.5]}},{"properties":{"publicID":"E2","time":"2026-01-02T03:04:05.000Z","depth":1,"magnitude":1,"locality":"x"},"geometry":{"coordinates":[1,2,3]}},{"properties":{"publicID":"x"},"geometry":{"coordinates":[1,2]}},{"id":"skip"}]}`,
		lat:   -43.5, lon: 172.6, depth: 10, region: "Canterbury", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: GFZ_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: HKO_ID, malformed: `<`, malformedErr: true,
		valid: `<Earthquake><EventGroup>` +
			`<Event><EventId>E1</EventId><Verify>Y</Verify><HKTDate>20260102</HKTDate><HKTTime>0304</HKTTime><City>Taipei</City><Region>Taiwan</Region><Lat>25.1</Lat><Lon>121.5</Lon><Mag>4.5</Mag><Depth>10</Depth></Event>` +
			`<Event><EventId>skip</EventId></Event>` +
			`<Event><EventId>E2</EventId><Verify>N</Verify><HKTDate>bad</HKTDate><HKTTime>0304</HKTTime><City>X</City><Region>Y</Region><Lat>1</Lat><Lon>2</Lon><Mag>1</Mag><Depth>1</Depth></Event>` +
			`</EventGroup></Earthquake>`,
		lat: 25.1, lon: 121.5, depth: 10, region: "Taipei - Taiwan", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: ICL_ID, malformed: `{`, malformedErr: true,
		valid: `{"code":0,"data":[{"eventId":"E1","latitude":30.5,"longitude":104.0,"depth":10,"epicenter":"四川","startAt":1735689845000,"magnitude":4.5},{"eventId":"","epicenter":""}]}`,
		lat:   30.5, lon: 104, depth: 10, region: "四川", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: INFP_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: INGV_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: JMA_OFFICIAL_ID, malformed: `{`, malformedErr: true,
		valid: `[{"eid":"20260102030412","anm":"福島県","mag":"4.5","cod":"+37.5+141.2-10000/","at":"2026-01-02T03:04:00+09:00"},{"eid":"ab","anm":"短","mag":"1","cod":"37.5","at":"bad"},{"eid":""}]`,
		lat:   37.5, lon: 141.2, depth: 10, region: "福島県", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: JMA_P2PQUAKE_ID, malformed: `{`, malformedErr: true,
		valid: `[{"id":"E1","earthquake":{"time":"2026/01/02 03:04:05","hypocenter":{"name":"福島","latitude":37.5,"longitude":141.2,"depth":10,"magnitude":4.5}}},{"id":"E3","earthquake":{"hypocenter":{"name":"x","latitude":1,"longitude":2,"depth":3,"magnitude":4}}},{"id":"E4","earthquake":{"time":"2026/01/02 03:04:05","hypocenter":{"latitude":1}}},{"id":"E2","earthquake":{"time":"bad","hypocenter":{"name":"x","latitude":1,"longitude":2,"depth":3,"magnitude":4}}}]`,
		lat:   37.5, lon: 141.2, depth: 10, region: "福島", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: JMA_WOLFX_ID, malformed: `{`, malformedErr: true,
		valid: `{"No1":{"EventID":"E1","time_full":"2026/01/02 03:04:05","location":"福島","magnitude":"4.5","depth":"10km","latitude":"37.5","longitude":"141.2"},"count":1,"skip":{"EventID":""},"bad":{"EventID":"E2","time_full":"bad","location":"x","magnitude":"1","depth":"1km","latitude":"1","longitude":"2"}}`,
		lat:   37.5, lon: 141.2, depth: 10, region: "福島", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: KMA_ID, malformed: `<html></html>`, malformedErr: false,
		valid: `<div id="excel_body"><table><tbody><tr><td></td><td>2026/01/02 03:04:05</td><td>4.5</td><td>10</td><td></td><td>37.5 N</td><td>129.1 E</td><td>경주</td></tr><tr><td></td><td>2026/01/02 04:05:06</td><td>3</td><td>5</td><td></td><td>12.0 S</td><td>70.0 W</td><td>south</td></tr><tr><td></td><td>bad</td><td>x</td><td>x</td><td></td><td>12</td><td>70</td><td>plain</td></tr></tbody></table></div>`,
		lat:   37.5, lon: 129.1, depth: 10, region: "경주", verified: true, magType: "Ml", mag: 4.5,
	},
	{
		id: KNDC_ID, malformed: `{`, malformedErr: true,
		valid: `[{"id":"E1","epochtime":"1735689845","lat":"43.2","lon":"76.9","depth":10,"mb":4.5,"gregion":"Almaty","sregion":"Kazakhstan"},{"id":""}]`,
		lat:   43.2, lon: 76.9, depth: 10, region: "Almaty (Kazakhstan)", verified: true, magType: "Mb", mag: 4.5,
	},
	{
		id: KNMI_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: KRDAE_ID, malformed: `<html></html>`, malformedErr: false,
		valid: "<pre>\n--------------\n2026.01.02 03:04:05 37.50 28.10 10.0 2.1 3.2 4.3 TURKEY\n</pre>",
		lat:   37.5, lon: 28.1, depth: 10, region: "TURKEY", verified: true, magType: "Mw", mag: 4.3,
	},
	{
		id: NCS_ID, malformed: `<html></html>`, malformedErr: false,
		valid: `<div id="sidebar-wrapper">` +
			`<div class="event_list"></div>` +
			`<div class="event_list" data-json='{'></div>` +
			`<div class="event_list" data-json='{"event_name":"x"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E","event_name":"N"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E","event_name":"N","origin_time":"t"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E","event_name":"N","origin_time":"t","lat_long":"1,2"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E","event_name":"N","origin_time":"t","lat_long":"1,2","magnitude_depth":"M:1, Depth:1 km"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E2","event_name":"Bad","origin_time":"bad","lat_long":"28.6","magnitude_depth":"4.5","event_type":"Automatic"}'></div>` +
			`<div class="event_list" data-json='{"event_id":"E1","event_name":"Delhi","origin_time":"2026-01-02 03:04:05 IST","lat_long":"28.6, 77.2","magnitude_depth":"M:4.5, Depth:10 km","event_type":"Reviewed"}'></div>` +
			`</div>`,
		lat: 28.6, lon: 77.2, depth: 10, region: "Delhi", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: NRCAN_ID, malformed: `<`, malformedErr: true,
		valid: `<quakeml xmlns="http://quakeml.org/xmlns/bed/1.2"><eventParameters>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>25.1</value></latitude><longitude><value>121.5</value></longitude><depth><value>10</value></depth></origin>` +
			`<magnitude><mag><value>4.5</value></mag><type>mww</type></magnitude></event>` +
			`<event><type>quarry blast</type><typeCertainty>known</typeCertainty><description><text>Quarry</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>1</value></latitude><longitude><value>2</value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value>1.0</value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text lang="en">Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>1</value></latitude><longitude><value>2</value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value>1</value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value><x>1</x></value></time><latitude><value>1</value></latitude><longitude><value>2</value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value>1</value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value><x>1</x></value></latitude><longitude><value>2</value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value>1</value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>1</value></latitude><longitude><value><x>2</x></value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value>1</value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>1</value></latitude><longitude><value>2</value></longitude><depth><value><x>3</x></value></depth></origin>` +
			`<magnitude><mag><value>1</value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>1</value></latitude><longitude><value>2</value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value><x>1</x></value></mag><type>ml</type></magnitude></event>` +
			`<event><type>earthquake</type><typeCertainty>known</typeCertainty><description><text>Near Coast</text></description>` +
			`<origin><time><value>2026-01-02T03:04:05.000000Z</value></time><latitude><value>1</value></latitude><longitude><value>2</value></longitude><depth><value>3</value></depth></origin>` +
			`<magnitude><mag><value>1</value></mag><type lang="en">ml</type></magnitude></event>` +
			`</eventParameters></quakeml>`,
		lat: 25.1, lon: 121.5, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 4.5,
	},
	{
		id: PALERT_ID, malformed: `{`, malformedErr: true,
		valid: `{"data":{"eventList":[{"DateUTC":"2026-01-02T03:04:05","Depth":10,"Latitude":23.537,"Longitude":120.1305,"ML":4.5},{"DateUTC":"2026-01-02T04:05:06","Depth":10,"Latitude":0,"Longitude":0,"ML":2},{"DateUTC":"bad","Depth":1,"Latitude":1,"Longitude":1,"ML":1},{"DateUTC":""}]}}`,
		lat:   23.537, lon: 120.1305, depth: 10, region: "雲林縣", verified: true, magType: "Ml", mag: 4.5,
	},
	{
		id: SCEA_ID, malformed: `{`, malformedErr: true,
		valid: `{"code":0,"msg":"ok","data":[{"eventId":"E1","shockTime":1735689845000,"longitude":104.0,"latitude":30.5,"placeName":"四川","magnitude":4.5,"depth":10,"infoTypeName":"[正式]"},{"eventId":"","shockTime":1,"placeName":"x"},{"eventId":"E2","shockTime":1735689846000,"longitude":1,"latitude":2,"placeName":"云南","magnitude":3,"infoTypeName":"auto"}]}`,
		lat:   30.5, lon: 104, depth: 10, region: "四川", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: SED_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
	{
		id: SSN_ID, malformed: `<html></html>`, malformedErr: false,
		valid: `<div class="table-latest"><table>` +
			`<tr><td>1</td></tr>` +
			`<tr title="bad-time"><td>1</td><td><span>nope</span><span>nope</span></td><td>NoColon</td><td>1</td></tr>` +
			`<tr title="bad-coord"><td>1</td><td><span>2026-01-02</span><span>03:04:05</span></td><td>Place: onlylat</td><td>1</td></tr>` +
			`<tr title="ssn-1"><td>4.5</td><td><span>2026-01-02</span><span>03:04:05</span></td><td>Oaxaca: 16.5°, -95.1°</td><td>10 km</td></tr>` +
			`</table></div>`,
		lat: 16.5, lon: -95.1, depth: 10, region: "Oaxaca", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: TMD_ID, malformed: `<html></html>`, malformedErr: false,
		valid: `<table><tr class="tbis_leq1"><td>2026-01-02 03:04:05</td><td>4.5</td><td>13.5°N</td><td>100.5°E</td><td>10</td><td></td><td><span>กรุงเทพ<br/>Bangkok</span></td></tr>` +
			`<tr class="tbis_leq1"><td>bad</td><td>1</td><td>bad</td><td>bad</td><td>1</td><td></td><td>noregion</td></tr>` +
			`<tr class="tbis_leq2"><td>2026-01-02 04:05:06</td><td>3.2</td><td>14.0°S</td><td>101.0°W</td><td>20</td><td></td><td><span>เชียงใหม่<br/>Chiang Mai</span></td></tr>` +
			`<tr class="tbis_leq2"><td>also-bad</td><td>1</td><td>95°N</td><td>200°E</td><td>1</td><td></td><td>noregion</td></tr></table>`,
		lat: 13.5, lon: 100.5, depth: 10, region: "กรุงเทพ (Bangkok)", verified: true, magType: "M", mag: 4.5,
	},
	{
		id: USGS_ID, malformed: `{`, malformedErr: true,
		valid: `{"features":[{"id":"us1","properties":{"mag":4.5,"place":"Taiwan","time":1735689845000,"type":"earthquake","title":"M 4.5","status":"reviewed","magType":"ml"},"geometry":{"coordinates":[121.5,25.1,10]}},{"id":"q","properties":{"mag":1,"place":"x","time":1,"type":"quarry","title":"t","status":"automatic","magType":"ml"},"geometry":{"coordinates":[1,2,3]}},{"id":"short","properties":{"mag":1,"place":"x","time":1,"type":"earthquake","title":"t","status":"automatic","magType":"ml"},"geometry":{"coordinates":[1,2]}},{"id":"partial","properties":{"mag":1},"geometry":{"coordinates":[1,2,3]}},{"id":"skip"}]}`,
		lat:   25.1, lon: 121.5, depth: 10, region: "Taiwan", verified: true, magType: "Ml", mag: 4.5,
	},
	{
		id: USP_ID, malformed: "not-an-event", malformedErr: true,
		valid: fdsnValidBody(),
		lat:   -12.5, lon: 123.75, depth: 10, region: "Near Coast", verified: true, magType: "Mw", mag: 5.6,
	},
}

func TestParserEdgeBranches(t *testing.T) {
	if _, err := (&AFAD{}).getTimestamp("bad"); err == nil {
		t.Fatal("AFAD timestamp")
	}
	bgs := &BGS{}
	if _, err := bgs.getTimestamp("no comma here"); err == nil {
		t.Fatal("BGS timestamp shape")
	}
	if _, err := bgs.getTimestamp("Sun, not-a-date"); err == nil {
		t.Fatal("BGS timestamp parse")
	}
	if _, err := bgs.getDepth("nocolon"); err == nil {
		t.Fatal("BGS depth")
	}
	if _, err := bgs.getMagnitude("nocolon"); err == nil {
		t.Fatal("BGS magnitude")
	}
	bmkg := &BMKG{}
	if _, _, err := bmkg.getCoordinates("106.8"); err == nil {
		t.Fatal("BMKG coordinates")
	}
	if _, err := bmkg.getTimestamp("bad", "bad"); err == nil {
		t.Fatal("BMKG timestamp")
	}
	cea := &CEA{}
	if len(cea.getMagnitude("nomag")) != 0 || cea.getLatitude("25") != 0 || cea.getLongitude("121") != 0 {
		t.Fatal("CEA partial fields")
	}
	web := &CENC_WEB{}
	if _, err := web.getTimestamp("bad"); err == nil {
		t.Fatal("CENC web timestamp")
	}
	if web.getDepth(8.5) != 8.5 || web.getDepth(true) != -1 {
		t.Fatal("CENC web depth")
	}
	if _, err := (&CENC_WOLFX{}).getTimestamp("bad"); err == nil {
		t.Fatal("CENC wolfx timestamp")
	}
	if (&DOST{}).getTimestamp("bad") != 0 {
		t.Fatal("DOST timestamp")
	}
	if _, err := (&HKO{}).getTimestamp("bad"); err == nil {
		t.Fatal("HKO timestamp")
	}
	jma := &JMA_OFFICIAL{}
	if _, err := jma.getTimestamp("bad", "00"); err == nil {
		t.Fatal("JMA timestamp")
	}
	if jma.getDepth("37.5") != 0 || jma.getLatitude("37.5") != 0 || jma.getLongitude("37.5") != 0 {
		t.Fatal("JMA short code")
	}
	if _, err := (&JMA_P2PQUAKE{}).getTimestamp("bad"); err == nil {
		t.Fatal("P2P timestamp")
	}
	if _, err := (&JMA_WOLFX{}).getTimestamp("bad"); err == nil {
		t.Fatal("JMA wolfx timestamp")
	}
	kma := &KMA{}
	if kma.getLatitude("12.0 S") != -12 || kma.getLatitude("12") != 0 {
		t.Fatal("KMA latitude")
	}
	if kma.getLongitude("70.0 W") != -70 || kma.getLongitude("70") != 0 {
		t.Fatal("KMA longitude")
	}
	ncs := &NCS{}
	if ncs.getLatitude("28.6") != -1 || ncs.getLongitude("77.2") != -1 {
		t.Fatal("NCS coordinates")
	}
	if ncs.getDepth("M:4.5") != -1 || ncs.getDepth("M:4.5, 10km") != -1 {
		t.Fatal("NCS depth")
	}
	if len(ncs.getMagnitude("4.5")) != 0 || len(ncs.getMagnitude("4.5, Depth:10")) != 0 {
		t.Fatal("NCS magnitude")
	}
	if _, err := (&SSN{}).getTimestamp("bad"); err == nil {
		t.Fatal("SSN timestamp")
	}
	tmd := &TMD{}
	for _, lat := range []string{"bad", "abc°N", "13°X", "95°N"} {
		if _, err := tmd.getLatitude(lat); err == nil {
			t.Fatalf("TMD latitude %q", lat)
		}
	}
	for _, lon := range []string{"bad", "abc°E", "13°X", "200°E"} {
		if _, err := tmd.getLongitude(lon); err == nil {
			t.Fatalf("TMD longitude %q", lon)
		}
	}
	if _, err := tmd.getTimestamp("bad"); err == nil {
		t.Fatal("TMD timestamp")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<td>plain</td>`))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := tmd.getRegion(doc.Find("td")); err == nil {
		t.Fatal("TMD region")
	}

	palert := &PALERT{}
	if palert.pointToSegmentDistance(orb.Point{0, 0}, orb.Point{1, 1}, orb.Point{1, 1}) <= 0 {
		t.Fatal("degenerate segment")
	}
	poly := orb.Polygon{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}
	fc := geojson.NewFeatureCollection()
	inside := geojson.NewFeature(poly)
	inside.Properties["COUNTYNAME"] = "測試縣"
	outside := geojson.NewFeature(orb.MultiPolygon{poly})
	outside.Properties["COUNTYNAME"] = "外海縣"
	fc.Append(inside)
	if got := palert.getRegion(fc, 0.5, 0.5); got != "測試縣" {
		t.Fatalf("polygon region %q", got)
	}
	fc = geojson.NewFeatureCollection()
	fc.Append(outside)
	if got := palert.getRegion(fc, 0.5, 0.5); got != "外海縣" {
		t.Fatalf("multipolygon region %q", got)
	}
	if got := palert.getRegion(fc, 10, 10); !strings.Contains(got, "附近海域") {
		t.Fatalf("offshore region %q", got)
	}
	if palert.distanceToPolygon(orb.Point{5, 5}, poly) <= 0 {
		t.Fatal("polygon distance")
	}
	if palert.pointToSegmentDistance(orb.Point{0.5, -1}, orb.Point{0, 0}, orb.Point{1, 0}) <= 0 {
		t.Fatal("segment projection")
	}
}

func TestGetEventsEdgeBodies(t *testing.T) {
	rt := &stubRT{}
	restore := request.SetRoundTripperOverrideForTest(rt)
	t.Cleanup(restore)

	sources, err := New(nil, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	for _, src := range sources {
		injectTransport(src, rt)
	}

	cases := []struct {
		id      string
		body    string
		wantErr bool
	}{
		{AFAD_ID, `[{"eventID":"E1","location":"VAN","latitude":"38.5","longitude":"43.4","depth":"10","type":"ML","magnitude":"4.5","date":"bad"}]`, true},
		{AFAD_ID, `[]`, false},
		{CENC_APP_ID, `{"result":"NO"}`, true},
		{CENC_APP_ID, `{"result":"OK"}`, true},
		{CENC_WEB_ID, `{}`, true},
		{CENC_WEB_ID, `{"features":[]}`, false},
		{GA_ID, `{}`, true},
		{GA_ID, `{"features":[]}`, false},
		{GEONET_ID, `{}`, true},
		{GEONET_ID, `{"features":[]}`, false},
		{USGS_ID, `{}`, true},
		{USGS_ID, `{"features":[]}`, false},
		{ICL_ID, `{"code":1}`, true},
		{ICL_ID, `{"code":0}`, true},
		{SCEA_ID, `{"code":1,"msg":"nope","data":[]}`, true},
		{SCEA_ID, `{"code":0,"msg":"ok","data":[]}`, true},
		{PALERT_ID, `{}`, true},
		{PALERT_ID, `{"data":{}}`, true},
		{CWA_SC_ID, `{}`, true},
		{CWA_SC_ID, `{"data":[]}`, false},
		{HKO_ID, `<Earthquake><EventGroup><Other>x</Other></EventGroup></Earthquake>`, true},
		{NRCAN_ID, `<quakeml xmlns="http://quakeml.org/xmlns/bed/1.2"></quakeml>`, true},
		{NRCAN_ID, `<quakeml xmlns="http://quakeml.org/xmlns/bed/1.2"><eventParameters><event><type>earthquake</type></event></eventParameters></quakeml>`, true},
		{BMKG_ID, `<gempa><info><eventid>E1</eventid><date>02-01-26</date><time>03:04:05 WIB</time><magnitude>4.5</magnitude><depth>10 Km</depth><area>Jawa</area><coordinates>bad</coordinates></info></gempa>`, true},
		{BMKG_ID, `<gempa><info><eventid>E1</eventid><date>bad</date><time>bad</time><magnitude>4.5</magnitude><depth>10 Km</depth><area>Jawa</area><coordinates>106.8,-6.2</coordinates></info></gempa>`, true},
		{NRCAN_ID, `<other></other>`, true},
		{JMA_OFFICIAL_ID, `[]`, false},
		{JMA_P2PQUAKE_ID, `[]`, false},
		{KNDC_ID, `[]`, false},
	}
	for _, tc := range cases {
		t.Run(tc.id+"/"+tc.body[:min(12, len(tc.body))], func(t *testing.T) {
			src := sources[tc.id]
			resetCache(src)
			rt.use(http.StatusOK, tc.body, nil)
			events, err := src.GetEvents(25, 121)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(events) != 0 {
				t.Fatalf("got %#v", events)
			}
		})
	}
}
