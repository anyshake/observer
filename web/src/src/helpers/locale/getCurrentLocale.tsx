import { i18n } from 'i18next';

import { localeConfig } from '../../config/locale';

export const getCurrentLocale = async (i18n: Promise<i18n>) => {
    const currentLang = (await i18n).language;
    const availableLocales = Object.keys(localeConfig.resources);
    return availableLocales.includes(currentLang) ? currentLang : localeConfig.fallback;
};
