import {
    mdiAccessPointNetwork,
    mdiChip,
    mdiDesktopTower,
    mdiHarddisk,
    mdiMapMarkerRadius,
    mdiMemory,
    mdiRouterWireless,
    mdiSatelliteUplink,
    mdiTooltipText
} from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Banner } from '../../components/Banner';
import { GaugeChart } from '../../components/GaugeChart';
import { LineChart } from '../../components/LineChart';
import { MapContainer } from '../../components/MapContainer';
import { Skeleton } from '../../components/Skeleton';
import { StatusCard } from '../../components/StatusCard';
import { StatusList } from '../../components/StatusList';
import { useGetHomeDataQuery, useGetUpgradeStatusLazyQuery } from '../../graphql';
import { getTimeString } from '../../helpers/utils/getTimeString';

const lineChartRetention = 100;
const pollInterval = 2000;
const maxGapSeconds = 5 * (pollInterval / 1000);

const Home = () => {
    const { t } = useTranslation();
    const lineChartConfig = useMemo(
        () => ({
            height: 250,
            lineColor: '#8A3EED',
            lineWidth: 5,
            yMax: 100,
            yMin: 0,
            yInterval: 20,
            yPosition: 'left' as const
        }),
        []
    );

    const {
        data: getHomeDataData,
        error: getHomeDataError,
        loading: getHomeDataLoading
    } = useGetHomeDataQuery({ pollInterval });
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

    const [memoryData, setMemoryData] = useState<Array<[number, number | null]>>([]);
    const [cpuData, setCpuData] = useState<Array<[number, number | null]>>([]);
    const memoryDataRef = useRef(memoryData);
    const cpuDataRef = useRef(cpuData);
    const processChartData = useCallback(
        (
            prevData: Array<[number, number | null]>,
            newPoint: [number, number],
            retention: number,
            maxGap: number
        ): Array<[number, number | null]> => {
            if (prevData.length === 0) {
                return [newPoint];
            }

            const lastValidPoint = [...prevData].reverse().find((point) => point[1] !== null);
            if (!lastValidPoint) {
                return [...prevData.slice(-retention), newPoint];
            }

            const timeDiff = (newPoint[0] - lastValidPoint[0]) / 1000;
            if (timeDiff > maxGap) {
                return [...prevData.slice(-retention), [lastValidPoint[0] + 1, null], newPoint];
            }

            return [...prevData.slice(-retention), newPoint];
        },
        []
    );
    useEffect(() => {
        if (getHomeDataData) {
            const newCpuData = processChartData(
                cpuDataRef.current,
                [getHomeDataData.getCurrentTime, getHomeDataData.getSystemStatus.cpu ?? 0],
                lineChartRetention,
                maxGapSeconds
            );
            cpuDataRef.current = newCpuData;
            setCpuData(newCpuData);

            const newMemoryData = processChartData(
                memoryDataRef.current,
                [getHomeDataData.getCurrentTime, getHomeDataData.getSystemStatus.memory ?? 0],
                lineChartRetention,
                maxGapSeconds
            );
            memoryDataRef.current = newMemoryData;
            setMemoryData(newMemoryData);
        }
    }, [getHomeDataData, processChartData]);

    const [getUpgradeStatus] = useGetUpgradeStatusLazyQuery();
    const [upgradeMessage, setUpgradeMessage] = useState('');

    const generateUpgradeMessage = useCallback(
        (
            current: string,
            latest: string,
            required: string,
            eligible: boolean,
            applied: boolean
        ) => {
            if (eligible) {
                if (applied) {
                    return t('views.Home.upgrade.restart_needed', { latest, current });
                }
                return t('views.Home.upgrade.update_available', { latest });
            }
            if (current !== latest) {
                return t('views.Home.upgrade.manual_upgrade_needed', { latest, current, required });
            }

            return '';
        },
        [t]
    );

    useEffect(() => {
        if (getHomeDataData?.getCurrentUser.admin) {
            getUpgradeStatus().then(({ data }) => {
                let msg = '';
                if (data?.getUpgradeStatus) {
                    const { current, latest, required, eligible, applied } = data.getUpgradeStatus;
                    msg = generateUpgradeMessage(current, latest, required, eligible, applied);
                }
                setUpgradeMessage(msg);
            });
        }
    }, [generateUpgradeMessage, getHomeDataData?.getCurrentUser.admin, getUpgradeStatus, t]);

    return (
        <div className="container mx-auto space-y-6 p-4">
            <div className="space-y-3">
                {upgradeMessage.length > 0 && <Banner status="warning" message={upgradeMessage} />}
                <Banner
                    status={getHomeDataLoading ? 'warning' : getHomeDataError ? 'error' : 'ok'}
                    message={
                        getHomeDataLoading
                            ? t('views.Home.connectivity.loading')
                            : getHomeDataError
                              ? t('views.Home.connectivity.error')
                              : t('views.Home.connectivity.success', {
                                    uptime: (
                                        getHomeDataData!.getSystemStatus.uptime / 1000
                                    ).toFixed(1),
                                    updatedAt: getHomeDataData
                                        ? getTimeString(getHomeDataData.getCurrentTime)
                                        : 'N/A'
                                })
                    }
                />
            </div>

            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4">
                <StatusCard
                    iconPath={mdiDesktopTower}
                    title={t('views.Home.station_info.title')}
                    fields={getStationInfo}
                />
                <StatusCard
                    iconPath={mdiRouterWireless}
                    title={t('views.Home.device_info.title')}
                    fields={getDeviceInfo}
                />
                <StatusCard
                    iconPath={mdiAccessPointNetwork}
                    title={t('views.Home.device_status.title')}
                    fields={getDeviceStatus}
                />
                <StatusCard
                    iconPath={mdiSatelliteUplink}
                    title={t('views.Home.link_status.title')}
                    fields={getLinkStatus}
                />
            </div>

            <div className="grid grid-cols-1 gap-6 md:grid-cols-16">
                <div className="card bg-base-100 flex flex-col space-y-4 p-5 text-gray-700 shadow-md md:col-span-7">
                    <h2 className="flex items-center space-x-2 text-lg font-bold">
                        <Icon className="flex-shrink-0" path={mdiMapMarkerRadius} size={1} />
                        <span>{t('views.Home.station_location.title')}</span>
                    </h2>
                    <MapContainer
                        scrollWheelZoom
                        zoomControl
                        dragging
                        zoom={6}
                        maxZoom={7}
                        minZoom={3}
                        height={300}
                        tile="/tiles/{z}/{x}/{y}.webp"
                        coordinates={getStationCoordinates}
                    />
                </div>

                <div className="card bg-base-100 flex flex-col space-y-4 p-5 text-gray-700 shadow-md md:col-span-9">
                    <h2 className="flex items-center space-x-2 text-lg font-bold">
                        <Icon className="flex-shrink-0" path={mdiTooltipText} size={1} />
                        <span>{t('views.Home.service_status.title')}</span>
                    </h2>
                    <StatusList data={getServiceStatus} />
                </div>
            </div>

            <div className="grid grid-cols-1 gap-6 md:grid-cols-19">
                <div className="card bg-base-100 flex flex-col space-y-8 p-5 text-gray-700 shadow-md md:col-span-7">
                    <h2 className="flex items-center space-x-2 text-lg font-bold">
                        <Icon className="flex-shrink-0" path={mdiChip} size={1} />
                        <span>{t('views.Home.cpu_usage.title')}</span>
                    </h2>
                    {cpuData.length > 1 ? (
                        <LineChart {...lineChartConfig} data={cpuData} />
                    ) : (
                        <Skeleton height={lineChartConfig.height} />
                    )}
                </div>
                <div className="card bg-base-100 flex flex-col space-y-8 p-5 text-gray-700 shadow-md md:col-span-7">
                    <h2 className="flex items-center space-x-2 text-lg font-bold">
                        <Icon className="flex-shrink-0" path={mdiMemory} size={1} />
                        <span>{t('views.Home.memory_usage.title')}</span>
                    </h2>
                    {memoryData.length > 1 ? (
                        <LineChart {...lineChartConfig} data={memoryData} />
                    ) : (
                        <Skeleton height={lineChartConfig.height} />
                    )}
                </div>
                <div className="card bg-base-100 flex flex-col space-y-8 p-5 text-gray-700 shadow-md md:col-span-5">
                    <h2 className="flex items-center space-x-2 text-lg font-bold">
                        <Icon className="flex-shrink-0" path={mdiHarddisk} size={1} />
                        <span>{t('views.Home.disk_usage.title')}</span>
                    </h2>
                    <GaugeChart
                        height={lineChartConfig.height}
                        valueMax={100}
                        animation={true}
                        value={getHomeDataData?.getSystemStatus.disk ?? 0}
                    />
                </div>
            </div>
        </div>
    );
};

export default Home;
