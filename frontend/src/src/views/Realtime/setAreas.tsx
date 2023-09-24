import getVelocity from "../../helpers/seismic/getVelocity";
import getQueueArray from "../../helpers/array/getQueueArray";
import getVoltage from "../../helpers/seismic/getVoltage";
import setObjectByPath from "../../helpers/utils/setObjectByPath";
import { RealtimeArea } from "./Index";
import getAcceleration from "../../helpers/seismic/getAcceleration";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import GLOBAL_CONFIG from "../../config/global";
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";

const setAreas = (
    obj: RealtimeArea[],
    data: any,
    prevTs: number,
    length: number,
    adc: ADC,
    gp: Geophone,
    scale: IntensityStandardProperty
): RealtimeArea[] => {
    const tags = ["ehz", "ehe", "ehn"];
    const { ts } = data;

    for (let i of tags) {
        // Get sample rate
        const { ehz, ehe, ehn } = data;
        let sampleRate = (ehz.length + ehe.length + ehn.length) / 3;

        // Get time difference and time span
        const timeDiff = prevTs !== 0 ? ts - prevTs : 1000;
        const timeSpan = timeDiff / sampleRate;

        // Get voltage, velocity, acceleration
        const channelData = data[i];
        const voltage = getVoltage(channelData, adc.resolution, adc.fullscale);
        const velocity = getVelocity(voltage, gp[i]);
        const acceleration = getAcceleration(velocity, timeSpan / timeDiff);

        // Fill data queue with raw count
        const srcArr = obj.find((item) => item.tag === i)?.chart.series?.data;
        const newArr = [];
        for (let j = 0; j < channelData.length; j++) {
            newArr.push([ts - (sampleRate - j) * timeSpan, channelData[j]]);
        }

        // Merge data queue with raw count
        const resultArr = getQueueArray(srcArr, newArr, length * sampleRate);
        setObjectByPath(obj, `[tag:${i}]>chart>series>data`, resultArr);

        // Get PGV, PGA
        const pgv = velocity.reduce((a, b) => {
            const absA = Math.abs(a);
            const absB = Math.abs(b);
            return Math.max(absA, absB);
        }, 0);
        const pga = acceleration.reduce((a, b) => {
            const absA = Math.abs(a);
            const absB = Math.abs(b);
            return Math.max(absA, absB);
        }, 0);

        // Match scale standard
        const { scales } = GLOBAL_CONFIG.app_settings;
        const scaleStandard = scales.find(
            (item) => item.property().value === scale.value
        );
        // Get intensity
        const intensity = scaleStandard?.intensity(pgv, pga);

        // Set PGV, PGA, intensity in area field
        setObjectByPath(obj, `[tag:${i}]>area>text`, {
            id: `views.realtime.areas.${i}.text`,
            format: {
                pga: pga.toFixed(5),
                pgv: pgv.toFixed(5),
                intensity: `${scale.value} ${intensity}`,
            },
        });
    }

    return obj;
};

export default setAreas;
