import {
    mdiAccessPointNetwork,
    mdiChartLine,
    mdiCrosshairsGps,
    mdiDesktopTower,
    mdiMapMarkerRadius,
    mdiRouterWireless,
    mdiSatelliteUplink,
    mdiTooltipText
} from '@mdi/js';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { LineChart } from '../../components/chart/LineChart';
import { Skeleton } from '../../components/ui/Skeleton';
import { Card } from '../../components/widget/Card';
import { MapContainer } from '../../components/widget/MapContainer';
import { type IPageHeaderNotice, PageHeader } from '../../components/widget/PageHeader';
import { StatusCard } from '../../components/widget/StatusCard';
import { HomeConstraints } from '../../config/constraints';
import { type GetHomeDataQuery, useGetHomeDataQuery } from '../../graphql';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { useCredentialStore } from '../../stores/credential';
import { HealthTile } from './HealthTile';
import { ServiceTable } from './ServiceTable';
import { SystemUsage } from './SystemUsage';
import { useHomeHealth } from './useHomeHealth';
import { useLinkHistory } from './useLinkHistory';
import { useStationStatus } from './useStationStatus';
import { useUpgradeMessage } from './useUpgradeMessage';

const Home = () => {
    const { t } = useTranslation();
    const { data, error, loading } = useGetHomeDataQuery({
        pollInterval: HomeConstraints.pollInterval
    });
    const { upgradeMessage } = useUpgradeMessage(data);

    if (loading && !data) {
        return (
            <div className="app-page app-home-page">
                <PageHeader
                    title="—"
                    status="warning"
                    statusText={t('views.Home.connectivity.loading')}
                />
                <Skeleton />
            </div>
        );
    }
    if (!data) {
        return (
            <div className="app-page app-home-page">
                <PageHeader
                    title="—"
                    status="error"
                    statusText={t('views.Home.station_header.offline')}
                    notices={[{ status: 'error', message: t('views.Home.connectivity.error') }]}
                />
            </div>
        );
    }

    return <HomeContent data={data} stale={!!error} upgradeMessage={upgradeMessage} />;
};

interface IHomeContent {
    readonly data: GetHomeDataQuery;
    readonly stale: boolean;
    readonly upgradeMessage: string;
}

const HomeContent = ({ data, stale, upgradeMessage }: IHomeContent) => {
    const { t } = useTranslation();
    const { getCredential } = useCredentialStore();
    const { throughput, recentErrors, currentRate } = useLinkHistory(data);
    const { station, coordinates, stationInfo, deviceInfo, deviceStatus, linkStatus, gnss } =
        useStationStatus(data);
    const health = useHomeHealth({ data, stale, currentRate, recentErrors });
    const { gnssEnabled } = data.getDeviceConfig;
    const notices = useMemo(() => {
        const list: IPageHeaderNotice[] = [];
        if (stale) {
            list.push({ status: 'error', message: t('views.Home.connectivity.error') });
        }
        if (upgradeMessage.length > 0) {
            list.push({ status: 'warning', message: upgradeMessage });
        }
        return list;
    }, [stale, t, upgradeMessage]);

    return (
        <div className="app-page app-home-page">
            <PageHeader
                title={station.name}
                meta={[station.code, station.place].filter(Boolean).join(' · ')}
                status={stale ? 'error' : 'ok'}
                statusText={t(stale ? 'views.Home.station_header.offline' : 'views.Home.station_header.online')}
                detail={`${t('views.Home.station_header.updated_at')} ${getTimeString(data.getCurrentTime)}`}
                notices={notices}
            />

            <div className="app-health-grid">
                <HealthTile title={t('views.Home.data_reception.title')} {...health.reception} />
                <HealthTile title={t('views.Home.time_source.title')} {...health.timing} />
                <HealthTile title={t('views.Home.reception_errors.title')} {...health.errors} />
                <HealthTile title={t('views.Home.service_summary.title')} {...health.service} />
            </div>

            <div className="app-status-grid">
                <StatusCard
                    iconPath={mdiDesktopTower}
                    title={t('views.Home.station_info.title')}
                    fields={stationInfo}
                />
                <StatusCard
                    iconPath={mdiRouterWireless}
                    title={t('views.Home.device_info.title')}
                    fields={deviceInfo}
                />
                <StatusCard
                    iconPath={mdiAccessPointNetwork}
                    title={t('views.Home.device_status.title')}
                    fields={deviceStatus}
                />
                <StatusCard
                    iconPath={mdiSatelliteUplink}
                    title={t('views.Home.link_status.title')}
                    fields={linkStatus}
                />
            </div>

            <div className="app-home-grid">
                <div className="app-home-col-wide">
                    <Card
                        title={t('views.Home.packet_throughput.title')}
                        iconPath={mdiChartLine}
                        className="app-home-card-body"
                    >
                        <p className="app-home-footnote">
                            {t('views.Home.packet_throughput.note', {
                                minutes: HomeConstraints.linkHistoryMinutes
                            })}
                        </p>
                        {throughput.length > 1 ? (
                            <LineChart height={240} lineWidth={2} yMin={0} data={throughput} />
                        ) : (
                            <div className="app-home-chart-empty">
                                <span className="loading loading-dots loading-sm" />
                                <span>{t('views.Home.packet_throughput.collecting')}</span>
                            </div>
                        )}
                    </Card>

                    <Card title={t('views.Home.service_status.title')} iconPath={mdiTooltipText}>
                        <ServiceTable services={data.getServiceStatus} stale={stale} />
                    </Card>
                </div>

                <div className="app-home-col-narrow">
                    {gnssEnabled && gnss.length > 0 ? (
                        <StatusCard
                            title={t('views.Home.gnss_status.title')}
                            iconPath={mdiCrosshairsGps}
                            fields={gnss}
                        >
                            <p className="app-home-footnote">{t('views.Home.gnss_status.note')}</p>
                        </StatusCard>
                    ) : (
                        <Card title={t('views.Home.gnss_status.title')} iconPath={mdiCrosshairsGps}>
                            <p className="app-help-text">
                                {t(
                                    gnssEnabled
                                        ? 'views.Home.gnss_status.pending'
                                        : 'views.Home.gnss_status.disabled'
                                )}
                            </p>
                        </Card>
                    )}

                    <div className="app-home-map-slot">
                        <Card
                            title={t('views.Home.station_location.title')}
                            iconPath={mdiMapMarkerRadius}
                            className="app-home-map"
                        >
                            <MapContainer
                                scrollWheelZoom
                                zoomControl
                                dragging
                                height={320}
                                zoom={HomeConstraints.mapDefaultZoom}
                                maxZoom={HomeConstraints.mapMaxZoom}
                                minZoom={HomeConstraints.mapMinZoom}
                                tileUrl={HomeConstraints.mapTileUrl}
                                tileToken={getCredential().token}
                                layers={HomeConstraints.mapTileLayers}
                                coordinates={coordinates}
                            />
                        </Card>
                    </div>
                </div>
            </div>

            <SystemUsage data={data} />
        </div>
    );
};

export default Home;
