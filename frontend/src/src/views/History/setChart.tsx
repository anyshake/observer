import { ChartProps } from "../../components/Chart";
import { ADC } from "../../config/adc";
import { Geophone } from "../../config/geophone";
import getAcceleration from "../../helpers/getAcceleration";
import getSortedArray from "../../helpers/getSortedArray";
import getVelocity from "../../helpers/getVelocity";
import getVoltage from "../../helpers/getVoltage";
import setObject from "../../helpers/setObject";

const setChart = (
    obj: ChartProps,
    data: any,
    adc: ADC,
    gp: Geophone
): ChartProps => {
    const tag = ["ehz", "ehe", "ehn"];
    const { fullscale, resolution } = adc;

    // Sort data by timestamp
    const sortedData: any = getSortedArray(data, "ts", "asc");

    // Store channel acceleration
    const ehzAcceleration = [];
    const eheAcceleration = [];
    const ehnAcceleration = [];

    // Calculate acceleration
    let prevTs = 0;
    for (let i of sortedData) {
        for (let j of tag) {
            // Get data sample rate
            const sampleLength = i[j].length;
            let sampleRate = 0;
            if (prevTs !== 0) {
                sampleRate = (1000 * sampleLength) / (prevTs - i.ts);
            } else {
                sampleRate = i[j].length;
            }

            // Get time difference and time span
            const timeDiff = prevTs !== 0 ? prevTs - i.ts : 1000;
            const timeSpan = timeDiff / sampleRate;

            // Get voltage, velocity, acceleration
            const voltage = getVoltage(i[j], resolution, fullscale);
            const velocity = getVelocity(voltage, gp[j]);
            const acceleration = getAcceleration(velocity, timeSpan / timeDiff);

            // Store acceleration
            for (let k = 0; k < acceleration.length; k++) {
                switch (j) {
                    case "ehz":
                        ehzAcceleration.push([
                            i.ts + k * timeSpan,
                            acceleration[k],
                        ]);
                        break;
                    case "ehe":
                        eheAcceleration.push([
                            i.ts + k * timeSpan,
                            acceleration[k],
                        ]);
                        break;
                    case "ehn":
                        ehnAcceleration.push([
                            i.ts + k * timeSpan,
                            acceleration[k],
                        ]);
                        break;
                }
            }
        }

        prevTs = i.ts;
    }

    setObject(obj, `series>[name:EHZ]>data`, ehzAcceleration);
    setObject(obj, `series>[name:EHE]>data`, eheAcceleration);
    setObject(obj, `series>[name:EHN]>data`, ehnAcceleration);
    return obj;
};

export default setChart;
