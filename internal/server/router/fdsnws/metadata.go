package fdsnws

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/bclswl0827/obsgo/inventory/stationxml"
)

func parseStationMetadata(data []byte) (*stationxml.Document, error) {
	// Remove incomplete optional nodes before Parse turns "None" frequencies into zero.
	var filtered bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(&filtered)
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if start, ok := token.(xml.StartElement); ok && (start.Name.Local == "InstrumentSensitivity" || start.Name.Local == "StageGain") {
			var fields struct {
				Value     string
				Frequency string
				XML       string `xml:",innerxml"`
			}
			if err := decoder.DecodeElement(&fields, &start); err != nil {
				return nil, err
			}
			if !hasMetadataValue(fields.Value) || !hasMetadataValue(fields.Frequency) {
				continue
			}
			if err := encoder.EncodeElement(struct {
				XML string `xml:",innerxml"`
			}{fields.XML}, start); err != nil {
				return nil, err
			}
			continue
		}
		if err := encoder.EncodeToken(token); err != nil {
			return nil, err
		}
	}
	if err := encoder.Flush(); err != nil {
		return nil, err
	}
	document, err := stationxml.Parse(filtered.Bytes())
	if err != nil {
		return nil, err
	}
	for _, network := range document.Networks {
		for _, station := range network.Stations {
			for i := range station.Channels {
				channel := &station.Channels[i]
				response := channel.Response
				if response == nil {
					continue
				}
				stages := response.Stages[:0]
				for _, stage := range response.Stages {
					if stage.PolesZeros != nil || stage.Coefficients != nil || stage.FIR != nil || stage.Decimation != nil || stage.Gain != nil {
						stages = append(stages, stage)
					}
				}
				response.Stages = stages
				if response.Sensitivity == nil {
					if len(stages) == 0 {
						channel.Response = nil
					}
					continue
				}
				// Overall units can also be supplied by the first and last filter stages.
				var input, output stationxml.Units
				firstFilter := true
				for _, stage := range response.Stages {
					var stageInput, stageOutput stationxml.Units
					switch {
					case stage.PolesZeros != nil:
						stageInput, stageOutput = stage.PolesZeros.Input, stage.PolesZeros.Output
					case stage.Coefficients != nil:
						stageInput, stageOutput = stage.Coefficients.Input, stage.Coefficients.Output
					case stage.FIR != nil:
						stageInput, stageOutput = stage.FIR.Input, stage.FIR.Output
					default:
						continue
					}
					if firstFilter {
						input = stageInput
						firstFilter = false
					}
					output = stageOutput
				}
				response.Sensitivity.Input = metadataUnits(response.Sensitivity.Input, input)
				response.Sensitivity.Output = metadataUnits(response.Sensitivity.Output, output)
				if !hasMetadataValue(response.Sensitivity.Input.Name) || !hasMetadataValue(response.Sensitivity.Output.Name) {
					response.Sensitivity = nil
					if len(stages) == 0 {
						channel.Response = nil
					}
				}
			}
		}
	}
	return document, nil
}

func hasMetadataValue(value string) bool {
	value = strings.TrimSpace(value)
	return value != "" && !strings.EqualFold(value, "None")
}

func metadataUnits(value, fallback stationxml.Units) stationxml.Units {
	if !hasMetadataValue(value.Name) {
		value = fallback
	}
	if !hasMetadataValue(value.Description) {
		value.Description = ""
	}
	return value
}
