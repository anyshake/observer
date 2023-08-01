import { LabelProps } from "../../components/Label";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import getAcceleration from "../../helpers/getAcceleration";
import getIntensity, {
    IntensityScaleStandard,
} from "../../helpers/getIntensity";
import getSortedArray from "../../helpers/getSortedArray";
import getVelocity from "../../helpers/getVelocity";
import getVoltage from "../../helpers/getVoltage";
import setObject from "../../helpers/setObject";

const setLabels = (
    obj: LabelProps[],
    data: any,
    adc: ADC,
    gp: Geophone,
    scale: IntensityScaleStandard
): LabelProps[] => {
    const { ehz, ehe, ehn } = gp;
    const { fullscale, resolution } = adc;

    const ehzVoltage: number[] = [];
    const eheVoltage: number[] = [];
    const ehnVoltage: number[] = [];

    const ehzVelocity: number[] = [];
    const eheVelocity: number[] = [];
    const ehnVelocity: number[] = [];

    const ehzAcceleration: number[] = [];
    const eheAcceleration: number[] = [];
    const ehnAcceleration: number[] = [];

    const sortedData: any = getSortedArray(data, "ts", "asc");
    const durationMS = sortedData[sortedData.length - 1].ts - sortedData[0].ts;

    let sampleRate = 1000 / 100;
    for (let i of data) {
        ehzVoltage.push(...getVoltage(i["ehz"], resolution, fullscale));
        eheVoltage.push(...getVoltage(i["ehe"], resolution, fullscale));
        ehnVoltage.push(...getVoltage(i["ehn"], resolution, fullscale));
        sampleRate += (i["ehz"].length + i["ehe"].length + i["ehn"].length) / 3;
    }
    sampleRate /= durationMS / 1000;

    ehzVelocity.push(...getVelocity(ehzVoltage, ehz));
    eheVelocity.push(...getVelocity(eheVoltage, ehe));
    ehnVelocity.push(...getVelocity(ehnVoltage, ehn));

    const interval = 1 / sampleRate;
    ehzAcceleration.push(...getAcceleration(ehzVelocity, interval));
    eheAcceleration.push(...getAcceleration(eheVelocity, interval));
    ehnAcceleration.push(...getAcceleration(ehnVelocity, interval));

    const ehzPGV = ehzVelocity.reduce((a, b) => {
        const absA = Math.abs(a);
        const absB = Math.abs(b);
        return Math.max(absA, absB);
    }, 0);
    const ehePGV = eheVelocity.reduce((a, b) => {
        const absA = Math.abs(a);
        const absB = Math.abs(b);
        return Math.max(absA, absB);
    }, 0);
    const ehnPGV = ehnVelocity.reduce((a, b) => {
        const absA = Math.abs(a);
        const absB = Math.abs(b);
        return Math.max(absA, absB);
    }, 0);

    const ehzPGA = ehzAcceleration.reduce((a, b) => {
        const absA = Math.abs(a);
        const absB = Math.abs(b);
        return Math.max(absA, absB);
    }, 0);
    const ehePGA = eheAcceleration.reduce((a, b) => {
        const absA = Math.abs(a);
        const absB = Math.abs(b);
        return Math.max(absA, absB);
    }, 0);
    const ehnPGA = ehnAcceleration.reduce((a, b) => {
        const absA = Math.abs(a);
        const absB = Math.abs(b);
        return Math.max(absA, absB);
    }, 0);

    setObject(obj, "[tag:ehz-pga]>value", ehzPGA.toFixed(2));
    setObject(obj, "[tag:ehz-intensity]>unit", scale);
    setObject(
        obj,
        "[tag:ehz-intensity]>value",
        getIntensity(ehzPGV, ehzPGA, scale)
    );

    setObject(obj, "[tag:ehe-pga]>value", ehePGA.toFixed(2));
    setObject(obj, "[tag:ehe-intensity]>unit", scale);
    setObject(
        obj,
        "[tag:ehe-intensity]>value",
        getIntensity(ehePGV, ehePGA, scale)
    );

    setObject(obj, "[tag:ehn-pga]>value", ehnPGA.toFixed(2));
    setObject(obj, "[tag:ehn-intensity]>unit", scale);
    setObject(
        obj,
        "[tag:ehn-intensity]>value",
        getIntensity(ehnPGV, ehnPGA, scale)
    );

    return obj;
};

export default setLabels;
