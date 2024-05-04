import { i18n } from "i18next";

import { i18nConfig } from "../../config/i18n";

export const getCurrentLocale = async (i18n: Promise<i18n>) => {
	const currentLang = (await i18n).language;
	const availableLocales = Object.keys(i18nConfig.resources);
	return availableLocales.includes(currentLang) ? currentLang : i18nConfig.fallback;
};
