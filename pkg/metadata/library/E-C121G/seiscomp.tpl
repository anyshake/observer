<?xml version="1.0" encoding="UTF-8"?>
<seiscomp xmlns="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.10" version="0.10">
    <Inventory>
        <sensor publicID="Sensor-E-C121G-VEL" name="S-E-C121G-VEL"
            response="ResponsePAZ-E-C121G-VEL">
            <model>LGT-4.5C</model>
            <unit>m/s</unit>
            <remark>{"unit":"Velocity in Meters Per Second"}</remark>
        </sensor>
        <sensor publicID="Sensor-E-C121G-ACC" name="S-E-C121G-ACC"
            response="ResponsePAZ-E-C121G-ACC">
            <model>ICM-42688-P</model>
            <unit>m/s**2</unit>
            <remark>{"unit":"Acceleration in Meters Per Second Squared"}</remark>
        </sensor>
        <datalogger publicID="Datalogger-E-C121G-VEL" name="DL-E-C121G-VEL">
            <recorderModel>E-C121G</recorderModel>
            <recorderManufacturer>SensePlex Limited</recorderManufacturer>
            <!--
                Input amplifier gain: x1
                ADC PGA gain: x1
                ADC full-scale: 2.5 V
                ==> ((2 ** 31) - 1) * 2 / 2.5
            -->
            <gain>1717986917.6</gain>
            <maxClockDrift>0</maxClockDrift>
            <decimation sampleRateNumerator="{{.SampleRate}}" sampleRateDenominator="1" />
        </datalogger>
        <datalogger publicID="Datalogger-E-C121G-ACC" name="DL-E-C121G-ACC">
            <recorderModel>E-C121G</recorderModel>
            <recorderManufacturer>SensePlex Limited</recorderManufacturer>
            <gain>1</gain>
            <maxClockDrift>0</maxClockDrift>
            <decimation sampleRateNumerator="{{.SampleRate}}" sampleRateDenominator="1" />
        </datalogger>
        <responsePAZ publicID="ResponsePAZ-E-C121G-VEL" name="AS-E-C121G-VEL">
            <type>A</type>
            <!--
                Gain = Sensor sensitivity = 100 V/m/s
            -->
            <gain>100</gain>
            <gainFrequency>4.5</gainFrequency>
            <normalizationFactor>171.99852139050935</normalizationFactor>
            <normalizationFrequency>4.5</normalizationFrequency>
            <numberOfZeros>2</numberOfZeros>
            <numberOfPoles>2</numberOfPoles>
            <zeros>(0,0) (0,0)</zeros>
            <poles>(-2.2211060060879837,2.221776881419294) (-2.2211060060879837,-2.221776881419294)</poles>
        </responsePAZ>
        <responsePAZ publicID="ResponsePAZ-E-C121G-ACC" name="AS-E-C121G-ACC">
            <type>D</type>
            <!--
                1 g = 32768 / 2 = 16384 counts ≈ 9.80665 m/s²
                Measurement range: ±2 g = ±32768 counts
                Gain = 16384 / 1 g
                     = 16384 / 9.80665
                     ≈ 1670.7030433430377
            -->
            <gain>1670.7030433430377</gain>
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
                    <stream publicID="Stream/E-C121G-CH1" code="{{.ChannelCode1}}"
                        datalogger="Datalogger-E-C121G-VEL"
                        sensor="Sensor-E-C121G-VEL">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>0</dataloggerChannel>
                        <sensorChannel>0</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>-90</dip>
                        <gain>171798691760</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C121G-CH2" code="{{.ChannelCode2}}"
                        datalogger="Datalogger-E-C121G-VEL"
                        sensor="Sensor-E-C121G-VEL">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>1</dataloggerChannel>
                        <sensorChannel>1</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>90</azimuth>
                        <dip>0</dip>
                        <gain>171798691760</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C121G-CH3" code="{{.ChannelCode3}}"
                        datalogger="Datalogger-E-C121G-VEL"
                        sensor="Sensor-E-C121G-VEL">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>2</dataloggerChannel>
                        <sensorChannel>2</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>0</dip>
                        <gain>171798691760</gain>
                        <gainFrequency>4.5</gainFrequency>
                        <gainUnit>m/s</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C121G-CH4" code="{{.ChannelCode4}}"
                        datalogger="Datalogger-E-C121G-ACC"
                        sensor="Sensor-E-C121G-ACC">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>3</dataloggerChannel>
                        <sensorChannel>3</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>-90</dip>
                        <gain>1670.7030433430377</gain>
                        <gainUnit>m/s**2</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C121G-CH5" code="{{.ChannelCode5}}"
                        datalogger="Datalogger-E-C121G-ACC"
                        sensor="Sensor-E-C121G-ACC">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>4</dataloggerChannel>
                        <sensorChannel>4</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>90</azimuth>
                        <dip>0</dip>
                        <gain>1670.7030433430377</gain>
                        <gainUnit>m/s**2</gainUnit>
                    </stream>
                    <stream publicID="Stream/E-C121G-CH6" code="{{.ChannelCode6}}"
                        datalogger="Datalogger-E-C121G-ACC"
                        sensor="Sensor-E-C121G-ACC">
                        <!-- <start>{{.StartTime}}</start> -->
                        <start>1970-01-01T00:00:00.0000Z</start>
                        <dataloggerChannel>5</dataloggerChannel>
                        <sensorChannel>5</sensorChannel>
                        <sampleRateNumerator>{{.SampleRate}}</sampleRateNumerator>
                        <sampleRateDenominator>1</sampleRateDenominator>
                        <depth>0</depth>
                        <azimuth>0</azimuth>
                        <dip>0</dip>
                        <gain>1670.7030433430377</gain>
                        <gainUnit>m/s**2</gainUnit>
                    </stream>
                </sensorLocation>
            </station>
        </network>
    </Inventory>
</seiscomp>