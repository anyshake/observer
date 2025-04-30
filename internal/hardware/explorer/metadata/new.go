package metadata

import (
	"fmt"

	"github.com/anyshake/observer/internal/hardware/explorer/metadata/ec111g"
	"github.com/anyshake/observer/internal/hardware/explorer/metadata/ec121g"
)

func New(model string, options Options) (IMetadata, error) {
	switch model {
	case "E-C111G":
		return &ec111g.EC111G_MetadataImpl{
			StartTime:          options.StartTime,
			SampleRate:         options.SampleRate,
			Latitude:           options.Latitude,
			Longitude:          options.Longitude,
			Elevation:          options.Elevation,
			NetworkCode:        options.NetworkCode,
			StationCode:        options.StationCode,
			LocationCode:       options.LocationCode,
			ChannelCodes:       options.ChannelCodes,
			StationPlace:       options.StationPlace,
			StationCountry:     options.StationCountry,
			StationAffiliation: options.StationAffiliation,
			StationDescription: options.StationDescription,
		}, nil
	case "E-C121G":
		return &ec121g.EC121G_MetadataImpl{
			StartTime:          options.StartTime,
			SampleRate:         options.SampleRate,
			Latitude:           options.Latitude,
			Longitude:          options.Longitude,
			Elevation:          options.Elevation,
			NetworkCode:        options.NetworkCode,
			StationCode:        options.StationCode,
			LocationCode:       options.LocationCode,
			ChannelCodes:       options.ChannelCodes,
			StationPlace:       options.StationPlace,
			StationCountry:     options.StationCountry,
			StationAffiliation: options.StationAffiliation,
			StationDescription: options.StationDescription,
		}, nil
		// case "E-D001":
		// 	return &ed001.ED001_MetadataImpl{
		// 		StartTime:          options.StartTime,
		// 		SampleRate:         options.SampleRate,
		// 		Latitude:           options.Latitude,
		// 		Longitude:          options.Longitude,
		// 		Elevation:          options.Elevation,
		// 		NetworkCode:        options.NetworkCode,
		// 		StationCode:        options.StationCode,
		// 		LocationCode:       options.LocationCode,
		// 		ChannelCodes:       options.ChannelCodes,
		// 		StationPlace:       options.StationPlace,
		// 		StationCountry:     options.StationCountry,
		// 		StationAffiliation: options.StationAffiliation,
		// 		StationDescription: options.StationDescription,
		// 	}, nil
	}

	return nil, fmt.Errorf("unknown hardware model %s for generating metadata", model)
}
