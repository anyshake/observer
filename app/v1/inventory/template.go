package inventory

import (
	"fmt"
	"math"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/publisher"
)

func getInventoryString(config *config.Conf, status *publisher.Status) string {
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
			<description>AnyShake_Project_Seismic_Network</description>
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

	startTime := status.ReadyTime
	currentSampleRate := (len(status.Geophone.EHZ) + len(status.Geophone.EHE) + len(status.Geophone.EHN)) / 3
	if startTime.IsZero() || currentSampleRate == 0 {
		return ""
	}

	sensorHighFrequency := currentSampleRate / 2
	dataloggerGain := math.Pow(2, float64(config.ADC.Resolution-1)) / config.ADC.FullScale
	dataloggerSampleRateNumerator := currentSampleRate
	responsePAZGain := config.Geophone.Sensitivity
	responsePAZGainFrequency := config.Geophone.Frequency
	responsePAZGainNormalizationFrequency := config.Geophone.Frequency
	networkCode := config.Station.Network
	networkStart := status.ReadyTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	networkRegion := config.Station.Region
	stationCode := config.Station.Station
	stationStart := status.ReadyTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	stationDescription := fmt.Sprintf("AnyShake Station %s", config.Station.UUID)
	stationLatitude := config.Station.Latitude
	stationLongitude := config.Station.Longitude
	stationElevation := config.Station.Elevation
	stationCity := config.Station.City
	stationCountry := config.Station.Country
	stationAffiliation := config.Station.Owner
	sensorLocationCode := config.Station.Location
	sensorLocationStart := status.ReadyTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	sensorLocationLatitude := config.Station.Latitude
	sensorLocationLongitude := config.Station.Longitude
	sensorLocationElevation := config.Station.Elevation

	// Stream settings
	streamStart := status.ReadyTime.UTC().Format("2006-01-02T15:04:05.0000Z")
	streamSampleRateNumerator := currentSampleRate
	streamGain := dataloggerGain * config.Geophone.Sensitivity

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
