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
    // Metadata for getting voltage & velocity
    const { ehz, ehe, ehn } = gp;
    const { fullscale, resolution } = adc;

    // Channel voltage arrays
    const ehzVoltage: number[] = [];
    const eheVoltage: number[] = [];
    const ehnVoltage: number[] = [];
    // Channel velocity arrays
    const ehzVelocity: number[] = [];
    const eheVelocity: number[] = [];
    const ehnVelocity: number[] = [];
    // Channel acceleration arrays
    const ehzAcceleration: number[] = [];
    const eheAcceleration: number[] = [];
    const ehnAcceleration: number[] = [];

    // Sort data by timestamp and get data duration in seconds
    const sortedData: any = getSortedArray(data, "ts", "asc");
    const durationSec =
        (sortedData[sortedData.length - 1].ts - sortedData[0].ts) / 1000;

    // Get voltage arrays and sample rate
    let sampleRate = 0;
    for (let i of data) {
        // getCounts offsets data to center around 0
        ehzVoltage.push(...getVoltage(i["ehz"], resolution, fullscale));
        eheVoltage.push(...getVoltage(i["ehe"], resolution, fullscale));
        ehnVoltage.push(...getVoltage(i["ehn"], resolution, fullscale));

        // Use average samples
        sampleRate += (i["ehz"].length + i["ehe"].length + i["ehn"].length) / 3;
    }

    // Get sample rate
    sampleRate /= durationSec;
    // Get channel velocities
    ehzVelocity.push(...getVelocity(ehzVoltage, ehz));
    eheVelocity.push(...getVelocity(eheVoltage, ehe));
    ehnVelocity.push(...getVelocity(ehnVoltage, ehn));

    // Get channel accelerations by differentiation
    const interval = 1 / sampleRate;
    ehzAcceleration.push(...getAcceleration(ehzVelocity, interval));
    eheAcceleration.push(...getAcceleration(eheVelocity, interval));
    ehnAcceleration.push(...getAcceleration(ehnVelocity, interval));

    // Get peak ground velocity and acceleration
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

    // Set labels data
    setObject(obj, "[tag:ehz-pga]>value", ehzPGA.toFixed(2));
    setObject(
        obj,
        "[tag:ehz-intensity]>value",
        `${scale.value} ${scaleStandard?.intensity(ehzPGV, ehzPGA)}`
    );
    setObject(obj, "[tag:ehe-pga]>value", ehePGA.toFixed(2));
    setObject(
        obj,
        "[tag:ehe-intensity]>value",
        `${scale.value} ${scaleStandard?.intensity(ehePGV, ehePGA)}`
    );
    setObject(obj, "[tag:ehn-pga]>value", ehnPGA.toFixed(2));
    setObject(
        obj,
        "[tag:ehn-intensity]>value",
        `${scale.value} ${scaleStandard?.intensity(ehnPGV, ehnPGA)}`
    );

    return obj;
};

export default setLabels;
