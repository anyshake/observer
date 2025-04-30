package ed001

import (
	"fmt"
	"strings"
	"time"
)

func (m *ED001_MetadataImpl) StationXML() string {
	if len(m.ChannelCodes) != 3 {
		return "only 3 channels are supported on E-D001"
	}

	tpl, err := getStationXMLTemplate()
	if err != nil {
		return err.Error()
	}

	dataMap := map[string]string{
		"SampleRate":         fmt.Sprintf("%d", m.SampleRate),
		"NetworkCode":        m.NetworkCode,
		"StationCode":        m.StationCode,
		"LocationCode":       m.LocationCode,
		"StartTime":          m.StartTime.UTC().Format(time.RFC3339Nano),
		"Latitude":           fmt.Sprintf("%f", m.Latitude),
		"Longitude":          fmt.Sprintf("%f", m.Longitude),
		"Elevation":          fmt.Sprintf("%f", m.Elevation),
		"StationPlace":       m.StationPlace,
		"StationCountry":     m.StationCountry,
		"StationAffiliation": m.StationAffiliation,
		"StationDescription": m.StationDescription,
	}
	for idx, channelCode := range m.ChannelCodes {
		dataMap[fmt.Sprintf("ChannelCode%d", idx)] = channelCode
	}

	var stringBuf strings.Builder
	err = tpl.Execute(&stringBuf, dataMap)
	if err != nil {
		return err.Error()
	}

	return stringBuf.String()
}
