import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import type { GetHomeDataQuery } from '../../graphql';
import { getStationIdentity } from '../../helpers/app/getStationIdentity';
import { getTimeString } from '../../helpers/utils/getTimeString';

interface IField {
    readonly label: string;
    readonly value: string;
}

const text = (value: unknown) => (typeof value === 'string' && value.length > 0 ? value : '');
const fixed = (value: number | null | undefined, digits: number, unit: string) =>
    typeof value === 'number' ? `${value.toFixed(digits)}${unit}` : 'N/A';
const timestamp = (value: number) => (value > 0 ? getTimeString(value) : 'N/A');

export const useStationStatus = (data: GetHomeDataQuery) => {
    const { t } = useTranslation();

    const station = useMemo(() => {
        const identity = getStationIdentity(data.getStationConfig);
        return { ...identity, name: identity.name || 'N/A' };
    }, [data.getStationConfig]);

    const coordinates = useMemo(() => {
        const { latitude, longitude, elevation } = data.getDeviceInfo;
        return [latitude ?? 0, longitude ?? 0, elevation ?? 0];
    }, [data.getDeviceInfo]);

    const stationInfo = useMemo<IField[]>(() => {
        const config = data.getStationConfig ?? {};
        return [
            { label: t('views.Home.station_info.name'), value: text(config.station_name) || 'N/A' },
            {
                label: t('views.Home.station_info.country'),
                value: text(config.station_country) || 'N/A'
            },
            {
                label: t('views.Home.station_info.place'),
                value: text(config.station_place) || 'N/A'
            },
            {
                label: t('views.Home.station_info.affiliation'),
                value: text(config.station_affiliation) || 'N/A'
            },
            {
                label: t('views.Home.station_info.channels'),
                value: data.getDeviceConfig.channelCodes.join(', ') || 'N/A'
            }
        ];
    }, [t, data.getStationConfig, data.getDeviceConfig.channelCodes]);

    const deviceInfo = useMemo<IField[]>(() => {
        const config = data.getDeviceConfig;
        return [
            { label: t('views.Home.device_info.product_model'), value: config.model || 'N/A' },
            { label: t('views.Home.device_info.device_id'), value: data.getDeviceId || 'N/A' },
            { label: t('views.Home.device_info.link_protocol'), value: config.protocol || 'N/A' },
            {
                label: t('views.Home.device_info.packet_interval'),
                value: `${config.packetInterval} ms`
            },
            {
                label: t('views.Home.device_info.gnss_active'),
                value: t(
                    config.gnssEnabled ? 'views.Home.device_info.yes' : 'views.Home.device_info.no'
                )
            }
        ];
    }, [t, data.getDeviceConfig, data.getDeviceId]);

    const deviceStatus = useMemo<IField[]>(() => {
        const info = data.getDeviceInfo;
        return [
            {
                label: t('views.Home.device_status.sample_rate'),
                value: `${data.getDeviceConfig.sampleRate} SPS`
            },
            { label: t('views.Home.device_status.latitude'), value: fixed(info.latitude, 2, '°') },
            {
                label: t('views.Home.device_status.longitude'),
                value: fixed(info.longitude, 2, '°')
            },
            {
                label: t('views.Home.device_status.elevation'),
                value: fixed(info.elevation, 2, ' m')
            },
            {
                label: t('views.Home.device_status.temperature'),
                value: fixed(info.temperature, 2, ' °C')
            }
        ];
    }, [t, data.getDeviceInfo, data.getDeviceConfig.sampleRate]);

    const linkStatus = useMemo<IField[]>(() => {
        const status = data.getDeviceStatus;
        return [
            { label: t('views.Home.link_status.started_at'), value: timestamp(status.startedAt) },
            { label: t('views.Home.link_status.last_frame'), value: timestamp(status.updatedAt) },
            { label: t('views.Home.link_status.frames'), value: `${status.frames}` },
            { label: t('views.Home.link_status.errors'), value: `${status.errors}` },
            { label: t('views.Home.link_status.messages'), value: `${status.messages}` }
        ];
    }, [t, data.getDeviceStatus]);

    const gnss = useMemo<IField[]>(() => {
        const status = data.getDeviceStatus.gnssStatus;
        if (!status) {
            return [];
        }
        const valid = (ok: boolean) => t(ok ? 'views.Home.gnss_status.valid' : 'views.Home.gnss_status.invalid');
        return [
            { label: t('views.Home.gnss_status.position_fix'), value: valid(status.fixValid) },
            { label: t('views.Home.gnss_status.nmea_time'), value: valid(status.timeValid) },
            {
                label: t('views.Home.gnss_status.satellites'),
                value: status.satellites >= 31 ? '31+' : `${status.satellites}`
            },
            { label: 'HDOP', value: fixed(status.hdop, 1, '') },
            { label: t('views.Home.gnss_status.pps_age'), value: fixed(status.ppsAgeSeconds, 0, ' s') }
        ];
    }, [t, data.getDeviceStatus.gnssStatus]);

    return { station, coordinates, stationInfo, deviceInfo, deviceStatus, linkStatus, gnss };
};
