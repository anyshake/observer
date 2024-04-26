import { i18n } from "i18next";
import { i18nConfig } from "../../config/i18n";

export const setUserLocale = async (i18n: Promise<i18n>, lang: string) => {
    const availableLocales = Object.keys(i18nConfig.resources);
    (await i18n).changeLanguage(
        availableLocales.includes(lang) ? lang : i18nConfig.fallback
    );
};
