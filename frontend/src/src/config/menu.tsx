import IndexIcon from "../assets/icons/server-solid.svg";
import RealtimeIcon from "../assets/icons/wave-square-solid.svg";
import HistoryIcon from "../assets/icons/file-waveform-solid.svg";

const MENU_CONFIG: MenuConfig[] = [
    {
        uri: "/",
        label: "测站状态",
        icon: IndexIcon,
    },
    {
        uri: "/realtime",
        label: "实时波形",
        icon: RealtimeIcon,
    },
    {
        uri: "/history",
        label: "历史数据",
        icon: HistoryIcon,
    },
    // {
    //     uri: "/settings",
    //     label: "系统设定",
    //     icon: SettingsIcon,
    // },
];

export interface MenuConfig {
    readonly uri: string;
    readonly label: string;
    readonly icon: string;
}

export default MENU_CONFIG;
