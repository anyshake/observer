import createI18n from "../helpers/i18n/createI18n";

const I18N_CONFIG: I18nConfig = {
    fallback: "zh-CN",
    key: "i18n",
    list: [
        { name: "简体中文", value: "zh-CN" },
        { name: "正體中文", value: "zh-TW" },
        { name: "English", value: "en-US" },
    ],
    uri: "/i18n/{{lng}}.json",
};

export interface I18nTranslation {
    readonly id: string;
    readonly format?: { [key: string]: string };
}

export interface I18nConfigItem {
    readonly name: string;
    readonly value: string;
}

export interface I18nConfig {
    readonly key: string;
    readonly uri: string;
    readonly fallback: string;
    readonly list: I18nConfigItem[];
}

export default I18N_CONFIG;
export const i18n = createI18n(
    I18N_CONFIG.uri,
    I18N_CONFIG.fallback,
    I18N_CONFIG.key
);
