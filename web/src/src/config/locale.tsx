import { ResourceLanguage } from 'i18next';

import { createI18n } from '../helpers/locale/createI18n';
import deDE from '../locales/de-DE.json';
import enUS from '../locales/en-US.json';
import frFR from '../locales/fr-FR.json';
import idID from '../locales/id-ID.json';
import jaJP from '../locales/ja-JP.json';
import plPL from '../locales/pl-PL.json';
import ptPT from '../locales/pt-PT.json';
import ruRU from '../locales/ru-RU.json';
import trTR from '../locales/tr-TR.json';
import zhTW from '../locales/zh-TW.json';

interface ILocaleConfig {
    key: string;
    fallback: string;
    resources: Record<string, { label: string; translation: ResourceLanguage }>;
}

export const localeConfig: ILocaleConfig = {
    fallback: 'en-US',
    key: 'i18n',
    resources: {
        'de-DE': { label: 'Deutsch', translation: deDE },
        'en-US': { label: 'English', translation: enUS },
        'fr-FR': { label: 'Français', translation: frFR },
        'id-ID': { label: 'Bahasa Indonesia', translation: idID },
        'ja-JP': { label: '日本語', translation: jaJP },
        'pl-PL': { label: 'Polski', translation: plPL },
        'pt-PT': { label: 'Português', translation: ptPT },
        'ru-RU': { label: 'Русский', translation: ruRU },
        'tr-TR': { label: 'Türkçe', translation: trTR },
        'zh-TW': { label: '正體中文', translation: zhTW }
    }
};

export type Translation = Record<keyof typeof localeConfig.resources, string>;

const i18n = createI18n(localeConfig.fallback, localeConfig.key, localeConfig.resources);

export default i18n;
