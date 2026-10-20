import { Resource, ResourceLanguage } from 'i18next';

import { createI18n } from '../helpers/locale/createI18n';
import arSA from '../locales/ar-SA.json';
import deDE from '../locales/de-DE.json';
import enUS from '../locales/en-US.json';
import esES from '../locales/es-ES.json';
import esMX from '../locales/es-MX.json';
import faIR from '../locales/fa-IR.json';
import frFR from '../locales/fr-FR.json';
import hiIN from '../locales/hi-IN.json';
import idID from '../locales/id-ID.json';
import itIT from '../locales/it-IT.json';
import jaJP from '../locales/ja-JP.json';
import koKR from '../locales/ko-KR.json';
import myMM from '../locales/my-MM.json';
import neNP from '../locales/ne-NP.json';
import nlNL from '../locales/nl-NL.json';
import plPL from '../locales/pl-PL.json';
import ptBR from '../locales/pt-BR.json';
import ptPT from '../locales/pt-PT.json';
import ruRU from '../locales/ru-RU.json';
import thTH from '../locales/th-TH.json';
import trTR from '../locales/tr-TR.json';
import ukUA from '../locales/uk-UA.json';
import viVN from '../locales/vi-VN.json';
import zhTW from '../locales/zh-TW.json';

// A type alias retains implicit index-signature compatibility with i18next's ResourceLanguage.
type ILocaleResource = {
    label: string;
    translation: ResourceLanguage;
    // Export name from date-fns/locale, used by the date picker adapter.
    key4datefns: keyof typeof import('date-fns/locale');
    // MUI packages support different locales; explicitly use enUS when unavailable.
    key4mui: {
        dataGrid: keyof typeof import('@mui/x-data-grid/locales');
        datePickers: Exclude<keyof typeof import('@mui/x-date-pickers/locales'), 'DEFAULT_LOCALE'>;
    };
};

interface ILocaleConfig {
    key: string;
    fallback: string;
    resources: Record<string, ILocaleResource>;
}

export const localeConfig: ILocaleConfig = {
    fallback: 'en-US',
    key: 'i18n',
    resources: {
        'ar-SA': {
            label: 'العربية',
            translation: arSA,
            key4datefns: 'enUS',
            // Only DataGrid provides an Arabic locale.
            key4mui: { dataGrid: 'arSD', datePickers: 'enUS' }
        },
        'de-DE': {
            label: 'Deutsch',
            translation: deDE,
            key4datefns: 'de',
            key4mui: { dataGrid: 'deDE', datePickers: 'deDE' }
        },
        'en-US': {
            label: 'English',
            translation: enUS,
            key4datefns: 'enUS',
            key4mui: { dataGrid: 'enUS', datePickers: 'enUS' }
        },
        'es-ES': {
            label: 'Español (España)',
            translation: esES,
            key4datefns: 'es',
            key4mui: { dataGrid: 'esES', datePickers: 'esES' }
        },
        'es-MX': {
            label: 'Español (México)',
            translation: esMX,
            key4datefns: 'es',
            key4mui: { dataGrid: 'esES', datePickers: 'esES' }
        },
        'fa-IR': {
            label: 'فارسی',
            translation: faIR,
            key4datefns: 'faIR',
            key4mui: { dataGrid: 'faIR', datePickers: 'faIR' }
        },
        'fr-FR': {
            label: 'Français',
            translation: frFR,
            key4datefns: 'fr',
            key4mui: { dataGrid: 'frFR', datePickers: 'frFR' }
        },
        'hi-IN': {
            label: 'हिन्दी',
            translation: hiIN,
            key4datefns: 'hi',
            // MUI does not provide a Hindi locale.
            key4mui: { dataGrid: 'enUS', datePickers: 'enUS' }
        },
        'id-ID': {
            label: 'Bahasa Indonesia',
            translation: idID,
            key4datefns: 'id',
            // MUI does not provide an Indonesian locale.
            key4mui: { dataGrid: 'enUS', datePickers: 'enUS' }
        },
        'it-IT': {
            label: 'Italiano',
            translation: itIT,
            key4datefns: 'it',
            key4mui: { dataGrid: 'itIT', datePickers: 'itIT' }
        },
        'ja-JP': {
            label: '日本語',
            translation: jaJP,
            key4datefns: 'ja',
            key4mui: { dataGrid: 'jaJP', datePickers: 'jaJP' }
        },
        'ko-KR': {
            label: '한국어',
            translation: koKR,
            key4datefns: 'ko',
            key4mui: { dataGrid: 'koKR', datePickers: 'koKR' }
        },
        'my-MM': {
            label: 'မြန်မာဘာသာ',
            translation: myMM,
            // Neither date-fns nor MUI provides a Burmese locale.
            key4datefns: 'enUS',
            key4mui: { dataGrid: 'enUS', datePickers: 'enUS' }
        },
        'ne-NP': {
            label: 'नेपाली',
            translation: neNP,
            // Neither date-fns nor MUI provides a Nepali locale.
            key4datefns: 'enUS',
            key4mui: { dataGrid: 'enUS', datePickers: 'enUS' }
        },
        'nl-NL': {
            label: 'Nederlands',
            translation: nlNL,
            key4datefns: 'nl',
            key4mui: { dataGrid: 'nlNL', datePickers: 'nlNL' }
        },
        'pl-PL': {
            label: 'Polski',
            translation: plPL,
            key4datefns: 'pl',
            key4mui: { dataGrid: 'plPL', datePickers: 'plPL' }
        },
        'pt-BR': {
            label: 'Português (Brasil)',
            translation: ptBR,
            key4datefns: 'ptBR',
            key4mui: { dataGrid: 'ptBR', datePickers: 'ptBR' }
        },
        'pt-PT': {
            label: 'Português (Portugal)',
            translation: ptPT,
            key4datefns: 'pt',
            key4mui: { dataGrid: 'ptPT', datePickers: 'ptPT' }
        },
        'ru-RU': {
            label: 'Русский',
            translation: ruRU,
            key4datefns: 'ru',
            key4mui: { dataGrid: 'ruRU', datePickers: 'ruRU' }
        },
        'th-TH': {
            label: 'ไทย',
            translation: thTH,
            key4datefns: 'th',
            // MUI does not provide a Thai locale.
            key4mui: { dataGrid: 'enUS', datePickers: 'enUS' }
        },
        'tr-TR': {
            label: 'Türkçe',
            translation: trTR,
            key4datefns: 'tr',
            key4mui: { dataGrid: 'trTR', datePickers: 'trTR' }
        },
        'uk-UA': {
            label: 'Українська',
            translation: ukUA,
            key4datefns: 'uk',
            key4mui: { dataGrid: 'ukUA', datePickers: 'ukUA' }
        },
        'vi-VN': {
            label: 'Tiếng Việt',
            translation: viVN,
            key4datefns: 'vi',
            key4mui: { dataGrid: 'viVN', datePickers: 'viVN' }
        },
        'zh-TW': {
            label: '正體中文',
            translation: zhTW,
            key4datefns: 'zhTW',
            key4mui: { dataGrid: 'zhTW', datePickers: 'zhTW' }
        }
    }
};

const i18n = createI18n(
    localeConfig.fallback,
    localeConfig.key,
    localeConfig.resources satisfies Resource
);

export default i18n;
