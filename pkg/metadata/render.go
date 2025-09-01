package metadata

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

type Render struct {
	options            *Options
	channels           int
	templateSeisComP   *template.Template
	templateStationXML *template.Template
}

func (r *Render) SeisComP() string {
	dataMap := map[string]string{
		"SampleRate":         fmt.Sprintf("%d", r.options.SampleRate),
		"NetworkCode":        r.options.NetworkCode,
		"StationCode":        r.options.StationCode,
		"LocationCode":       r.options.LocationCode,
		"StartTime":          r.options.StartTime.UTC().Format(time.RFC3339Nano),
		"Latitude":           fmt.Sprintf("%f", r.options.Latitude),
		"Longitude":          fmt.Sprintf("%f", r.options.Longitude),
		"Elevation":          fmt.Sprintf("%f", r.options.Elevation),
		"StationPlace":       r.options.StationPlace,
		"StationCountry":     r.options.StationCountry,
		"StationAffiliation": r.options.StationAffiliation,
		"StationDescription": r.options.StationDescription,
	}
	for idx := 1; idx <= r.channels; idx++ {
		if idx > len(r.options.ChannelCodes) {
			dataMap[fmt.Sprintf("ChannelCode%d", idx)] = fmt.Sprintf("CH%d", idx)
		} else {
			dataMap[fmt.Sprintf("ChannelCode%d", idx)] = r.options.ChannelCodes[idx-1]
		}
	}

	var stringBuf strings.Builder
	if err := r.templateSeisComP.Execute(&stringBuf, dataMap); err != nil {
		return err.Error()
	}

	return stringBuf.String()
}

func (r *Render) StationXML() string {
	dataMap := map[string]string{
		"SampleRate":         fmt.Sprintf("%d", r.options.SampleRate),
		"NetworkCode":        r.options.NetworkCode,
		"StationCode":        r.options.StationCode,
		"LocationCode":       r.options.LocationCode,
		"StartTime":          r.options.StartTime.UTC().Format(time.RFC3339Nano),
		"Latitude":           fmt.Sprintf("%f", r.options.Latitude),
		"Longitude":          fmt.Sprintf("%f", r.options.Longitude),
		"Elevation":          fmt.Sprintf("%f", r.options.Elevation),
		"StationPlace":       r.options.StationPlace,
		"StationCountry":     r.options.StationCountry,
		"StationAffiliation": r.options.StationAffiliation,
		"StationDescription": r.options.StationDescription,
	}
	for idx := 1; idx <= r.channels; idx++ {
		if idx > len(r.options.ChannelCodes) {
			dataMap[fmt.Sprintf("ChannelCode%d", idx)] = fmt.Sprintf("CH%d", idx)
		} else {
			dataMap[fmt.Sprintf("ChannelCode%d", idx)] = r.options.ChannelCodes[idx-1]
		}
	}

	var stringBuf strings.Builder
	if err := r.templateStationXML.Execute(&stringBuf, dataMap); err != nil {
		return err.Error()
	}

	return stringBuf.String()
}
