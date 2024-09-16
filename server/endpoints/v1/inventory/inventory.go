package inventory

import (
	"fmt"
	"math"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/drivers/explorer"
)

func (i *Inventory) handleInventory(config *config.Config, explorerDeps *explorer.ExplorerDependency) string {
	const xmlTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<seiscomp xmlns="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10" version="0.10">
	<Inventory>
		<sensor publicID="Sensor@AnyShake" name="Geophone" response="ResponsePAZ@AnyShake">
			<type>SPS</type>
			<unit>M/S</unit>
			<lowFrequency>0.1</lowFrequency>
			<highFrequency>%d</highFrequency>
		</sensor>
		<datalogger publicID="Datalogger@AnyShake" name="ADC_WITH_FULLSCALE">
			<description>Analog to digital converter</description>
			<gain>%f</gain>
			<decimation sampleRateNumerator="%d" sampleRateDenominator="1" />
		</datalogger>
		<responsePAZ publicID="ResponsePAZ@AnyShake" name="Geophone">
			<type>A</type>
			<gain>%f</gain>
			<gainFrequency>%f</gainFrequency>
			<normalizationFactor>1</normalizationFactor>
			<normalizationFrequency>%f</normalizationFrequency>
		</responsePAZ>
		<network publicID="Network@AnyShake" code="%s">
			<start>%s</start>
			<description>AnyShake seismic network</description>
			<institutions>AnyShake Project</institutions>
			<region>%s</region>
			<shared>true</shared>
			<type>BB</type>
			<station publicID="Station@AnyShake" code="%s">
				<start>%s</start>
				<description>%s</description>
				<latitude>%f</latitude>
				<longitude>%f</longitude>
				<elevation>%f</elevation>
				<place>%s</place>
				<country>%s</country>
				<affiliation>%s</affiliation>
				<type>BB</type>
				<shared>true</shared>
				<sensorLocation publicID="SensorLocation@AnyShake" code="%s">
					<start>%s</start>
					<latitude>%f</latitude>
					<longitude>%f</longitude>
					<elevation>%f</elevation>
					<stream publicID="Stream#EHN" code="EHN" datalogger="Datalogger@AnyShake"
						sensor="Sensor@AnyShake">
						<start>%s</start>
						<sampleRateNumerator>%d</sampleRateNumerator>
						<sampleRateDenominator>1</sampleRateDenominator>
						<depth>0</depth>
						<azimuth>0</azimuth>
						<dip>0</dip>
						<gain>%f</gain>
						<gainFrequency>1</gainFrequency>
						<gainUnit>M/S</gainUnit>
						<format>Steim2</format>
						<flags>GC</flags>
						<shared>true</shared>
					</stream>
					<stream publicID="Stream#EHE" code="EHE" datalogger="Datalogger@AnyShake"
						sensor="Sensor@AnyShake">
						<start>%s</start>
						<sampleRateNumerator>%d</sampleRateNumerator>
						<sampleRateDenominator>1</sampleRateDenominator>
						<depth>0</depth>
						<azimuth>90</azimuth>
						<dip>0</dip>
						<gain>%f</gain>
						<gainFrequency>1</gainFrequency>
						<gainUnit>M/S</gainUnit>
						<format>Steim2</format>
						<flags>GC</flags>
						<shared>true</shared>
					</stream>
					<stream publicID="Stream#EHZ" code="EHZ" datalogger="Datalogger@AnyShake"
						sensor="Sensor@AnyShake">
						<start>%s</start>
						<sampleRateNumerator>%d</sampleRateNumerator>
						<sampleRateDenominator>1</sampleRateDenominator>
						<depth>0</depth>
						<azimuth>0</azimuth>
						<dip>90</dip>
						<gain>%f</gain>
						<gainFrequency>1</gainFrequency>
						<gainUnit>M/S</gainUnit>
						<format>Steim2</format>
						<flags>GC</flags>
						<shared>true</shared>
					</stream>
				</sensorLocation>
			</station>
		</network>
	</Inventory>
</seiscomp>
`

	var (
		currentSampleRate = explorerDeps.Health.GetSampleRate()
		startTime         = explorerDeps.Health.GetStartTime()
	)

	sensorHighFrequency := currentSampleRate / 2
	dataloggerGain := math.Pow(2, float64(config.Sensor.Resolution-1)) / config.Sensor.FullScale
	dataloggerSampleRateNumerator := currentSampleRate
	responsePAZGain := config.Sensor.Sensitivity
	responsePAZGainFrequency := config.Sensor.Frequency
	responsePAZGainNormalizationFrequency := config.Sensor.Frequency
	networkCode := config.Stream.Network
	networkStart := startTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	networkRegion := config.Station.Region
	stationCode := config.Stream.Station
	stationStart := startTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	stationDescription := fmt.Sprintf("AnyShake Station in %s", config.Station.City)
	stationLatitude := config.Location.Latitude
	stationLongitude := config.Location.Longitude
	stationElevation := config.Location.Elevation
	stationCity := config.Station.City
	stationCountry := config.Station.Country
	stationAffiliation := config.Station.Owner
	sensorLocationCode := config.Stream.Location
	sensorLocationStart := startTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	sensorLocationLatitude := config.Location.Latitude
	sensorLocationLongitude := config.Location.Longitude
	sensorLocationElevation := config.Location.Elevation

	// Stream settings
	streamStart := startTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	streamSampleRateNumerator := currentSampleRate
	streamGain := dataloggerGain * config.Sensor.Sensitivity

	return fmt.Sprintf(
		xmlTemplate,
		sensorHighFrequency,
		dataloggerGain,
		dataloggerSampleRateNumerator,
		responsePAZGain,
		responsePAZGainFrequency,
		responsePAZGainNormalizationFrequency,
		networkCode,
		networkStart,
		networkRegion,
		stationCode,
		stationStart,
		stationDescription,
		stationLatitude,
		stationLongitude,
		stationElevation,
		stationCity,
		stationCountry,
		stationAffiliation,
		sensorLocationCode,
		sensorLocationStart,
		sensorLocationLatitude,
		sensorLocationLongitude,
		sensorLocationElevation,
		// Stream for EHZ
		streamStart,
		streamSampleRateNumerator,
		streamGain,
		// Stream for EHE
		streamStart,
		streamSampleRateNumerator,
		streamGain,
		// Stream for EHN
		streamStart,
		streamSampleRateNumerator,
		streamGain,
	)
}
