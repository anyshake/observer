import { LabelProps } from "../../components/Label";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import getAcceleration from "../../helpers/seismic/getAcceleration";
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";
import getSortedArray from "../../helpers/array/getSortedArray";
import getVelocity from "../../helpers/seismic/getVelocity";
import getVoltage from "../../helpers/seismic/getVoltage";
import setObject from "../../helpers/utils/setObjectByPath";
import GLOBAL_CONFIG from "../../config/global";

const setLabels = (
    obj: LabelProps[],
    data: any,
    adc: ADC,
    gp: Geophone,
    scale: IntensityStandardProperty
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

    // Match scale standard
    const { scales } = GLOBAL_CONFIG.app_settings;
    const scaleStandard = scales.find(
        (item) => item.property().value === scale.value
    );

    setObject(obj, "[tag:ehz-pga]>value", ehzPGA.toFixed(2));
    setObject(
        obj,
        "[tag:ehz-intensity]>value",
        `${scale.value} 震度 ${scaleStandard?.intensity(ehzPGV, ehzPGA)}`
    );

    setObject(obj, "[tag:ehe-pga]>value", ehePGA.toFixed(2));
    setObject(
        obj,
        "[tag:ehe-intensity]>value",
        `${scale.value} 震度 ${scaleStandard?.intensity(ehePGV, ehePGA)}`
    );

    setObject(obj, "[tag:ehn-pga]>value", ehnPGA.toFixed(2));
    setObject(
        obj,
        "[tag:ehn-intensity]>value",
        `${scale.value} 震度 ${scaleStandard?.intensity(ehnPGV, ehnPGA)}`
    );

    return obj;
};

export default setLabels;
