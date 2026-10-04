package metadata_test

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/metadata"
)

func TestEmbeddedMetadataRendering(t *testing.T) {
	t.Parallel()
	options := metadata.Options{
		StartTime:  time.Date(2026, 1, 2, 3, 4, 5, 6, time.FixedZone("UTC+8", 8*60*60)),
		SampleRate: 100, Latitude: 25.1, Longitude: 121.5, Elevation: 12.25,
		NetworkCode: "TW", StationCode: "TEST", LocationCode: "00",
		ChannelCodes: []string{"HHZ", "HHE"}, StationPlace: "Taipei",
		StationCountry: "Taiwan", StationAffiliation: "SensePlex", StationDescription: "Test & Verify",
	}
	render, err := metadata.New("E-C111G", options)
	if err != nil {
		t.Fatal(err)
	}
	for name, document := range map[string]string{"SeisComP": render.SeisComP(), "StationXML": render.StationXML()} {
		if strings.Contains(document, "{{") {
			t.Errorf("%s contains an unrendered template value", name)
		}
		for _, text := range []string{"HHZ", "HHE", "CH3", "CH4", "CH5", "CH6", "TW", "TEST", "Taipei", "Test &amp; Verify"} {
			if !strings.Contains(document, text) {
				t.Errorf("%s is missing %q", name, text)
			}
		}
		var value struct{ XMLName xml.Name }
		if err := xml.Unmarshal([]byte(document), &value); err != nil {
			t.Errorf("%s is invalid XML: %v", name, err)
		}
	}
	if !strings.Contains(render.StationXML(), "2026-01-01T19:04:05.000000006Z") {
		t.Fatal("StationXML did not render StartTime in UTC")
	}
}

func TestMetadataModelsAndMissingModel(t *testing.T) {
	t.Parallel()
	for _, model := range []string{"E-C111G", "E-C121G", "E-C131G", "E-D001"} {
		render, err := metadata.New(model, metadata.Options{SampleRate: 100})
		if err != nil || render.SeisComP() == "" || render.StationXML() == "" {
			t.Errorf("metadata.New(%q) returned empty data or error: %v", model, err)
		}
	}
	if render, err := metadata.New("missing-model", metadata.Options{}); err == nil || render != nil {
		t.Fatalf("missing model = %#v, %v, want error", render, err)
	}
}
