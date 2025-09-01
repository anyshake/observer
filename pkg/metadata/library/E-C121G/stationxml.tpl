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
        <Sensor resourceId="Sensor-E-C111G-VEL">
          <Model>LGT-4.5C</Model>
        </Sensor>
        <DataLogger resourceId="Datalogger-E-C111G"/>
        <Response>
          <InstrumentSensitivity>
            <Value>429496729.6</Value>
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
            <PolesZeros name="AS-E-C111G-VEL" resourceId="ResponsePAZ-E-C111G-VEL">
              <InputUnits>
                <Name>m/s</Name>
              </InputUnits>
              <OutputUnits>
                <Name>V</Name>
              </OutputUnits>
              <PzTransferFunctionType>LAPLACE (RADIANS/SECOND)</PzTransferFunctionType>
              <NormalizationFactor>1.0</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
              <Zero number="2">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Zero number="3">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Pole number="0">
                <Real>-2.2211060060879837</Real>
                <Imaginary>2.221776881419294</Imaginary>
              </Pole>
              <Pole number="1">
                <Real>-2.2211060060879837</Real>
                <Imaginary>-2.221776881419294</Imaginary>
              </Pole>
            </PolesZeros>
            <StageGain>
              <Value>429496729.6</Value>
              <Frequency>4.5</Frequency>
            </StageGain>
          </Stage>
          <Stage number="2">
            <Coefficients name="DL-E-C111G" resourceId="Datalogger-E-C111G">
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
        <Sensor resourceId="Sensor-E-C111G-VEL">
          <Model>LGT-4.5C</Model>
        </Sensor>
        <DataLogger resourceId="Datalogger-E-C111G"/>
        <Response>
          <InstrumentSensitivity>
            <Value>429496729.6</Value>
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
            <PolesZeros name="AS-E-C111G-VEL" resourceId="ResponsePAZ-E-C111G-VEL">
              <InputUnits>
                <Name>m/s</Name>
              </InputUnits>
              <OutputUnits>
                <Name>V</Name>
              </OutputUnits>
              <PzTransferFunctionType>LAPLACE (RADIANS/SECOND)</PzTransferFunctionType>
              <NormalizationFactor>1.0</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
              <Zero number="2">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Zero number="3">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Pole number="0">
                <Real>-2.2211060060879837</Real>
                <Imaginary>2.221776881419294</Imaginary>
              </Pole>
              <Pole number="1">
                <Real>-2.2211060060879837</Real>
                <Imaginary>-2.221776881419294</Imaginary>
              </Pole>
            </PolesZeros>
            <StageGain>
              <Value>429496729.6</Value>
              <Frequency>4.5</Frequency>
            </StageGain>
          </Stage>
          <Stage number="2">
            <Coefficients name="DL-E-C111G" resourceId="Datalogger-E-C111G">
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
        <Sensor resourceId="Sensor-E-C111G-VEL">
          <Model>LGT-4.5C</Model>
        </Sensor>
        <DataLogger resourceId="Datalogger-E-C111G"/>
        <Response>
          <InstrumentSensitivity>
            <Value>429496729.6</Value>
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
            <PolesZeros name="AS-E-C111G-VEL" resourceId="ResponsePAZ-E-C111G-VEL">
              <InputUnits>
                <Name>m/s</Name>
              </InputUnits>
              <OutputUnits>
                <Name>V</Name>
              </OutputUnits>
              <PzTransferFunctionType>LAPLACE (RADIANS/SECOND)</PzTransferFunctionType>
              <NormalizationFactor>1.0</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">4.5</NormalizationFrequency>
              <Zero number="2">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Zero number="3">
                <Real>0.0</Real>
                <Imaginary>0.0</Imaginary>
              </Zero>
              <Pole number="0">
                <Real>-2.2211060060879837</Real>
                <Imaginary>2.221776881419294</Imaginary>
              </Pole>
              <Pole number="1">
                <Real>-2.2211060060879837</Real>
                <Imaginary>-2.221776881419294</Imaginary>
              </Pole>
            </PolesZeros>
            <StageGain>
              <Value>429496729.6</Value>
              <Frequency>4.5</Frequency>
            </StageGain>
          </Stage>
          <Stage number="2">
            <Coefficients name="DL-E-C111G" resourceId="Datalogger-E-C111G">
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
      <Channel code="{{.ChannelCode4}}" startDate="{{.StartTime}}" restrictedStatus="open" locationCode="{{.LocationCode}}">
        <ns3:format xmlns:ns3="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10">None</ns3:format>
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
        <Sensor resourceId="Sensor-E-C111G-ACC">
          <Model>ICM-42688-P</Model>
        </Sensor>
        <DataLogger resourceId="Datalogger-E-C111G"/>
        <Response>
          <InstrumentSensitivity>
            <Value>32768.0</Value>
            <Frequency>None</Frequency>
            <InputUnits>
              <Name>m/s**2</Name>
              <Description>None</Description>
            </InputUnits>
            <OutputUnits>
              <Name>None</Name>
              <Description>None</Description>
            </OutputUnits>
          </InstrumentSensitivity>
          <Stage number="1">
            <PolesZeros name="AS-E-C111G-ACC" resourceId="ResponsePAZ-E-C111G-ACC">
              <InputUnits>
                <Name>m/s**2</Name>
              </InputUnits>
              <OutputUnits>
                <Name>V</Name>
              </OutputUnits>
              <PzTransferFunctionType>DIGITAL (Z-TRANSFORM)</PzTransferFunctionType>
              <NormalizationFactor>1.0</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">10.0</NormalizationFrequency>
            </PolesZeros>
            <StageGain>
              <Value>32768.0</Value>
              <Frequency>10.0</Frequency>
            </StageGain>
          </Stage>
          <Stage number="2">
            <Coefficients name="DL-E-C111G" resourceId="Datalogger-E-C111G">
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
      <Channel code="{{.ChannelCode5}}" startDate="{{.StartTime}}" restrictedStatus="open" locationCode="{{.LocationCode}}">
        <ns4:format xmlns:ns4="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10">None</ns4:format>
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
        <Sensor resourceId="Sensor-E-C111G-ACC">
          <Model>ICM-42688-P</Model>
        </Sensor>
        <DataLogger resourceId="Datalogger-E-C111G"/>
        <Response>
          <InstrumentSensitivity>
            <Value>32768.0</Value>
            <Frequency>None</Frequency>
            <InputUnits>
              <Name>m/s**2</Name>
              <Description>None</Description>
            </InputUnits>
            <OutputUnits>
              <Name>None</Name>
              <Description>None</Description>
            </OutputUnits>
          </InstrumentSensitivity>
          <Stage number="1">
            <PolesZeros name="AS-E-C111G-ACC" resourceId="ResponsePAZ-E-C111G-ACC">
              <InputUnits>
                <Name>m/s**2</Name>
              </InputUnits>
              <OutputUnits>
                <Name>V</Name>
              </OutputUnits>
              <PzTransferFunctionType>DIGITAL (Z-TRANSFORM)</PzTransferFunctionType>
              <NormalizationFactor>1.0</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">10.0</NormalizationFrequency>
            </PolesZeros>
            <StageGain>
              <Value>32768.0</Value>
              <Frequency>10.0</Frequency>
            </StageGain>
          </Stage>
          <Stage number="2">
            <Coefficients name="DL-E-C111G" resourceId="Datalogger-E-C111G">
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
      <Channel code="{{.ChannelCode6}}" startDate="{{.StartTime}}" restrictedStatus="open" locationCode="{{.LocationCode}}">
        <ns5:format xmlns:ns5="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10">None</ns5:format>
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
        <Sensor resourceId="Sensor-E-C111G-ACC">
          <Model>ICM-42688-P</Model>
        </Sensor>
        <DataLogger resourceId="Datalogger-E-C111G"/>
        <Response>
          <InstrumentSensitivity>
            <Value>32768.0</Value>
            <Frequency>None</Frequency>
            <InputUnits>
              <Name>m/s**2</Name>
              <Description>None</Description>
            </InputUnits>
            <OutputUnits>
              <Name>None</Name>
              <Description>None</Description>
            </OutputUnits>
          </InstrumentSensitivity>
          <Stage number="1">
            <PolesZeros name="AS-E-C111G-ACC" resourceId="ResponsePAZ-E-C111G-ACC">
              <InputUnits>
                <Name>m/s**2</Name>
              </InputUnits>
              <OutputUnits>
                <Name>V</Name>
              </OutputUnits>
              <PzTransferFunctionType>DIGITAL (Z-TRANSFORM)</PzTransferFunctionType>
              <NormalizationFactor>1.0</NormalizationFactor>
              <NormalizationFrequency unit="HERTZ">10.0</NormalizationFrequency>
            </PolesZeros>
            <StageGain>
              <Value>32768.0</Value>
              <Frequency>10.0</Frequency>
            </StageGain>
          </Stage>
          <Stage number="2">
            <Coefficients name="DL-E-C111G" resourceId="Datalogger-E-C111G">
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
