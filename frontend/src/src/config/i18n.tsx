import { ResourceLanguage } from "i18next";

import { createI18n } from "../helpers/i18n/createI18n";
import enUS from "../locales/en-US.json";
import zhCN from "../locales/zh-CN.json";
import zhTW from "../locales/zh-TW.json";

interface I18nConfig {
	key: string;
	fallback: string;
	resources: Record<string, { label: string; translation: ResourceLanguage }>;
}

export const i18nConfig: I18nConfig = {
	key: "i18n",
	fallback: "en-US",
	resources: {
		"en-US": { label: "US English", translation: enUS },
		"zh-TW": { label: "正體中文", translation: zhTW },
		"zh-CN": { label: "简体中文", translation: zhCN }
	}
};

const i18n = createI18n(i18nConfig.fallback, i18nConfig.key, i18nConfig.resources);

export default i18n;
