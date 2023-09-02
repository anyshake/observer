import { BannerProps } from "../../components/Banner";
import { IntensityStandardProperty } from "../../helpers/seismic/intensityStandard";
import getTimeString from "../../helpers/utils/getTimeString";

const setBanner = (
    message?: any,
    prevTs?: number,
    scale?: IntensityStandardProperty
): BannerProps => {
    const { ts, ehz, ehe, ehn } = message || { ts: -1, ehz: 0, ehe: 0, ehn: 0 };
    const sampleLength = (ehz.length + ehe.length + ehn.length) / 3;
    const sampleRate = (1000 * sampleLength) / (ts - (prevTs || 0));

    // Display error message if ts is -1 and all data are 0
    if (ts === -1 && ehz * ehn * ehe === 0) {
        return {
            type: "error",
            label: "连线断开",
            text: "正在尝试重新连线，请稍候",
        };
    }

    return {
        type: "success",
        label: `连接成功：实际采样率 ${sampleRate.toFixed(2)} Sps`,
        text: `当前震度标准为${
            scale?.name || "未知标准"
        }\n数据最后更新于 ${getTimeString(ts)}`,
    };
};

export default setBanner;
