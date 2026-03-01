import { mdiCog, mdiDatabaseExport, mdiFileClock, mdiServerNetwork, mdiWaveform } from '@mdi/js';

import { routerConfig } from './router';

export interface IMenuItem {
    readonly url: string;
    readonly icon: string;
    readonly home?: boolean;
    readonly label: string;
}

export const menuConfig: IMenuItem[] = [
    {
        home: true,
        url: routerConfig.routes.home.uri,
        label: routerConfig.routes.home.title,
        icon: mdiServerNetwork
    },
    {
        url: routerConfig.routes.realtime.uri,
        label: routerConfig.routes.realtime.title,
        icon: mdiWaveform
    },
    {
        url: routerConfig.routes.history.uri,
        label: routerConfig.routes.history.title,
        icon: mdiFileClock
    },
    {
        url: routerConfig.routes.download.uri,
        label: routerConfig.routes.download.title,
        icon: mdiDatabaseExport
    },
    {
        url: routerConfig.routes.settings.uri,
        label: routerConfig.routes.settings.title,
        icon: mdiCog
    }
];
