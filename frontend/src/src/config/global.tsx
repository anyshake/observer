import getBackend from "../helpers/getBackend";
import { IntensityScaleStandard } from "../helpers/getIntensity";
import getRelease from "../helpers/getRelease";
import getVersion from "../helpers/getVersion";

const GLOBAL_CONFIG: GlobalConfig = {
    app_settings: {
        router: "hash",
        name: "G-Observer",
        author: "通信实验室",
        title: "G-Observer 测站面板",
        description: "Constructing Realtime Seismic Network Ambitiously.",
        version: getVersion(),
        release: getRelease(),
        scale: "JMA",
    },
    api_settings: {
        version: "v1",
        prefix: "/api",
        types: ["http", "ws"],
        backend: getBackend(),
    },
};

export interface AppSettings {
    readonly router: "hash" | "history";
    readonly name: string;
    readonly author: string;
    readonly title: string;
    readonly version: string;
    readonly release: string;
    readonly description: string;
    readonly scale: IntensityScaleStandard;
}

export interface ApiSettings {
    readonly prefix: string;
    readonly version: string;
    readonly backend: string;
    readonly types: string[];
}

export interface GlobalConfig {
    readonly app_settings: AppSettings;
    readonly api_settings: ApiSettings;
}

export default GLOBAL_CONFIG;
