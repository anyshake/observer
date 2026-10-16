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

                <div className="card bg-base-100 flex flex-col space-y-4 p-5 text-gray-700 shadow-md md:col-span-9">
                    <h2 className="flex items-center space-x-2 text-lg font-bold">
                        <Icon className="flex-shrink-0" path={mdiTooltipText} size={1} />
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
