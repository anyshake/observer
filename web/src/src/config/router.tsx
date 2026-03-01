import { JSX, lazy, LazyExoticComponent } from 'react';

import { RouterMode } from '../components/ui/RouterWrapper';

export type RouterProp<T> = Record<string, T>;

export interface IRouterComponent {
    currentLocale: string;
    onSwitchLocale: (newLocale: string) => void;
    locales: Record<string, string>;
}

export interface IRoute {
    readonly uri: string;
    readonly title: string;
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
            title: 'config.router.home'
        },
        realtime: {
            uri: '/realtime',
            element: RealTime,
            title: 'config.router.realtime'
        },
        history: {
            uri: '/history',
            element: History,
            title: 'config.router.history'
        },
        download: {
            uri: '/download',
            element: Download,
            title: 'config.router.download'
        },
        settings: {
            uri: '/settings',
            element: Settings,
            title: 'config.router.settings'
        },
        default: {
            uri: '*',
            element: NotFound,
            title: 'config.router.default'
        }
    }
};
