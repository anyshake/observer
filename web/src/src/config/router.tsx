import { JSX, lazy, LazyExoticComponent } from 'react';

import { RouterMode } from '../components/RouterWrapper';
import { Translation } from './locale';

export type RouterProp<T> = Record<string, T>;

export interface IRouterComponent {
    currentLocale: string;
}

export interface IRoute {
    readonly uri: string;
    readonly title: Translation;
    readonly element: LazyExoticComponent<(props: IRouterComponent) => JSX.Element>;
}

interface IRouterConfig {
    readonly mode: RouterMode;
    readonly basename: string;
    readonly routes: Record<string, IRoute>;
}

const Home = lazy(() => import('../views/Home'));
const RealTime = lazy(() => import('../views/RealTime'));
const History = lazy(() => import('../views/History'));
const Download = lazy(() => import('../views/Download'));
const Settings = lazy(() => import('../views/Settings'));
const NotFound = lazy(() => import('../views/NotFound'));

export const routerConfig: IRouterConfig = {
    basename: '/',
    mode: 'hash',
    routes: {
        home: {
            uri: '/',
            element: Home,
            title: {
                'de-DE': 'Stationsstatus',
                'en-US': 'Station Status',
                'fr-FR': 'Etat de la station',
                'ja-JP': '現在の状況',
                'zh-TW': '當前站況',
                'pt-PT': 'Status da Estação',
                'ru-RU': 'Статус станции',
                'id-ID': 'Status Stasiun',
                'tr-TR': 'İstasyon Durumu'
            }
        },
        realtime: {
            uri: '/realtime',
            element: RealTime,
            title: {
                'de-DE': 'Echtzeit-Wellenform',
                'en-US': 'Real-time Waveform',
                'fr-FR': 'Ondes en direct',
                'ja-JP': 'リアルタイム波形',
                'zh-TW': '即時波形',
                'pt-PT': 'Forma de Onda em Tempo Real',
                'ru-RU': 'Волновая форма в реальном времени',
                'id-ID': 'Gelombang Waktu Nyata',
                'tr-TR': 'Gerçek Zamanlı Dalga Formu'
            }
        },
        history: {
            uri: '/history',
            element: History,
            title: {
                'de-DE': 'Historische Wellenform',
                'en-US': 'History Waveform',
                'fr-FR': 'Historique des ondes',
                'ja-JP': '履歴波形',
                'zh-TW': '歷史資料',
                'pt-PT': 'Forma de Onda Histórica',
                'ru-RU': 'Историческая волновая форма',
                'id-ID': 'Gelombang Sejarah',
                'tr-TR': 'Geçmiş Dalga Formu'
            }
        },
        export: {
            uri: '/download',
            element: Download,
            title: {
                'de-DE': 'Daten herunterladen',
                'en-US': 'Data Download',
                'fr-FR': 'Téléchargement des données',
                'ja-JP': 'データダウンロード',
                'zh-TW': '資料下載',
                'pt-PT': 'Download de Dados',
                'ru-RU': 'Загрузка данных',
                'id-ID': 'Unduh Data',
                'tr-TR': 'Veri İndir'
            }
        },
        setting: {
            uri: '/settings',
            element: Settings,
            title: {
                'de-DE': 'Systemeinstellungen',
                'en-US': 'System Settings',
                'fr-FR': 'Paramètres du système',
                'ja-JP': 'システム設定',
                'zh-TW': '系統設定',
                'pt-PT': 'Configurações do Sistema',
                'ru-RU': 'Настройки системы',
                'id-ID': 'Pengaturan Sistem',
                'tr-TR': 'Sistem Ayarları'
            }
        },
        default: {
            uri: '*',
            element: NotFound,
            title: {
                'en-US': 'Page Not Found',
                'fr-FR': 'Page introuvable',
                'ja-JP': 'ページが見つかりません',
                'zh-TW': '找不到頁面',
                'de-DE': 'Seite nicht gefunden',
                'pt-PT': 'Página não encontrada',
                'ru-RU': 'Страница не найдена',
                'id-ID': 'Halaman Tidak Ditemukan',
                'tr-TR': 'Sayfa Bulunamadı'
            }
        }
    }
};
