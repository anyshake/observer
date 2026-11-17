import { useTranslation } from 'react-i18next';

import { HomeConstraints } from '../../config/constraints';
import type { GetHomeDataQuery } from '../../graphql';
import { getTimeString } from '../../helpers/utils/getTimeString';
import type { IHealth } from './healthStatus';

interface IHomeHealth {
    readonly data: GetHomeDataQuery;
    readonly stale: boolean;
    readonly currentRate: number | null;
    readonly recentErrors: number;
}

// Summarizes reception, timing, error and service health for the overview tiles.
export const useHomeHealth = ({ data, stale, currentRate, recentErrors }: IHomeHealth) => {
    const { t } = useTranslation();
    const unknown = t('views.Home.health_level.unknown');
    const { packetInterval, gnssEnabled } = data.getDeviceConfig;
    const status = data.getDeviceStatus;
    const gnss = status.gnssStatus;
    const services = data.getServiceStatus ?? [];
    const running = services.filter(({ isRunning }) => isRunning).length;
    const fresh =
        status.frames > 0 &&
        Math.abs(data.getCurrentTime - status.updatedAt) <=
            Math.max(HomeConstraints.staleAfter, packetInterval * 5);
    const lastFrame = t('views.Home.data_reception.last_frame', {
        time: status.updatedAt > 0 ? getTimeString(status.updatedAt) : '—'
    });

    const reception: IHealth = stale
        ? { level: 'unknown', value: unknown, detail: t('views.Home.data_reception.unavailable') }
        : status.frames === 0
          ? { level: 'warning', value: t('views.Home.data_reception.waiting'), detail: lastFrame }
          : fresh
            ? {
                  level: 'ok',
                  value:
                      currentRate !== null
                          ? t('views.Home.data_reception.packet_rate', { value: currentRate.toFixed(1) })
                          : t('views.Home.data_reception.receiving'),
                  detail: lastFrame
              }
            : { level: 'error', value: t('views.Home.data_reception.stale'), detail: lastFrame };

    const timing: IHealth =
        stale || !fresh
            ? { level: 'unknown', value: unknown, detail: t('views.Home.time_source.unavailable') }
            : !gnssEnabled
              ? { level: 'ok', value: 'NTP', detail: t('views.Home.time_source.ntp_note') }
              : !gnss
                ? { level: 'unknown', value: unknown, detail: t('views.Home.gnss_status.pending') }
                : gnss.locked
                  ? {
                        level: 'ok',
                        value: t('views.Home.time_source.gnss_locked'),
                        detail: t('views.Home.time_source.gnss_locked_note')
                    }
                  : {
                        level: 'warning',
                        value: t('views.Home.time_source.holdover'),
                        detail: t('views.Home.time_source.holdover_note')
                    };

    const errors: IHealth = {
        level: stale ? 'unknown' : recentErrors > 0 ? 'warning' : 'ok',
        value: stale ? unknown : `${recentErrors}`,
        detail: t('views.Home.reception_errors.note', {
            minutes: HomeConstraints.linkHistoryMinutes,
            total: status.errors
        })
    };

    const service: IHealth = {
        // Stopped services may be disabled on purpose, so they never raise a warning here.
        level: stale || services.length === 0 ? 'unknown' : 'ok',
        value: stale
            ? unknown
            : t('views.Home.service_summary.running_count', { running, total: services.length }),
        detail: t('views.Home.service_summary.note')
    };

    return { reception, timing, errors, service };
};
