import logo from '/anyshake.svg';

import { getVersionTag } from '../helpers/app/getVersionTag';
import { Translation } from './locale';

interface IGlobalConfig {
    readonly repository: string;
    readonly logo: string;
    readonly version: string;
    readonly name: Translation;
    readonly copyright: string;
    readonly homepage: string;
    readonly update: number;
    readonly footer: Translation;
}

export const globalConfig: IGlobalConfig = {
    version: getVersionTag(),
    logo,
    name: {
        'de-DE': 'AnyShake Observer',
        'en-US': 'AnyShake Observer',
        'fr-FR': 'AnyShake Observer',
        'ja-JP': 'AnyShake Observer',
        'zh-TW': 'AnyShake Observer',
        'pt-PT': 'AnyShake Observer',
        'ru-RU': 'AnyShake Observer',
        'id-ID': 'AnyShake Observer',
        'tr-TR': 'AnyShake Observer',
        'pl-PL': 'AnyShake Observer'
    },
    footer: {
        'de-DE': '"Hoher Geist der Erde."',
        'en-US': '"Listen to the whispering earth."',
        'fr-FR': '"Ecouter le murmure de la terre."',
        'ja-JP': '「地球のささやき」',
        'zh-TW': '「聽見地球」',
        'pt-PT': '"Ouça a terra sussurrar."',
        'ru-RU': '"Прислушайтесь к шепоту Земли."',
        'id-ID': '"Dengarkan bisikan bumi."',
        'tr-TR': '"Dünyanın fısıltısını dinleyin."',
        'pl-PL': '"Wsłuchaj się w szept Ziemi."'
    },
    update: 10 * 60 * 1000,
    copyright: 'SensePlex Limited',
    homepage: 'https://anyshake.org',
    repository: 'https://github.com/anyshake/observer'
};
