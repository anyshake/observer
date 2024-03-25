import {
    CSISIntensityStandard,
    CWAIntensityStandard,
    IntensityStandard,
    JMAIntensityStandard,
    MMIIntensityStandard,
} from "../helpers/seismic/intensityStandard";
import { getRelease } from "../helpers/app/getRelease";
import { getVersion } from "../helpers/app/getVersion";
import { i18nConfig } from "./i18n";

interface GlobalConfig {
    readonly name: string;
    readonly title: string;
    readonly author: string;
    readonly version: string;
    readonly release: string;
    readonly homepage: string;
    readonly repository: string;
    readonly duration: {
        readonly default: number;
        readonly maximum: number;
        readonly minimum: number;
    };
    readonly retention: {
        readonly default: number;
        readonly maximum: number;
        readonly minimum: number;
    };
    readonly scales: IntensityStandard[];
    readonly footer: Record<keyof typeof i18nConfig.resources, string>;
}

const version = getVersion();
const release = getRelease();
const scales = [
    new JMAIntensityStandard(),
    new CWAIntensityStandard(),
    new MMIIntensityStandard(),
    new CSISIntensityStandard(),
];

export const globalConfig: GlobalConfig = {
    scales,
    version,
    release,
    name: "Observer",
    author: "AnyShake",
    title: "AnyShake Observer",
    homepage: "https://anyshake.org",
    repository: "https://github.com/AnyShake",
    duration: { default: 300, maximum: 3600, minimum: 10 },
    retention: { default: 180, maximum: 600, minimum: 10 },
    footer: {
        "en-US": "Constructing Realtime Seismic Network Ambitiously.",
        "zh-CN": "雄心勃勃，致力于构建实时地震网络",
        "zh-TW": "雄心勃勃，致力於建置即時地震網路",
    },
};

export const fallbackScale = scales[0];
