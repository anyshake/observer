export interface IntensityStandardProperty {
    readonly name: string;
    readonly value: string;
}

export interface IntensityStandard {
    property(): IntensityStandardProperty;
    intensity(pgv: number, pga: number): string;
}

class JMAIntensityStandard implements IntensityStandard {
    property(): IntensityStandardProperty {
        return {
            name: "日本気象庁震度",
            value: "JMA",
        };
    }

    intensity(_: number, pga: number): string {
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
    }
}

class CWBIntensityStandard implements IntensityStandard {
    property(): IntensityStandardProperty {
        return {
            name: "台湾中央氣象局新震度",
            value: "CWB",
        };
    }

    intensity(pgv: number, pga: number): string {
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
    }
}

class MMIIntensityStandard implements IntensityStandard {
    property(): IntensityStandardProperty {
        return {
            name: "修訂麥加利地震震度表",
            value: "MMI",
        };
    }

    intensity(pgv: number, pga: number): string {
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
    }
}

class CSISIntensityStandard implements IntensityStandard {
    property(): IntensityStandardProperty {
        return {
            name: "中国地震局地震震度",
            value: "CSIS",
        };
    }

    intensity(pgv: number, pga: number): string {
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

        return `${II.toFixed(1)} 级`;
    }
}

export {
    JMAIntensityStandard,
    CWBIntensityStandard,
    MMIIntensityStandard,
    CSISIntensityStandard,
};
