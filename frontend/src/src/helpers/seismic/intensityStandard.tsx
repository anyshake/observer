import { ADC } from "../../stores/adc";
import { Geophone } from "../../stores/geophone";

interface Property {
    readonly name: string;
    readonly value: string;
}

interface GetIntensityData {
    rawData: number[];
    currentPGA: number;
    currentPGV: number;
}

interface GetIntensityAttributes {
    adc: ADC;
    geophone: Geophone;
}

export interface IntensityStandard {
    property: () => Property;
    getIntensity: (
        data: GetIntensityData,
        attributes: GetIntensityAttributes
    ) => string;
}

export class JMAIntensityStandard implements IntensityStandard {
    property = () => ({
        name: "気象庁震度階級",
        value: "JMA",
    });

    getIntensity = (
        { currentPGA: pga }: GetIntensityData,
        _attributes: GetIntensityAttributes
    ) => {
        let intensity = parseFloat((2 * Math.log10(pga) + 0.94).toFixed(3));
        intensity = parseFloat(intensity.toFixed(2));
        switch (true) {
            case intensity < 0.5:
                return "0";
            case intensity < 1.5:
                return "1";
            case intensity < 2.5:
                return "2";
            case intensity < 3.5:
                return "3";
            case intensity < 4.5:
                return "4";
            case intensity < 5.0:
                return "5 弱";
            case intensity < 5.5:
                return "5 強";
            case intensity < 6.0:
                return "6 弱";
            case intensity < 6.5:
                return "6 強";
            default:
                return "7";
        }
    };
}

export class CWAIntensityStandard implements IntensityStandard {
    property = (): Property => ({
        name: "交通部中央氣象署地震震度分級",
        value: "CWA",
    });

    getIntensity = (
        { currentPGA: pga, currentPGV: pgv }: GetIntensityData,
        _attributes: GetIntensityAttributes
    ) => {
        if (pga < 80) {
            switch (true) {
                case pga < 0.8:
                    return "0 級";
                case pga < 2.5:
                    return "1 級";
                case pga < 8:
                    return "2 級";
                case pga < 25:
                    return "3 級";
                case pga < 80:
                    return "4 級";
            }
        } else {
            switch (true) {
                case pgv < 15:
                    return "4 級";
                case pgv < 30:
                    return "5 弱";
                case pgv < 50:
                    return "5 強";
                case pgv < 80:
                    return "6 弱";
                case pgv < 140:
                    return "6 強";
            }
        }

        return "7 級";
    };
}

export class MMIIntensityStandard implements IntensityStandard {
    property = (): Property => ({
        name: "The Modified Mercalli Intensity",
        value: "MMI",
    });

    getIntensity = (
        { currentPGA: pga }: GetIntensityData,
        _attributes: GetIntensityAttributes
    ) => (2.33 * Math.log10(pga) + 1.5).toFixed(0);
}

export class CSISIntensityStandard implements IntensityStandard {
    property = (): Property => ({
        name: "中国地震烈度表",
        value: "CSIS",
    });

    getIntensity = (
        { currentPGA: pga, currentPGV: pgv }: GetIntensityData,
        _attributes: GetIntensityAttributes
    ) => {
        const IA = 3.17 * Math.log10(pga) + 6.59;
        const IV = 3 * Math.log10(pgv) + 9.77;

        let II = 1;
        if (IA >= 6 && IV >= 6) {
            II = IV;
        } else {
            II = (IA + IV) / 2;
        }

        if (II < 1) {
            II = 1;
        } else if (II > 12) {
            II = 12;
        }

        return `${II.toFixed(0)}`;
    };
}
