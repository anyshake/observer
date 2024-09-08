import { mdiCog, mdiDatabaseExport, mdiFileClock, mdiServerNetwork, mdiWaveform } from "@mdi/js";

import { i18nConfig } from "./i18n";

export interface MenuItem {
	readonly url: string;
	readonly icon: string;
	readonly label: Record<keyof typeof i18nConfig.resources, string>;
}

export const menuConfig: MenuItem[] = [
	{
		url: "/",
		label: {
			"en-US": "Station Status",
			"zh-TW": "當前站況",
			"zh-CN": "测站状态"
		},
		icon: mdiServerNetwork
	},
	{
		url: "/realtime",
		label: {
			"en-US": "Realtime Waveform",
			"zh-TW": "即時波形",
			"zh-CN": "实时波形"
		},
		icon: mdiWaveform
	},
	{
		url: "/history",
		label: {
			"en-US": "History Waveform",
			"zh-TW": "歷史資料",
			"zh-CN": "历史数据"
		},
		icon: mdiFileClock
	},
	{
		url: "/export",
		label: {
			"en-US": "Waveform Export",
			"zh-TW": "波形匯出",
			"zh-CN": "波形导出"
		},
		icon: mdiDatabaseExport
	},
	{
		url: "/setting",
		label: {
			"en-US": "Panel Settings",
			"zh-TW": "面板設定",
			"zh-CN": "面板设置"
		},
		icon: mdiCog
	}
];
