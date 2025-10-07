<?xml version='1.0' encoding='UTF-8'?>
<FDSNStationXML xmlns="http://www.fdsn.org/xml/station/1" schemaVersion="1.2">
  <Source>scxml import</Source>
  <Sender>AnyShake Project</Sender>
  <Module/>
  <ModuleURI/>
  <Created>{{.StartTime}}</Created>
  <Network code="{{.NetworkCode}}" startDate="{{.StartTime}}" restrictedStatus="open">
    <Description>Realtime seismic network of AnyShake Project.</Description>
    <Station code="{{.StationCode}}" startDate="{{.StartTime}}" restrictedStatus="open">
      <Latitude unit="DEGREES">{{.Latitude}}</Latitude>
      <Longitude unit="DEGREES">{{.Longitude}}</Longitude>
      <Elevation>{{.Elevation}}</Elevation>
      <Site>
        <Name>{{.StationDescription}}</Name>
        <Town>{{.StationPlace}}</Town>
        <Country>{{.StationCountry}}</Country>
      </Site>
      <CreationDate>{{.StartTime}}</CreationDate>
      <Channel code="{{.ChannelCode1}}" startDate="{{.StartTime}}" restrictedStatus="open" locationCode="{{.LocationCode}}">
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
        <DataLogger resourceId="Datalogger-E-D001"/>
        <Response>
          <InstrumentSensitivity>
            <Value>343597383520.0</Value>
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
              <NormalizationFactor>171.99852139050935</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
              <Zero number="3">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Zero number="4">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Pole number="0">
                <Real>-3.553769609740774</Real>
                <Imaginary>3.554843010270871</Imaginary>
              </Pole>
              <Pole number="1">
                <Real>-3.553769609740774</Real>
                <Imaginary>-3.554843010270871</Imaginary>
              </Pole>
              <Pole number="2">
                <Real>-169.64600329384882</Real>
                <Imaginary>0.0</Imaginary>
              </Pole>
            </PolesZeros>
            <StageGain>
              <Value>343597383520.0</Value>
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
      <Channel code="{{.ChannelCode2}}" startDate="{{.StartTime}}" restrictedStatus="open" locationCode="{{.LocationCode}}">
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
        <DataLogger resourceId="Datalogger-E-D001"/>
        <Response>
          <InstrumentSensitivity>
            <Value>343597383520.0</Value>
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
              <NormalizationFactor>171.99852139050935</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
              <Zero number="3">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Zero number="4">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Pole number="0">
                <Real>-3.553769609740774</Real>
                <Imaginary>3.554843010270871</Imaginary>
              </Pole>
              <Pole number="1">
                <Real>-3.553769609740774</Real>
                <Imaginary>-3.554843010270871</Imaginary>
              </Pole>
              <Pole number="2">
                <Real>-169.64600329384882</Real>
                <Imaginary>0.0</Imaginary>
              </Pole>
            </PolesZeros>
            <StageGain>
              <Value>343597383520.0</Value>
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
      <Channel code="{{.ChannelCode3}}" startDate="{{.StartTime}}" restrictedStatus="open" locationCode="{{.LocationCode}}">
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
        <DataLogger resourceId="Datalogger-E-D001"/>
        <Response>
          <InstrumentSensitivity>
            <Value>343597383520.0</Value>
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
              <NormalizationFactor>171.99852139050935</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
              <Zero number="3">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Zero number="4">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Pole number="0">
                <Real>-3.553769609740774</Real>
                <Imaginary>3.554843010270871</Imaginary>
              </Pole>
              <Pole number="1">
                <Real>-3.553769609740774</Real>
                <Imaginary>-3.554843010270871</Imaginary>
              </Pole>
              <Pole number="2">
                <Real>-169.64600329384882</Real>
                <Imaginary>0.0</Imaginary>
              </Pole>
            </PolesZeros>
            <StageGain>
              <Value>343597383520.0</Value>
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
</FDSNStationXML>
