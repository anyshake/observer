import getBackend from "../helpers/app/getBackend";
import {
    CSISIntensityStandard,
    CWBIntensityStandard,
    IntensityStandard,
    JMAIntensityStandard,
    MMIIntensityStandard,
} from "../helpers/seismic/intensityStandard";
import getRelease from "../helpers/app/getRelease";
import getVersion from "../helpers/app/getVersion";

const version = getVersion();
const release = getRelease();
const backend = getBackend();
const scales = [
    new JMAIntensityStandard(),
    new CWBIntensityStandard(),
    new MMIIntensityStandard(),
    new CSISIntensityStandard(),
];

const GLOBAL_CONFIG: GlobalConfig = {
    app_settings: {
        version,
        release,
        scales,
        router: "hash",
        name: "Observer",
        author: "Project ES",
        title: "Observer 测站面板",
        description: "Constructing Realtime Seismic Network Ambitiously.",
    },
    api_settings: {
        backend,
        version: "v1",
        prefix: "/api",
        types: ["http", "ws"],
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
    readonly scales: IntensityStandard[];
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
export const fallbackScale = scales[0];
