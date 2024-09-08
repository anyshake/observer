import { JSX, lazy, LazyExoticComponent, RefObject } from "react";

import { RouterMode } from "../components/RouterWrapper";
import { i18nConfig } from "./i18n";

const Home = lazy(() => import("../views/Home"));
const Realtime = lazy(() => import("../views/Realtime"));
const History = lazy(() => import("../views/History"));
const Export = lazy(() => import("../views/Export"));
const Setting = lazy(() => import("../views/Setting"));
const NotFound = lazy(() => import("../views/NotFound"));

export type RouterProp<T> = Record<string, T>;

export interface RouterComponentProps {
	refs?: RouterProp<RefObject<HTMLElement>>;
	locale?: string;
}

export interface RouterConfigRoutes {
	readonly uri: string;
	readonly prefix: string;
	readonly suffix: string;
	readonly element: LazyExoticComponent<(props: RouterComponentProps) => JSX.Element>;
	readonly title: Record<keyof typeof i18nConfig.resources, string>;
}

type RouterConfig = {
	readonly mode: RouterMode;
	readonly basename: string;
	readonly routes: Record<string, RouterConfigRoutes>;
};

export const routerConfig: RouterConfig = {
	basename: "/",
	mode: "hash",
	routes: {
		home: {
			prefix: "/",
			uri: "",
			suffix: "",
			element: Home,
			title: {
				"en-US": "Station Status",
				"zh-TW": "當前站況",
				"zh-CN": "测站状态"
			}
		},
		realtime: {
			prefix: "/realtime",
			uri: "",
			suffix: "",
			element: Realtime,
			title: {
				"en-US": "Realtime Waveform",
				"zh-TW": "即時波形",
				"zh-CN": "实时波形"
			}
		},
		history: {
			prefix: "/history",
			uri: "",
			suffix: "",
			element: History,
			title: {
				"en-US": "History Waveform",
				"zh-TW": "歷史資料",
				"zh-CN": "历史数据"
			}
		},
		export: {
			prefix: "/export",
			uri: "",
			suffix: "",
			element: Export,
			title: {
				"en-US": "Waveform Export",
				"zh-TW": "波形匯出",
				"zh-CN": "波形导出"
			}
		},
		setting: {
			prefix: "/setting",
			uri: "",
			suffix: "",
			element: Setting,
			title: {
				"en-US": "Panel Settings",
				"zh-TW": "面板設定",
				"zh-CN": "面板设置"
			}
		},
		default: {
			prefix: "*",
			uri: "",
			suffix: "",
			element: NotFound,
			title: {
				"en-US": "Not Found",
				"zh-TW": "找不到頁面",
				"zh-CN": "找不到页面"
			}
		}
	}
};
