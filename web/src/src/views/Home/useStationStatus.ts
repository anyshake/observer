import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import type { GetHomeDataQuery } from '../../graphql';
import { getTimeString } from '../../helpers/utils/getTimeString';

export const useStationStatus = (getHomeDataData: GetHomeDataQuery | undefined) => {
    const { t } = useTranslation();
    const getStationCoordinates = useMemo(
        () =>
            getHomeDataData
                ? [
                      getHomeDataData.getDeviceInfo.latitude ?? 0,
                      getHomeDataData.getDeviceInfo.longitude ?? 0,
                      getHomeDataData.getDeviceInfo.elevation ?? 0
                  ]
                : [0, 0, 0],
        [getHomeDataData]
    );
    const getStationInfo = useMemo(
        () => [
            {
                label: t('views.Home.station_info.name'),
                value: getHomeDataData?.getStationConfig.station_name ?? 'N/A'
            },
            {
                label: t('views.Home.station_info.country'),
                value: getHomeDataData?.getStationConfig.station_country ?? 'N/A'
            },
            {
                label: t('views.Home.station_info.place'),
                value: getHomeDataData?.getStationConfig.station_place ?? 'N/A'
            },
            {
                label: t('views.Home.station_info.affiliation'),
                value: getHomeDataData?.getStationConfig.station_affiliation ?? 'N/A'
            },
            {
                label: t('views.Home.station_info.channels'),
                value:
                    (getHomeDataData?.getDeviceConfig.channelCodes ?? []).join(', ').length > 0
                        ? getHomeDataData?.getDeviceConfig.channelCodes.join(', ')
                        : 'N/A'
            }
        ],
        [t, getHomeDataData]
    );
    const getDeviceInfo = useMemo(
        () => [
            {
                label: t('views.Home.device_info.product_model'),
                value: getHomeDataData?.getDeviceConfig.model ?? 'N/A'
            },
            {
                label: t('views.Home.device_info.device_id'),
                value: getHomeDataData?.getDeviceId ?? 'N/A'
            },
            {
                label: t('views.Home.device_info.link_protocol'),
                value: getHomeDataData?.getDeviceConfig.protocol ?? 'N/A'
            },
            {
                label: t('views.Home.device_info.packet_interval'),
                value: `${getHomeDataData?.getDeviceConfig.packetInterval ?? 0} ms`
            },
            {
                label: t('views.Home.device_info.gnss_active'),
                value: getHomeDataData?.getDeviceConfig.gnssEnabled
                    ? t('views.Home.device_info.yes')
                    : t('views.Home.device_info.no')
            }
        ],
        [t, getHomeDataData]
    );
    const getDeviceStatus = useMemo(
        () => [
            {
                label: t('views.Home.device_status.sample_rate'),
                value: `${getHomeDataData?.getDeviceConfig.sampleRate ?? 0} SPS`
            },
            {
                label: t('views.Home.device_status.latitude'),
                value: getHomeDataData?.getDeviceInfo.latitude
                    ? `${getHomeDataData.getDeviceInfo.latitude.toFixed(2)}°`
                    : 'N/A'
            },
            {
                label: t('views.Home.device_status.longitude'),
                value: getHomeDataData?.getDeviceInfo.longitude
                    ? `${getHomeDataData.getDeviceInfo.longitude?.toFixed(2)}°`
                    : 'N/A'
            },
            {
                label: t('views.Home.device_status.elevation'),
                value: getHomeDataData?.getDeviceInfo.elevation
                    ? `${getHomeDataData.getDeviceInfo.elevation?.toFixed(2)} m`
                    : 'N/A'
            },
            {
                label: t('views.Home.device_status.temperature'),
                value: getHomeDataData?.getDeviceInfo.temperature
                    ? `${getHomeDataData.getDeviceInfo.temperature.toFixed(2)} °C`
                    : 'N/A'
            }
        ],
        [t, getHomeDataData]
    );
    const getLinkStatus = useMemo(
        () => [
            {
                label: t('views.Home.link_status.started_at'),
                value: getHomeDataData?.getDeviceStatus.startedAt
                    ? getTimeString(getHomeDataData.getDeviceStatus.startedAt)
                    : 'N/A'
            },
            {
                label: t('views.Home.link_status.last_frame'),
                value: getHomeDataData?.getDeviceStatus.updatedAt
                    ? getTimeString(getHomeDataData.getDeviceStatus.updatedAt)
                    : 'N/A'
            },
            {
                label: t('views.Home.link_status.frames'),
                value: `${getHomeDataData?.getDeviceStatus.frames ?? 0}`
            },
            {
                label: t('views.Home.link_status.errors'),
                value: `${getHomeDataData?.getDeviceStatus.errors ?? 0}`
            },
            {
                label: t('views.Home.link_status.messages'),
                value: `${getHomeDataData?.getDeviceStatus.messages ?? 0}`
            }
        ],
        [t, getHomeDataData]
    );
    const getServiceStatus = useMemo(
        () =>
            getHomeDataData?.getServiceStatus
                ? [...getHomeDataData.getServiceStatus]
                      .sort((a, b) => (a.serviceId ?? '').localeCompare(b.serviceId ?? ''))
                      .map((serviceStatus) => ({
                          title: serviceStatus.name!,
                          description: serviceStatus.description!,
                          running: serviceStatus.isRunning ?? false
                      }))
                : [],
        [getHomeDataData?.getServiceStatus]
    );

    return {
        getStationCoordinates,
        getStationInfo,
        getDeviceInfo,
        getDeviceStatus,
        getLinkStatus,
        getServiceStatus
    };
};
