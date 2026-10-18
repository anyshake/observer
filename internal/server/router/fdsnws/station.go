package fdsnws

import (
	"errors"
	"math"
	"time"

	"github.com/anyshake/observer/config"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/bclswl0827/obsgo/inventory/stationxml"
)

func (p *provider) Station(query fdsnwsgo.StationQuery) (*stationxml.Document, fdsnwsgo.ProviderErrorKind, error) {
	if query.UpdatedAfter != nil || query.MatchTimeSeries || query.IncludeAvailability {
		return nil, fdsnwsgo.ProviderErrorClient, errors.New("updatedafter, matchtimeseries and includeavailability are not supported")
	}
	network, station, location, err := p.stationCodes()
	if err != nil {
		return nil, fdsnwsgo.ProviderErrorServer, err
	}
	constraints := []config.IConstraint{
		&config.StationAffiliationConfigConstraintImpl{},
		&config.StationDescriptionConfigConstraintImpl{},
		&config.StationCountryConfigConstraintImpl{},
		&config.StationPlaceConfigConstraintImpl{},
	}
	values := make([]string, len(constraints))
	for i, constraint := range constraints {
		value, err := constraint.Get(p.handler)
		if err != nil {
			return nil, fdsnwsgo.ProviderErrorServer, err
		}
		values[i] = value.(string)
	}
	render, err := p.device.GetMetadata(values[0], values[1], values[2], values[3], network, station, location, false)
	if err != nil {
		return nil, fdsnwsgo.ProviderErrorServer, err
	}
	document, err := parseStationMetadata([]byte(render.StationXML()))
	if err != nil {
		return nil, fdsnwsgo.ProviderErrorServer, err
	}
	networks := make([]stationxml.Network, 0, len(document.Networks))
	for _, network := range document.Networks {
		stations := make([]stationxml.Station, 0, len(network.Stations))
		for _, station := range network.Stations {
			if !matchesGeographic(query.Geographic, station.Latitude, station.Longitude) {
				continue
			}
			channels := make([]stationxml.Channel, 0, len(station.Channels))
			for _, channel := range station.Channels {
				if !matchesEpoch(query, channel.StartDate, channel.EndDate) {
					continue
				}
				for _, selection := range query.Selections {
					if matches(selection.Network, network.Code) && matches(selection.Station, station.Code) &&
						matches(selection.Location, channel.LocationCode) && matches(selection.Channel, channel.Code) &&
						matchesTimeRange(selection.TimeRange, channel.StartDate, channel.EndDate) {
						channels = append(channels, channel)
						break
					}
				}
			}
			if len(channels) > 0 {
				station.Channels = channels
				stations = append(stations, station)
			}
		}
		if len(stations) > 0 {
			network.Stations = stations
			networks = append(networks, network)
		}
	}
	document.Networks = networks
	if len(networks) == 0 {
		return nil, fdsnwsgo.ProviderErrorNone, fdsnwsgo.ErrNoData
	}
	return document, fdsnwsgo.ProviderErrorNone, nil
}

func matchesTimeRange(query fdsnwsgo.TimeRange, startDate, endDate string) bool {
	start, _ := time.Parse(time.RFC3339Nano, startDate)
	end, _ := time.Parse(time.RFC3339Nano, endDate)
	return (query.Start == nil || end.IsZero() || !end.Before(*query.Start)) &&
		(query.End == nil || start.IsZero() || !start.After(*query.End))
}

func matchesEpoch(query fdsnwsgo.StationQuery, startDate, endDate string) bool {
	start, _ := time.Parse(time.RFC3339Nano, startDate)
	end, _ := time.Parse(time.RFC3339Nano, endDate)
	if query.StartBefore != nil && (start.IsZero() || !start.Before(*query.StartBefore)) {
		return false
	}
	if query.StartAfter != nil && (start.IsZero() || !start.After(*query.StartAfter)) {
		return false
	}
	if query.EndBefore != nil && (end.IsZero() || !end.Before(*query.EndBefore)) {
		return false
	}
	if query.EndAfter != nil && !end.IsZero() && !end.After(*query.EndAfter) {
		return false
	}
	return true
}

func matchesGeographic(query fdsnwsgo.GeographicConstraints, latitude, longitude float64) bool {
	if query.MinLatitude != nil && latitude < *query.MinLatitude ||
		query.MaxLatitude != nil && latitude > *query.MaxLatitude ||
		query.MinLongitude != nil && longitude < *query.MinLongitude ||
		query.MaxLongitude != nil && longitude > *query.MaxLongitude {
		return false
	}
	if query.Latitude != nil && query.Longitude != nil {
		toRadians := math.Pi / 180
		lat1, lat2 := latitude*toRadians, *query.Latitude*toRadians
		dlon := (longitude - *query.Longitude) * toRadians
		cosine := math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(dlon)
		radius := math.Acos(max(-1, min(1, cosine))) / toRadians
		if query.MinRadius != nil && radius < *query.MinRadius || query.MaxRadius != nil && radius > *query.MaxRadius {
			return false
		}
	}
	return true
}
