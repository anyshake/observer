import IndexIcon from "../assets/icons/server-solid.svg";
import RealtimeIcon from "../assets/icons/wave-square-solid.svg";
import HistoryIcon from "../assets/icons/file-waveform-solid.svg";
import ExportIcon from "../assets/icons/cloud-arrow-down-solid.svg";
import SettingIcon from "../assets/icons/gear-solid.svg";

const MENU_CONFIG: MenuConfig = {
    title: "config.menu.title",
    list: [
        {
            uri: "/",
            label: "config.menu.list.index",
            icon: IndexIcon,
        },
        {
            uri: "/realtime",
            label: "config.menu.list.realtime",
            icon: RealtimeIcon,
        },
        {
            uri: "/history",
            label: "config.menu.list.history",
            icon: HistoryIcon,
        },
        {
            uri: "/export",
            label: "config.menu.list.export",
            icon: ExportIcon,
        },
        {
            uri: "/setting",
            label: "config.menu.list.setting",
            icon: SettingIcon,
        },
    ],
};

export interface MenuConfigItem {
    readonly uri: string;
    readonly label: string;
    readonly icon: string;
}

export interface MenuConfig {
    readonly title: string;
    readonly list: MenuConfigItem[];
}

export default MENU_CONFIG;
