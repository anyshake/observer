import {
    mdiAccessPointNetwork,
    mdiDesktopTower,
    mdiMapMarkerRadius,
    mdiRouterWireless,
    mdiSatelliteUplink,
    mdiTooltipText
} from '@mdi/js';
import Icon from '@mdi/react';
import { useTranslation } from 'react-i18next';

import { Banner } from '../../components/widget/Banner';
import { MapContainer } from '../../components/widget/MapContainer';
import { StatusCard } from '../../components/widget/StatusCard';
import { StatusList } from '../../components/widget/StatusList';
import { HomeConstraints } from '../../config/constraints';
import { useGetHomeDataQuery } from '../../graphql';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { useCredentialStore } from '../../stores/credential';
import { SystemUsage } from './SystemUsage';
import { useStationStatus } from './useStationStatus';
import { useUpgradeMessage } from './useUpgradeMessage';

const Home = () => {
    const { t } = useTranslation();
    const { getCredential } = useCredentialStore();
    const {
        data: getHomeDataData,
        error: getHomeDataError,
        loading: getHomeDataLoading
    } = useGetHomeDataQuery({ pollInterval: HomeConstraints.pollInterval });
    const {
        getStationCoordinates,
        getStationInfo,
        getDeviceInfo,
        getDeviceStatus,
        getLinkStatus,
        getServiceStatus
    } = useStationStatus(getHomeDataData);
    const { upgradeMessage } = useUpgradeMessage(getHomeDataData);

    return (
        <div className="app-page">
            <div className="app-stack-md">
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

            <div className="app-status-grid">
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

            <div className="app-station-grid">
                <div className="card app-panel app-location-card">
                    <h2 className="app-section-title">
                        <Icon className="app-icon" path={mdiMapMarkerRadius} size={1} />
                        <span>{t('views.Home.station_location.title')}</span>
                    </h2>
                    <MapContainer
                        scrollWheelZoom
                        zoomControl
                        dragging
                        height={300}
                        zoom={HomeConstraints.mapDefaultZoom}
                        maxZoom={HomeConstraints.mapMaxZoom}
                        minZoom={HomeConstraints.mapMinZoom}
                        tileUrl={HomeConstraints.mapTileUrl}
                        tileToken={getCredential().token}
                        layers={HomeConstraints.mapTileLayers}
                        coordinates={getStationCoordinates}
                    />
                </div>

                <div className="card app-panel app-services-card">
                    <h2 className="app-section-title">
                        <Icon className="app-icon" path={mdiTooltipText} size={1} />
                        <span>{t('views.Home.service_status.title')}</span>
                    </h2>
                    <StatusList data={getServiceStatus} />
                </div>
            </div>

            <SystemUsage data={getHomeDataData} />
        </div>
    );
};

export default Home;
