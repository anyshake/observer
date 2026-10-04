package seisevent

import (
	"math"
	"reflect"
	"testing"
	"time"
)

const fdsnHeader = "#EventID|Time|Latitude|Longitude|Depth/km|Author|Catalog|Contributor|ContributorID|MagType|Magnitude|MagAuthor|EventLocationName\n"

func TestParseFdsnwsEvent(t *testing.T) {
	t.Parallel()
	data := fdsnHeader + "event-1|2026-01-02T03:04:05.123|-12.5|123.75|10|agency|catalog|contributor|id|mww|5.6|agency|Near Coast, Region\n" +
		"event-2|2026-01-02T04:05:06|1|2|3|agency|catalog|contributor|id|ml|2.1|agency|A \"quoted\" place\n"
	got, err := ParseFdsnwsEvent(data, "2006-01-02T15:04:05")
	if err != nil {
		t.Fatal(err)
	}
	want := []Event{
		{Event: "event-1", Timestamp: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC).UnixMilli(), Verfied: true, Latitude: -12.5, Longitude: 123.75, Depth: 10, Region: "Near Coast, Region", Magnitude: []Magnitude{{Type: "Mw", Value: 5.6}}},
		{Event: "event-2", Timestamp: time.Date(2026, 1, 2, 4, 5, 6, 0, time.UTC).UnixMilli(), Verfied: true, Latitude: 1, Longitude: 2, Depth: 3, Region: "A \"quoted\" place", Magnitude: []Magnitude{{Type: "Ml", Value: 2.1}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseFdsnwsEvent() = %#v, want %#v", got, want)
	}
}

func TestParseFdsnwsEventRejectsInvalidRecords(t *testing.T) {
	t.Parallel()
	for _, data := range []string{
		"", fdsnHeader,
		fdsnHeader + "event|bad-time|1|2|3|a|b|c|d|Ml|2|a|region\n",
		fdsnHeader + "event|2026-01-02T03:04:05\n",
		"#EventID|Time\nevent|2026-01-02T03:04:05\n",
	} {
		if events, err := ParseFdsnwsEvent(data, "2006-01-02T15:04:05"); err == nil || events != nil {
			t.Errorf("invalid FDSN data accepted: %q => %#v, %v", data, events, err)
		}
	}
}

func TestParseFdsnwsEventWithoutLocationName(t *testing.T) {
	t.Parallel()
	const header = "#EventID|Time|Latitude|Longitude|Depth/km|Author|Catalog|Contributor|ContributorID|MagType|Magnitude|MagAuthor|EventLocationName|EventType\n"
	for _, location := range []string{"", "   "} {
		data := header + "quakeml:apiv2.infp.ro/event/atlas24/2025253/10130032|2025-09-10T19:30:19.887|52.723800|159.314100|10.0|NIEP:rtMb||niep|quakeml:apiv2.infp.ro/event/atlas24/2025253/10130032|mb|5.25|dbevproc|" + location + "|earthquake\n"
		events, err := ParseFdsnwsEvent(data, "2006-01-02T15:04:05")
		if err != nil {
			t.Fatal(err)
		}
		if len(events) != 1 {
			t.Fatalf("got %d events, want 1", len(events))
		}
		if events[0].Region != "Latitude: 52.7238°, Longitude: 159.3141°" {
			t.Errorf("Region = %q, want epicenter coordinates", events[0].Region)
		}
	}
}

func TestMagnitudeNormalization(t *testing.T) {
	t.Parallel()
	for input, want := range map[string]MagnitudeType{"m": "M", "ml": "Ml", "MLv": "Ml", "ms": "MS", "mww": "Mw", "mb": "Mb", "md": "Md", "unknown": "UNKNOWN"} {
		if got := ParseMagnitude(input); got != want {
			t.Errorf("ParseMagnitude(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestEventDistanceAndSortOrder(t *testing.T) {
	t.Parallel()
	if got := getDistance(25, 25, 121, 121); got != 0 {
		t.Errorf("same-location distance = %v, want 0", got)
	}
	if got := getDistance(0, 0, 0, 1); math.IsNaN(got) || math.Abs(got-111.3195) > 0.0001 {
		t.Errorf("one equatorial degree = %v km, want 111.3195", got)
	}
	events := []Event{{Timestamp: 2}, {Timestamp: 1}, {Timestamp: 3}}
	got := sortSeismicEvents(events)
	if got[0].Timestamp != 3 || got[1].Timestamp != 2 || got[2].Timestamp != 1 {
		t.Fatalf("events are not newest first: %#v", got)
	}
}
