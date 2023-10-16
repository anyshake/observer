import { BannerProps } from "../../components/Banner";
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";
import getTimeString from "../../helpers/utils/getTimeString";

const setBanner = (
    message?: any,
    prevTs?: number,
    scale?: IntensityStandardProperty
): BannerProps => {
    // Parse response, empty response means error
    const { ts, ehz, ehe, ehn } = message || { ts: -1, ehz: 0, ehe: 0, ehn: 0 };
    const sampleLength = (ehz.length + ehe.length + ehn.length) / 3;
    const sampleRate = ((1000 * sampleLength) / (ts - (prevTs || 0))).toFixed(
        2
    );

    // Display error message if ts is -1 and all data are 0
    if (ts === -1 && ehz * ehn * ehe === 0) {
        return {
            type: "error",
            label: { id: "views.realtime.banner.error.label" },
            text: { id: "views.realtime.banner.error.text" },
        };
    }

    return {
        type: "success",
        label: {
            id: "views.realtime.banner.success.label",
            format: { sampleRate },
        },
        text: {
            id: "views.realtime.banner.success.text",
            format: {
                time: getTimeString(ts),
                scale: `${scale?.value} - ${scale?.name}` || "Unknown",
            },
        },
    };
};

export default setBanner;
