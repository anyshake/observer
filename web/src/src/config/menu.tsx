import { mdiCog, mdiDatabaseExport, mdiFileClock, mdiServerNetwork, mdiWaveform } from '@mdi/js';

import { Translation } from './locale';

export interface IMenuItem {
    readonly url: string;
    readonly icon: string;
    readonly home?: boolean;
    readonly label: Translation;
}

export const menuConfig: IMenuItem[] = [
    {
        home: true,
        url: '/',
        label: {
            'de-DE': 'Stationsstatus',
            'en-US': 'Station Status',
            'fr-FR': 'Etat de la station',
            'ja-JP': '現在の状況',
            'zh-TW': '當前站況',
            'pt-PT': 'Status da Estação',
            'ru-RU': 'Статус станции',
            'id-ID': 'Status Stasiun',
            'tr-TR': 'İstasyon Durumu'
        },
        icon: mdiServerNetwork
    },
    {
        url: '/realtime',
        label: {
            'de-DE': 'Echtzeit-Wellenform',
            'en-US': 'Real-time Waveform',
            'fr-FR': 'Ondes en direct',
            'ja-JP': 'リアルタイム波形',
            'zh-TW': '即時波形',
            'pt-PT': 'Forma de Onda em Tempo Real',
            'ru-RU': 'Волновая форма в реальном времени',
            'id-ID': 'Gelombang Waktu Nyata',
            'tr-TR': 'Gerçek Zamanlı Dalga Formu'
        },
        icon: mdiWaveform
    },
    {
        url: '/history',
        label: {
            'de-DE': 'Historische Wellenform',
            'en-US': 'History Waveform',
            'fr-FR': 'Historique des ondes',
            'ja-JP': '履歴波形',
            'zh-TW': '歷史資料',
            'pt-PT': 'Forma de Onda Histórica',
            'ru-RU': 'Историческая волновая форма',
            'id-ID': 'Gelombang Sejarah',
            'tr-TR': 'Geçmiş Dalga Formu'
        },
        icon: mdiFileClock
    },
    {
        url: '/download',
        label: {
            'de-DE': 'Daten herunterladen',
            'en-US': 'Data Download',
            'fr-FR': 'Téléchargement des données',
            'ja-JP': 'データダウンロード',
            'zh-TW': '資料下載',
            'pt-PT': 'Download de Dados',
            'ru-RU': 'Загрузка данных',
            'id-ID': 'Unduh Data',
            'tr-TR': 'Veri İndir'
        },
        icon: mdiDatabaseExport
    },
    {
        url: '/settings',
        label: {
            'de-DE': 'Systemeinstellungen',
            'en-US': 'System Settings',
            'fr-FR': 'Paramètres du système',
            'ja-JP': 'システム設定',
            'zh-TW': '系統設定',
            'pt-PT': 'Configurações do Sistema',
            'ru-RU': 'Настройки системы',
            'id-ID': 'Pengaturan Sistem',
            'tr-TR': 'Sistem Ayarları'
        },
        icon: mdiCog
    }
];
