package metadata

import (
	"fmt"
	"regexp"
	"text/template"
)

func New(model string, options Options) (*Render, error) {
	templateSeisComP, err := getMetadataFromLibrary(model, SEISCOMP_TEMPLATE)
	if err != nil {
		templateSeisComP, err = getMetadataFromLocalDisk(model, SEISCOMP_TEMPLATE)
		if err != nil {
			return nil, fmt.Errorf("failed to get SeisComP metadata template: %w", err)
		}
	}
	templateStationXML, err := getMetadataFromLibrary(model, STATIONXML_TEMPLATE)
	if err != nil {
		templateStationXML, err = getMetadataFromLocalDisk(model, STATIONXML_TEMPLATE)
		if err != nil {
			return nil, fmt.Errorf("failed to get StationXML metadata template: %w", err)
		}
	}

	channelCodeRegexp, err := regexp.Compile(`\{\{\.ChannelCode\d+\}\}`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile channel code regexp: %w", err)
	}
	channelsInTemplate := len(channelCodeRegexp.FindAllString(templateSeisComP, -1))

	templateSeisComPObj, err := template.New(fmt.Sprintf("%s-%s", model, SEISCOMP_TEMPLATE)).Parse(templateSeisComP)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SeisComP metadata template: %w", err)
	}
	templateStationXMLObj, err := template.New(fmt.Sprintf("%s-%s", model, STATIONXML_TEMPLATE)).Parse(templateStationXML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse StationXML metadata template: %w", err)
	}

	return &Render{
		options:            &options,
		channels:           channelsInTemplate,
		templateSeisComP:   templateSeisComPObj,
		templateStationXML: templateStationXMLObj,
	}, nil
}
