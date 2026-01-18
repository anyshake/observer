<?xml version="1.0" encoding="UTF-8"?>
<seiscomp xmlns="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10" version="0.10">
    <Inventory>
        <sensor publicID="Sensor-E-C111G-VEL" name="S-E-C111G-VEL"
            response="ResponsePAZ-E-C111G-VEL">
            <model>{{.VelocitySensorModel}}</model>
            <unit>m/s</unit>
            <remark>{"unit":"Velocity in Meters Per Second"}</remark>
        </sensor>
        <sensor publicID="Sensor-E-C111G-ACC" name="S-E-C111G-ACC"
            response="ResponsePAZ-E-C111G-ACC">
            <model>{{.AcceleroSensorModel}}</model>
            <unit>m/s**2</unit>
            <remark>{"unit":"Acceleration in Meters Per Second Squared"}</remark>
        </sensor>
        <datalogger publicID="Datalogger-E-C111G-VEL" name="DL-E-C111G-VEL">
            <recorderModel>E-C111G</recorderModel>
            <recorderManufacturer>SensePlex Limited</recorderManufacturer>
            <gain>{{.VelocityDataLoggerGain}}</gain>
            <maxClockDrift>0</maxClockDrift>
            <decimation sampleRateNumerator="{{.SampleRate}}" sampleRateDenominator="1" />
        </datalogger>
        <datalogger publicID="Datalogger-E-C111G-ACC" name="DL-E-C111G-ACC">
            <recorderModel>E-C111G</recorderModel>
            <recorderManufacturer>SensePlex Limited</recorderManufacturer>
            <gain>{{.AcceleroDataLoggerGain}}</gain>
            <maxClockDrift>0</maxClockDrift>
            <decimation sampleRateNumerator="{{.SampleRate}}" sampleRateDenominator="1" />
        </datalogger>
        <responsePAZ publicID="ResponsePAZ-E-C111G-VEL" name="AS-E-C111G-VEL">
            <type>A</type>
            <gain>{{.VelocitySensorGain}}</gain>
            <gainFrequency>4.5</gainFrequency>
            <normalizationFactor>171.99852139050935</normalizationFactor>
            <normalizationFrequency>4.5</normalizationFrequency>
            <numberOfZeros>2</numberOfZeros>
            <numberOfPoles>3</numberOfPoles>
            <zeros>(0,0) (0,0)</zeros>
            <poles>(-2.2211060060879837,2.221776881419294) (-2.2211060060879837,-2.221776881419294) (-169.64600329384882,0)</poles>
        </responsePAZ>
        <responsePAZ publicID="ResponsePAZ-E-C111G-ACC" name="AS-E-C111G-ACC">
            <type>D</type>
            <gain>{{.AcceleroSensorGain}}</gain>
            <gainFrequency>1</gainFrequency>
            <normalizationFactor>1</normalizationFactor>
            <normalizationFrequency>1</normalizationFrequency>
        </responsePAZ>
        <network publicID="{{.NetworkCode}}.Network" code="{{.NetworkCode}}">
            <!-- <start>{{.StartTime}}</start> -->
            <start>1970-01-01T00:00:00.0000Z</start>
            <description>Realtime seismic network of AnyShake Project.</description>
            <institutions>AnyShake Project</institutions>
            <type>SP</type>
            <station publicID="{{.StationCode}}.Station" code="{{.StationCode}}">
                <!-- <start>{{.StartTime}}</start> -->
                <start>1970-01-01T00:00:00.0000Z</start>
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
                    <!-- <start>{{.StartTime}}</start> -->
                    <start>1970-01-01T00:00:00.0000Z</start>
                    <latitude>{{.Latitude}}</latitude>
                    <longitude>{{.Longitude}}</longitude>
                    <elevation>{{.Elevation}}</elevation>
                    <stream publicID="Stream/E-C111G-CH1" code="{{.ChannelCode1}}"
                        datalogger="Datalogger-E-C111G-VEL"
                        sensor="Sensor-E-C111G-VEL">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>0</dataloggerChannel>
                        <sensorChannel>0</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>-90</dip>
                        <gain>{{.Channel1Gain}}</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C111G-CH2" code="{{.ChannelCode2}}"
                        datalogger="Datalogger-E-C111G-VEL"
                        sensor="Sensor-E-C111G-VEL">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>1</dataloggerChannel>
                        <sensorChannel>1</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>90</azimuth>
                        <dip>0</dip>
                        <gain>{{.Channel2Gain}}</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C111G-CH3" code="{{.ChannelCode3}}"
                        datalogger="Datalogger-E-C111G-VEL"
                        sensor="Sensor-E-C111G-VEL">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>2</dataloggerChannel>
                        <sensorChannel>2</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>0</dip>
                        <gain>{{.Channel3Gain}}</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C111G-CH4" code="{{.ChannelCode4}}"
                        datalogger="Datalogger-E-C111G-ACC"
                        sensor="Sensor-E-C111G-ACC">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>3</dataloggerChannel>
                        <sensorChannel>3</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>-90</dip>
                        <gain>{{.Channel4Gain}}</gain>
                        <gainUnit>m/s**2</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C111G-CH5" code="{{.ChannelCode5}}"
                        datalogger="Datalogger-E-C111G-ACC"
                        sensor="Sensor-E-C111G-ACC">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>4</dataloggerChannel>
                        <sensorChannel>4</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>90</azimuth>
                        <dip>0</dip>
                        <gain>{{.Channel5Gain}}</gain>
                        <gainUnit>m/s**2</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C111G-CH6" code="{{.ChannelCode6}}"
                        datalogger="Datalogger-E-C111G-ACC"
                        sensor="Sensor-E-C111G-ACC">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>5</dataloggerChannel>
                        <sensorChannel>5</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>0</dip>
                        <gain>{{.Channel6Gain}}</gain>
                        <gainUnit>m/s**2</gainUnit>
                    </stream>
                </sensorLocation>
            </station>
        </network>
    </Inventory>
</seiscomp>