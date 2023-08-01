export type IntensityScaleStandard = "JMA" | "CWB" | "MMI" | "CSIS";

const getJMAIntensity = (pga: number): string => {
    const intensity = Math.round(
        Math.round(2 * Math.log10(Math.abs(pga)) + 0.94)
    );

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

const getCWBIntensity = (pgv: number, pga: number): string => {
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

const getMMIIntensity = (pgv: number, pga: number): string => {
    switch (true) {
        case pga < 1:
            return "1 度";
        case pga < 2.1:
            return "2 度";
        case pga < 5:
            return "3 度";
        case pga < 10:
            return "4 度";
        case pga < 21:
            return "5 度";
        case pga < 44:
            return "6 度";
        case pga < 94:
            return "7 度";
        case pga < 202:
            return "8 度";
        case pga < 432:
            return "9 度";
    }

    if (pga >= 432 && pgv > 116) {
        return "10 度";
    } else if (pga >= 432) {
        return "11 度";
    }

    return "12 度";
};

const getCSISIntensity = (pgv: number, pga: number): string => {
    const IA = 3.17 * Math.log(pga) + 6.59;
    const IV = 3 * Math.log(pgv) + 9.77;

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

    return `${II.toFixed(1)} 级`;
};

const getIntensity = (
    pgv: number,
    pga: number,
    scale: IntensityScaleStandard
): string => {
    switch (scale) {
        case "MMI":
            return getMMIIntensity(pgv, pga);
        case "CSIS":
            return getCSISIntensity(pgv, pga);
        case "CWB":
            return getCWBIntensity(pgv, pga);
        case "JMA":
            return getJMAIntensity(pga);
        default:
            return "未知";
    }
};

export default getIntensity;
export { getJMAIntensity, getCWBIntensity, getMMIIntensity, getCSISIntensity };
