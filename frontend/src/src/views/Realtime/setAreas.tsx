import getVelocity from "../../helpers/getVelocity";
import getQueueArray from "../../helpers/getQueueArray";
import getVoltage from "../../helpers/getVoltage";
import setObject from "../../helpers/setObject";
import { RealtimeArea } from "./Index";
import getAcceleration from "../../helpers/getAcceleration";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import getIntensity, {
    IntensityScaleStandard,
} from "../../helpers/getIntensity";

const setAreas = (
    obj: RealtimeArea[],
    data: any,
    prevTs: number,
    length: number,
    adc: ADC,
    gp: Geophone,
    scale: IntensityScaleStandard
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
        const voltage = getVoltage(data[i], adc.resolution, adc.fullscale);
        const velocity = getVelocity(voltage, gp[i]);
        const acceleration = getAcceleration(velocity, timeSpan / timeDiff);

        // Fill data queue with acceleration
        const srcArr = obj.find((item) => item.tag === i)?.chart.series?.data;
        const newArr = [];
        for (let j = 0; j < acceleration.length; j++) {
            newArr.push([ts - (sampleRate - j) * timeSpan, acceleration[j]]);
        }

        // remove old data and add new data
        const resultArr = getQueueArray(srcArr, newArr, length * sampleRate);
        setObject(obj, `[tag:${i}]>chart>series>data`, resultArr);

        // Get PGV, PGA, Intensity
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
        const intensity = getIntensity(pgv, pga, scale);

        // Set PGV, PGA, Intensity in area
        setObject(
            obj,
            `[tag:${i}]>area>text`,
            `PGA：${pga.toFixed(5)} gal\n
            PGV：${pgv.toFixed(5)} cm/s\n
            震度：${scale} 震度 ${intensity}`
        );
    }

    return obj;
};

export default setAreas;
