package ed001

import (
	"text/template"
)

func getSeisCompTemplate() (*template.Template, error) {
	return template.New("E-D001").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<seiscomp xmlns="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10" version="0.10">
    <Inventory>
        <sensor publicID="Sensor-E-D001-VEL" name="S-E-D001-VEL"
            response="ResponsePAZ-E-D001-VEL">
            <model>LGT-4.5C</model>
            <unit>m/s</unit>
            <remark>{"unit":"Velocity in Meters Per Second"}</remark>
        </sensor>
        <datalogger publicID="Datalogger-E-D001" name="DL-E-D001">
            <recorderModel>E-D001</recorderModel>
            <recorderManufacturer>SensePlex Limited</recorderManufacturer>
            <gain>1</gain>
            <maxClockDrift>0</maxClockDrift>
            <decimation sampleRateNumerator="{{.SampleRate}}" sampleRateDenominator="1" />
        </datalogger>
        <responsePAZ publicID="ResponsePAZ-E-D001-VEL" name="AS-E-D001-VEL">
            <type>A</type>
            <gain>85899345920</gain>
            <gainFrequency>4.5</gainFrequency>
            <normalizationFactor>1</normalizationFactor>
            <normalizationFrequency>4.5</normalizationFrequency>
            <numberOfZeros>3</numberOfZeros>
            <numberOfPoles>4</numberOfPoles>
            <zeros>(0,0) (0,0) (28.27,0)</zeros>
            <poles>(-3.1416,0) (-19.991,19.999) (-19.991,-19.999)</poles>
        </responsePAZ>
        <network publicID="{{.NetworkCode}}.Network" code="{{.NetworkCode}}">
            <start>{{.StartTime}}</start>
            <description>Realtime seismic network of AnyShake Project.</description>
            <institutions>AnyShake Project</institutions>
            <type>SP</type>
            <station publicID="{{.StationCode}}.Station" code="{{.StationCode}}">
                <start>{{.StartTime}}</start>
                <description>{{.StationDescription}}</description>
                <latitude>{{.Latitude}}</latitude>
                <longitude>{{.Longitude}}</longitude>
                <elevation>{{.Elevation}}</elevation>
                <place>{{.StationPlace}}</place>
                <country>{{.StationCountry}}</country>
                <affiliation>{{.StationAffiliation}}</affiliation>
                <type>SP</type>
                <sensorLocation publicID="{{.StationCode}}.{{.LocationCode}}.Location"
                    code="{{.LocationCode}}">
                    <start>{{.StartTime}}</start>
                    <latitude>{{.Latitude}}</latitude>
                    <longitude>{{.Longitude}}</longitude>
                    <elevation>{{.Elevation}}</elevation>
                    <stream publicID="Stream/E-D001-{{.ChannelCode0}}" code="{{.ChannelCode0}}"
                        datalogger="Datalogger-E-D001"
                        sensor="Sensor-E-D001-VEL">
                        <start>{{.StartTime}}</start>
                        <dataloggerChannel>0</dataloggerChannel>
                        <sensorChannel>0</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>-90</dip>
                        <gain>85899345920</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-D001-{{.ChannelCode1}}" code="{{.ChannelCode1}}"
                        datalogger="Datalogger-E-D001"
                        sensor="Sensor-E-D001-VEL">
                        <start>{{.StartTime}}</start>
                        <dataloggerChannel>1</dataloggerChannel>
                        <sensorChannel>1</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>90</azimuth>
                        <dip>0</dip>
                        <gain>85899345920</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-D001-{{.ChannelCode2}}" code="{{.ChannelCode2}}"
                        datalogger="Datalogger-E-D001"
                        sensor="Sensor-E-D001-VEL">
                        <start>{{.StartTime}}</start>
                        <dataloggerChannel>2</dataloggerChannel>
                        <sensorChannel>2</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>0</dip>
                        <gain>85899345920</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                </sensorLocation>
            </station>
        </network>
    </Inventory>
</seiscomp>`)
}

func getStationXMLTemplate() (*template.Template, error) {
	return template.New("E-D001").Parse(`<?xml version='1.0' encoding='UTF-8'?>
<FDSNStationXML xmlns="http://www.fdsn.org/xml/station/1" schemaVersion="1.2">
    <Source>scxml import</Source>
    <Sender>AnyShake Project</Sender>
    <Module />
    <ModuleURI />
    <Created>{{.StartTime}}</Created>
    <Network code="AS" startDate="{{.StartTime}}" restrictedStatus="closed">
        <Description>Realtime seismic network of AnyShake Project.</Description>
        <Station code="{{.StationCode}}" startDate="{{.StartTime}}" restrictedStatus="closed">
            <Latitude unit="DEGREES">{{.Latitude}}</Latitude>
            <Longitude unit="DEGREES">{{.Longitude}}</Longitude>
            <Elevation>{{.Elevation}}</Elevation>
            <Site>
                <Name>{{.StationDescription}}</Name>
                <Town>{{.StationPlace}}</Town>
                <Country>{{.StationCountry}}</Country>
            </Site>
            <CreationDate>{{.StartTime}}</CreationDate>
            <Channel code="{{.ChannelCode0}}" startDate="{{.StartTime}}" restrictedStatus="closed"
                locationCode="{{.LocationCode}}">
                <ns0:format xmlns:ns0="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10">None</ns0:format>
                <Latitude unit="DEGREES">{{.Latitude}}</Latitude>
                <Longitude unit="DEGREES">{{.Longitude}}</Longitude>
                <Elevation>{{.Elevation}}</Elevation>
                <Depth>0.0</Depth>
                <Azimuth unit="DEGREES">0.0</Azimuth>
                <Dip unit="DEGREES">-90.0</Dip>
                <SampleRate unit="SAMPLES/S">{{.SampleRate}}</SampleRate>
                <SampleRateRatio>
                    <NumberSamples>{{.SampleRate}}</NumberSamples>
                    <NumberSeconds>1</NumberSeconds>
                </SampleRateRatio>
                <ClockDrift unit="SECONDS/SAMPLE">0.0</ClockDrift>
                <Sensor resourceId="Sensor-E-D001-VEL">
                    <Model>LGT-4.5C</Model>
                </Sensor>
                <DataLogger resourceId="Datalogger-E-D001" />
                <Response>
                    <InstrumentSensitivity>
                        <Value>85899345920.0</Value>
                        <Frequency>4.5</Frequency>
                        <InputUnits>
                            <Name>m/s</Name>
                            <Description>None</Description>
                        </InputUnits>
                        <OutputUnits>
                            <Name>None</Name>
                            <Description>None</Description>
                        </OutputUnits>
                    </InstrumentSensitivity>
                    <Stage number="1">
                        <PolesZeros name="AS-E-D001-VEL" resourceId="ResponsePAZ-E-D001-VEL">
                            <InputUnits>
                                <Name>m/s</Name>
                            </InputUnits>
                            <OutputUnits>
                                <Name>V</Name>
                            </OutputUnits>
                            <PzTransferFunctionType>LAPLACE (RADIANS/SECOND)</PzTransferFunctionType>
                            <NormalizationFactor>1.0</NormalizationFactor>
                            <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
                            <Zero number="3">
                                <Real>0.0</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Zero number="4">
                                <Real>0.0</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Zero number="5">
                                <Real>28.27</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Pole number="0">
                                <Real>-3.1416</Real>
                                <Imaginary>0.0</Imaginary>
                            </Pole>
                            <Pole number="1">
                                <Real>-19.991</Real>
                                <Imaginary>19.999</Imaginary>
                            </Pole>
                            <Pole number="2">
                                <Real>-19.991</Real>
                                <Imaginary>-19.999</Imaginary>
                            </Pole>
                        </PolesZeros>
                        <StageGain>
                            <Value>85899345920.0</Value>
                            <Frequency>4.5</Frequency>
                        </StageGain>
                    </Stage>
                    <Stage number="2">
                        <Coefficients name="DL-E-D001" resourceId="Datalogger-E-D001">
                            <InputUnits>
                                <Name>V</Name>
                            </InputUnits>
                            <OutputUnits>
                                <Name>COUNTS</Name>
                            </OutputUnits>
                            <CfTransferFunctionType>DIGITAL</CfTransferFunctionType>
                        </Coefficients>
                        <Decimation>
                            <InputSampleRate unit="HERTZ">{{.SampleRate}}</InputSampleRate>
                            <Factor>1</Factor>
                            <Offset>0</Offset>
                            <Delay>0.0</Delay>
                            <Correction>0.0</Correction>
                        </Decimation>
                        <StageGain>
                            <Value>1.0</Value>
                            <Frequency>0.0</Frequency>
                        </StageGain>
                    </Stage>
                </Response>
            </Channel>
            <Channel code="{{.ChannelCode1}}" startDate="{{.StartTime}}" restrictedStatus="closed"
                locationCode="{{.LocationCode}}">
                <ns1:format xmlns:ns1="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10">None</ns1:format>
                <Latitude unit="DEGREES">{{.Latitude}}</Latitude>
                <Longitude unit="DEGREES">{{.Longitude}}</Longitude>
                <Elevation>{{.Elevation}}</Elevation>
                <Depth>0.0</Depth>
                <Azimuth unit="DEGREES">90.0</Azimuth>
                <Dip unit="DEGREES">0.0</Dip>
                <SampleRate unit="SAMPLES/S">{{.SampleRate}}</SampleRate>
                <SampleRateRatio>
                    <NumberSamples>{{.SampleRate}}</NumberSamples>
                    <NumberSeconds>1</NumberSeconds>
                </SampleRateRatio>
                <ClockDrift unit="SECONDS/SAMPLE">0.0</ClockDrift>
                <Sensor resourceId="Sensor-E-D001-VEL">
                    <Model>LGT-4.5C</Model>
                </Sensor>
                <DataLogger resourceId="Datalogger-E-D001" />
                <Response>
                    <InstrumentSensitivity>
                        <Value>85899345920.0</Value>
                        <Frequency>4.5</Frequency>
                        <InputUnits>
                            <Name>m/s</Name>
                            <Description>None</Description>
                        </InputUnits>
                        <OutputUnits>
                            <Name>None</Name>
                            <Description>None</Description>
                        </OutputUnits>
                    </InstrumentSensitivity>
                    <Stage number="1">
                        <PolesZeros name="AS-E-D001-VEL" resourceId="ResponsePAZ-E-D001-VEL">
                            <InputUnits>
                                <Name>m/s</Name>
                            </InputUnits>
                            <OutputUnits>
                                <Name>V</Name>
                            </OutputUnits>
                            <PzTransferFunctionType>LAPLACE (RADIANS/SECOND)</PzTransferFunctionType>
                            <NormalizationFactor>1.0</NormalizationFactor>
                            <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
                            <Zero number="3">
                                <Real>0.0</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Zero number="4">
                                <Real>0.0</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Zero number="5">
                                <Real>28.27</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Pole number="0">
                                <Real>-3.1416</Real>
                                <Imaginary>0.0</Imaginary>
                            </Pole>
                            <Pole number="1">
                                <Real>-19.991</Real>
                                <Imaginary>19.999</Imaginary>
                            </Pole>
                            <Pole number="2">
                                <Real>-19.991</Real>
                                <Imaginary>-19.999</Imaginary>
                            </Pole>
                        </PolesZeros>
                        <StageGain>
                            <Value>85899345920.0</Value>
                            <Frequency>4.5</Frequency>
                        </StageGain>
                    </Stage>
                    <Stage number="2">
                        <Coefficients name="DL-E-D001" resourceId="Datalogger-E-D001">
                            <InputUnits>
                                <Name>V</Name>
                            </InputUnits>
                            <OutputUnits>
                                <Name>COUNTS</Name>
                            </OutputUnits>
                            <CfTransferFunctionType>DIGITAL</CfTransferFunctionType>
                        </Coefficients>
                        <Decimation>
                            <InputSampleRate unit="HERTZ">{{.SampleRate}}</InputSampleRate>
                            <Factor>1</Factor>
                            <Offset>0</Offset>
                            <Delay>0.0</Delay>
                            <Correction>0.0</Correction>
                        </Decimation>
                        <StageGain>
                            <Value>1.0</Value>
                            <Frequency>0.0</Frequency>
                        </StageGain>
                    </Stage>
                </Response>
            </Channel>
            <Channel code="{{.ChannelCode2}}" startDate="{{.StartTime}}" restrictedStatus="closed"
                locationCode="{{.LocationCode}}">
                <ns2:format xmlns:ns2="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10">None</ns2:format>
                <Latitude unit="DEGREES">{{.Latitude}}</Latitude>
                <Longitude unit="DEGREES">{{.Longitude}}</Longitude>
                <Elevation>{{.Elevation}}</Elevation>
                <Depth>0.0</Depth>
                <Azimuth unit="DEGREES">0.0</Azimuth>
                <Dip unit="DEGREES">0.0</Dip>
                <SampleRate unit="SAMPLES/S">{{.SampleRate}}</SampleRate>
                <SampleRateRatio>
                    <NumberSamples>{{.SampleRate}}</NumberSamples>
                    <NumberSeconds>1</NumberSeconds>
                </SampleRateRatio>
                <ClockDrift unit="SECONDS/SAMPLE">0.0</ClockDrift>
                <Sensor resourceId="Sensor-E-D001-VEL">
                    <Model>LGT-4.5C</Model>
                </Sensor>
                <DataLogger resourceId="Datalogger-E-D001" />
                <Response>
                    <InstrumentSensitivity>
                        <Value>85899345920.0</Value>
                        <Frequency>4.5</Frequency>
                        <InputUnits>
                            <Name>m/s</Name>
                            <Description>None</Description>
                        </InputUnits>
                        <OutputUnits>
                            <Name>None</Name>
                            <Description>None</Description>
                        </OutputUnits>
                    </InstrumentSensitivity>
                    <Stage number="1">
                        <PolesZeros name="AS-E-D001-VEL" resourceId="ResponsePAZ-E-D001-VEL">
                            <InputUnits>
                                <Name>m/s</Name>
                            </InputUnits>
                            <OutputUnits>
                                <Name>V</Name>
                            </OutputUnits>
                            <PzTransferFunctionType>LAPLACE (RADIANS/SECOND)</PzTransferFunctionType>
                            <NormalizationFactor>1.0</NormalizationFactor>
                            <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
                            <Zero number="3">
                                <Real>0.0</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Zero number="4">
                                <Real>0.0</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Zero number="5">
                                <Real>28.27</Real>
                                <Imaginary>0.0</Imaginary>
                            </Zero>
                            <Pole number="0">
                                <Real>-3.1416</Real>
                                <Imaginary>0.0</Imaginary>
                            </Pole>
                            <Pole number="1">
                                <Real>-19.991</Real>
                                <Imaginary>19.999</Imaginary>
                            </Pole>
                            <Pole number="2">
                                <Real>-19.991</Real>
                                <Imaginary>-19.999</Imaginary>
                            </Pole>
                        </PolesZeros>
                        <StageGain>
                            <Value>85899345920.0</Value>
                            <Frequency>4.5</Frequency>
                        </StageGain>
                    </Stage>
                    <Stage number="2">
                        <Coefficients name="DL-E-D001" resourceId="Datalogger-E-D001">
                            <InputUnits>
                                <Name>V</Name>
                            </InputUnits>
                            <OutputUnits>
                                <Name>COUNTS</Name>
                            </OutputUnits>
                            <CfTransferFunctionType>DIGITAL</CfTransferFunctionType>
                        </Coefficients>
                        <Decimation>
                            <InputSampleRate unit="HERTZ">{{.SampleRate}}</InputSampleRate>
                            <Factor>1</Factor>
                            <Offset>0</Offset>
                            <Delay>0.0</Delay>
                            <Correction>0.0</Correction>
                        </Decimation>
                        <StageGain>
                            <Value>1.0</Value>
                            <Frequency>0.0</Frequency>
                        </StageGain>
                    </Stage>
                </Response>
            </Channel>
        </Station>
    </Network>
</FDSNStationXML>`)
}
