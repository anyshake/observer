import { BannerProps } from "../../components/Banner";
import { IntensityScaleStandard } from "../../helpers/getIntensity";
import getTimeString from "../../helpers/getTimeString";

const setBanner = (
    message: any,
    prevTs: number,
    scale: IntensityScaleStandard
): BannerProps => {
    const { ts, ehz, ehe, ehn } = message;
    const sampleLength = (ehz.length + ehe.length + ehn.length) / 3;
    const sampleRate = (1000 * sampleLength) / (ts - prevTs);

    return {
        type: "success",
        label: `连接成功：实际采样率 ${sampleRate.toFixed(2)} Sps`,
        text: `当前震度标准为 ${scale}\n数据最后更新于 ${getTimeString(ts)}`,
    };
};

export default setBanner;
