package metadata

import (
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/bclswl0827/obsgo/inventory/stationxml"
)

const (
	UNIT_VELOCITY     = "m/s"
	UNIT_ACCELERATION = "m/s2"
)

// Sensitivity is the overall instrument sensitivity of one channel, in counts per SI unit.
type Sensitivity struct {
	Gain      float64
	Unit      string  // UNIT_VELOCITY or UNIT_ACCELERATION
	Frequency float64 // frequency (Hz) at which Gain is valid, 0 if unknown
	// Poles are the analog (Laplace, rad/s) poles of the sensor response, used
	// to correct the low-frequency roll-off of geophones. Empty when the
	// response is flat or not described in the Laplace domain.
	Poles []complex128
}

// NormalizeUnit maps StationXML input unit names to UNIT_VELOCITY or UNIT_ACCELERATION.
// Unsupported units yield an empty string.
func NormalizeUnit(name string) string {
	unit := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", ""))
	switch unit {
	case "m/s", "m*s**-1", "m/sec":
		return UNIT_VELOCITY
	case "m/s2", "m/s**2", "m/s^2", "m*s**-2", "m/sec**2":
		return UNIT_ACCELERATION
	}
	return ""
}

// ParseSensitivities returns the sensitivity of every channel in a StationXML document,
// keyed by channel code. Channels without a usable sensitivity are omitted.
func ParseSensitivities(document string) (map[string]Sensitivity, error) {
	var doc stationxml.Document
	if err := xml.NewDecoder(strings.NewReader(document)).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode StationXML: %w", err)
	}
	result := make(map[string]Sensitivity)
	for _, network := range doc.Networks {
		for _, station := range network.Stations {
			for _, channel := range station.Channels {
				if sensitivity, ok := channelSensitivity(channel); ok {
					result[channel.Code] = sensitivity
				}
			}
		}
	}
	return result, nil
}

func channelSensitivity(channel stationxml.Channel) (Sensitivity, bool) {
	if channel.Response == nil || channel.Response.Sensitivity == nil {
		return Sensitivity{}, false
	}
	gain := channel.Response.Sensitivity.Value
	unit := NormalizeUnit(channel.Response.Sensitivity.Input.Name)
	if unit == "" || gain <= 0 || math.IsNaN(gain) || math.IsInf(gain, 0) {
		return Sensitivity{}, false
	}
	return Sensitivity{
		Gain:      gain,
		Unit:      unit,
		Frequency: float64(channel.Response.Sensitivity.Frequency),
		Poles:     analogPoles(channel.Response.Stages),
	}, true
}

// analogPoles collects the poles of every Laplace pole-zero stage, in rad/s.
func analogPoles(stages []stationxml.Stage) []complex128 {
	var poles []complex128
	for _, stage := range stages {
		pz := stage.PolesZeros
		if pz == nil {
			continue
		}
		scale := 0.0
		switch kind := strings.ToUpper(pz.TransferFunctionType); {
		case strings.Contains(kind, "RADIANS"):
			scale = 1
		case strings.Contains(kind, "HERTZ"):
			scale = 2 * math.Pi
		default:
			continue // digital stages have no analog roll-off to correct
		}
		for _, pole := range pz.Poles {
			poles = append(poles, complex(float64(pole.Real)*scale, float64(pole.Imaginary)*scale))
		}
	}
	return poles
}

// Sensitivities renders the model's StationXML template and returns the channel
// sensitivities in template order: index 0 is ChannelCode1, and so on.
// Missing or unusable entries are nil.
func Sensitivities(model string) ([]*Sensitivity, error) {
	render, err := New(model, Options{})
	if err != nil {
		return nil, err
	}
	if render.channels == 0 {
		return nil, errors.New("metadata template declares no channels")
	}
	codes := make([]string, render.channels)
	for i := range codes {
		codes[i] = fmt.Sprintf("CH%d", i+1)
	}
	render.options.ChannelCodes = codes
	byCode, err := ParseSensitivities(render.StationXML())
	if err != nil {
		return nil, err
	}
	result := make([]*Sensitivity, len(codes))
	for i, code := range codes {
		if sensitivity, ok := byCode[code]; ok {
			result[i] = &sensitivity
		}
	}
	return result, nil
}
