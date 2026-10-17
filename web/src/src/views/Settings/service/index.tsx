import {
    mdiApplicationCogOutline,
    mdiBackupRestore,
    mdiPlay,
    mdiRestart,
    mdiSquareEditOutline,
    mdiStop
} from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import { ErrorPage } from '../../../components/ui/ErrorPage';
import { StatusCard } from '../../../components/widget/StatusCard';
import { getTimeString } from '../../../helpers/utils/getTimeString';
import { ServiceConfigModal } from './ServiceConfigModal';
import { useServiceActions } from './useServiceActions';
import { useServiceConfig } from './useServiceConfig';
import { useServiceData } from './useServiceData';

export const Service = () => {
    const { t } = useTranslation();

    const {
        getServiceDataRefetch,
        getServiceDataError,
        getServiceDataLoading,
        serviceConfig,
        serviceStatus
    } = useServiceData();
    const { handleRestartService, handleStartStopService, handleResetServiceConfig } =
        useServiceActions(getServiceDataRefetch);
    const configure = useServiceConfig(getServiceDataRefetch);
    const { setSelectedServiceId, setIsConfigureModalOpen } = configure;
    const getStatusCardFields = useCallback(
        (status: (typeof serviceStatus)[string]) => [
            {
                label: t('views.Settings.Service.service_status.is_running'),
                value: status?.isRunning
                    ? t('views.Settings.Service.service_status.yes')
                    : t('views.Settings.Service.service_status.no')
            },
            {
                label: t('views.Settings.Service.service_status.restarts'),
                value: status?.restarts
            },
            {
                label: t('views.Settings.Service.service_status.stopped_at'),
                value: status?.stoppedAt ? getTimeString(status?.stoppedAt) : 'N/A'
            },
            {
                label: t('views.Settings.Service.service_status.started_at'),
                value: status?.startedAt ? getTimeString(status?.startedAt) : 'N/A'
            },
            {
                label: t('views.Settings.Service.service_status.updated_at'),
                value: getTimeString(status?.updatedAt)
            }
        ],
        [t]
    );

    return getServiceDataError ? (
        <ErrorPage
            content={getServiceDataError.message}
            debug={JSON.stringify(getServiceDataError)}
        />
    ) : getServiceDataLoading ? (
        <div className="app-center">
            <span className="loading loading-spinner app-loading-delayed" />
        </div>
    ) : (
        <div className="app-service-grid">
            {Object.keys(serviceConfig)
                .sort((a, b) =>
                    serviceConfig[a].serviceName.localeCompare(serviceConfig[b].serviceName)
                )
                .map((serviceId, index) => {
                    const config = serviceConfig[serviceId];
                    const status = serviceStatus[serviceId];
                    return (
                        <StatusCard
                            key={`${serviceId}-${index}`}
                            title={status.serviceName}
                            iconPath={mdiApplicationCogOutline}
                            fields={getStatusCardFields(status)}
                        >
                            <div className="join app-service-actions">
                                {status?.isRunning && (
                                    <button
                                        className="join-item btn btn-sm app-button-toolbar"
                                        onClick={() => handleRestartService(serviceId)}
                                    >
                                        <Icon className="app-icon" path={mdiRestart} size={0.7} />
                                        <span>
                                            {t(
                                                'views.Settings.Service.restart_service.submit_button'
                                            )}
                                        </span>
                                    </button>
                                )}
                                <button
                                    className="join-item btn btn-sm app-button-toolbar"
                                    onClick={() =>
                                        handleStartStopService(
                                            serviceId,
                                            status?.isRunning ? 'stop' : 'start'
                                        )
                                    }
                                >
                                    {status?.isRunning ? (
                                        <>
                                            <Icon className="app-icon" path={mdiStop} size={0.7} />
                                            <span>
                                                {t(
                                                    'views.Settings.Service.stop_service.submit_button'
                                                )}
                                            </span>
                                        </>
                                    ) : (
                                        <>
                                            <Icon className="app-icon" path={mdiPlay} size={0.7} />
                                            <span>
                                                {t(
                                                    'views.Settings.Service.start_service.submit_button'
                                                )}
                                            </span>
                                        </>
                                    )}
                                </button>
                            </div>
                            <div className="join app-service-actions">
                                {config.constraints.length > 0 && (
                                    <button
                                        className="join-item btn btn-sm app-button-toolbar"
                                        onClick={() => {
                                            setSelectedServiceId(serviceId);
                                            setIsConfigureModalOpen(true);
                                        }}
                                    >
                                        <Icon
                                            className="app-icon"
                                            path={mdiSquareEditOutline}
                                            size={0.7}
                                        />
                                        <span>
                                            {t('views.Settings.Service.edit_service.submit_button')}
                                        </span>
                                    </button>
                                )}
                                <button
                                    className="join-item btn btn-sm app-button-toolbar"
                                    onClick={() => handleResetServiceConfig(serviceId)}
                                >
                                    <Icon className="app-icon" path={mdiBackupRestore} size={0.7} />
                                    <span>
                                        {t('views.Settings.Service.reset_service.submit_button')}
                                    </span>
                                </button>
                            </div>
                        </StatusCard>
                    );
                })}

            <ServiceConfigModal
                serviceConfig={serviceConfig}
                serviceStatus={serviceStatus}
                configure={configure}
            />
        </div>
    );
};
