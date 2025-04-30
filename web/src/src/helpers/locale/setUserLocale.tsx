import { i18n } from 'i18next';

import { localeConfig } from '../../config/locale';

export const setUserLocale = async (i18n: Promise<i18n>, lang: string) => {
    const availableLocales = Object.keys(localeConfig.resources);
    (await i18n).changeLanguage(availableLocales.includes(lang) ? lang : localeConfig.fallback);
};
