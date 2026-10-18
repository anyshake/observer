package fdsnws

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/metadata"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/bclswl0827/obsgo/inventory/stationxml"
)

type metadataHardware struct {
	*testsupport.Hardware
	render *metadata.Render
}

func (h metadataHardware) GetMetadata(string, string, string, string, string, string, string, bool) (*metadata.Render, error) {
	return h.render, nil
}

func TestStationResponseFromMetadata(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		model            string
		channels         int
		velocityGain     float64
		dataloggerGain   float64
		accelerationGain float64
	}{
		{"E-C111G", 6, 85899345880, 858993458.8, 1670.7030433430377},
		{"E-C121G", 6, 85899345880, 858993458.8, 1670.7030433430377},
		{"E-C131G", 6, 85899345880, 858993458.8, 26104.735052234966},
		{"E-D001", 3, 343597383520, 3435973835.2, 0},
	} {
		t.Run(tc.model, func(t *testing.T) {
			render, err := metadata.New(tc.model, metadata.Options{
				StartTime: now, SampleRate: 100,
				NetworkCode: "XX", StationCode: "TEST", LocationCode: "00",
				ChannelCodes: []string{"EHZ", "EHN", "EHE", "ENZ", "ENN", "ENE"},
			})
			if err != nil {
				t.Fatal(err)
			}
			server, err := fdsnwsgo.NewFDSNWS(&provider{
				handler: handler,
				device:  metadataHardware{Hardware: testsupport.NewHardware(), render: render},
				now:     func() time.Time { return now },
			}, fdsnwsgo.WithModules(fdsnwsgo.ModuleStation))
			if err != nil {
				t.Fatal(err)
			}
			result := server.Handle(context.Background(), fdsnwsgo.Request{
				Method: http.MethodGet, Path: "/fdsnws/station/1/query",
				Query: url.Values{"level": {"response"}},
			})
			if result.StatusCode != http.StatusOK {
				t.Fatalf("response status = %d: %s", result.StatusCode, result.Body)
			}
			document, err := stationxml.Parse(result.Body)
			if err != nil {
				t.Fatal(err)
			}
			channels := document.Networks[0].Stations[0].Channels
			if len(channels) != tc.channels {
				t.Fatalf("channels = %d; want %d", len(channels), tc.channels)
			}
			for i, channel := range channels {
				response := channel.Response
				if response == nil || len(response.Stages) != 2 {
					t.Fatalf("channel %s lost response stages: %+v", channel.Code, response)
				}
				first, last := response.Stages[0], response.Stages[1]
				if first.PolesZeros == nil || first.Gain == nil || last.Coefficients == nil || last.Decimation == nil || last.Gain == nil {
					t.Fatalf("channel %s lost filter, gain or decimation data", channel.Code)
				}
				if first.Number != 1 || last.Number != 2 || last.Decimation.InputSampleRate != 100 || last.Decimation.Factor != 1 || last.Gain.Frequency != 0 {
					t.Fatalf("channel %s changed stage data: numbers=%d,%d inputSampleRate=%g factor=%d gainFrequency=%g; want 1,2 100 1 0",
						channel.Code, first.Number, last.Number, last.Decimation.InputSampleRate, last.Decimation.Factor, last.Gain.Frequency)
				}
				if i < 3 {
					sensitivity := response.Sensitivity
					if sensitivity == nil || sensitivity.Value != tc.velocityGain || sensitivity.Frequency != 4.5 || sensitivity.Input.Name != "m/s" || sensitivity.Output.Name != "COUNTS" {
						t.Fatalf("channel %s sensitivity = %+v", channel.Code, sensitivity)
					}
					if first.Gain.Value != 100 || first.Gain.Frequency != 4.5 || last.Gain.Value != tc.dataloggerGain || len(first.PolesZeros.Poles) == 0 || len(first.PolesZeros.Zeros) == 0 {
						t.Fatalf("channel %s changed velocity response", channel.Code)
					}
				} else if response.Sensitivity != nil || first.Gain.Value != tc.accelerationGain || first.Gain.Frequency != 1 || last.Gain.Value != 1 {
					t.Fatalf("channel %s invented sensitivity frequency or changed acceleration gains: %+v", channel.Code, response)
				}
			}
		})
	}
}

func TestMissingMetadataResponseFieldsAreOmitted(t *testing.T) {
	for _, tc := range []struct {
		name        string
		response    string
		sensitivity bool
		gain        bool
	}{
		{name: "no response"},
		{name: "empty response", response: "<Response/>"},
		{name: "missing sensitivity frequency", response: "<Response><InstrumentSensitivity><Value>2</Value></InstrumentSensitivity></Response>"},
		{name: "unknown sensitivity frequency", response: "<Response><InstrumentSensitivity><Value>2</Value><Frequency>None</Frequency></InstrumentSensitivity></Response>"},
		{name: "missing stage gain value", response: "<Response><Stage number=\"1\"><StageGain><Frequency>1</Frequency></StageGain></Stage></Response>"},
		{name: "unknown stage gain frequency", response: "<Response><Stage number=\"1\"><StageGain><Value>2</Value><Frequency>None</Frequency></StageGain></Stage></Response>"},
		{name: "missing sensitivity units", response: "<Response><InstrumentSensitivity><Value>2</Value><Frequency>1</Frequency></InstrumentSensitivity></Response>"},
		{
			name: "explicit zero values", sensitivity: true, gain: true,
			response: `<Response>
  <InstrumentSensitivity><Value>0</Value><Frequency>0</Frequency><InputUnits><Name>m/s</Name></InputUnits><OutputUnits><Name>COUNTS</Name></OutputUnits></InstrumentSensitivity>
  <Stage number="1"><StageGain><Value>0</Value><Frequency>0</Frequency></StageGain></Stage>
</Response>`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := `<FDSNStationXML><Network code="XX"><Station code="TEST"><Channel code="EHZ">` + tc.response + `</Channel></Station></Network></FDSNStationXML>`
			document, err := parseStationMetadata([]byte(data))
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := stationxml.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			body := string(encoded)
			if strings.Contains(body, "<InstrumentSensitivity>") != tc.sensitivity || strings.Contains(body, "<StageGain>") != tc.gain ||
				strings.Contains(body, "<Response>") != (tc.sensitivity || tc.gain) {
				t.Fatalf("missing response data was invented or known data was lost: %s", body)
			}
		})
	}
}
