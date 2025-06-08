import {
    mdiAlertCircle,
    mdiApplicationCogOutline,
    mdiBackupRestore,
    mdiCheckCircle,
    mdiPlay,
    mdiRestart,
    mdiSquareEditOutline,
    mdiStop
} from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../components/DialogModal';
import { ErrorPage } from '../../components/ErrorPage';
import { StatusCard } from '../../components/StatusCard';
import { InputType, TypedInput } from '../../components/TypedInput';
import {
    useGetServiceDataQuery,
    useRestartServiceMutation,
    useRestoreServiceConfigMutation,
    useStartServiceMutation,
    useStopServiceMutation,
    useUpdateServiceConfigMutation
} from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { getTimeString } from '../../helpers/utils/getTimeString';

export const Service = () => {
    const { t } = useTranslation();

    const {
        data: getServiceDataData,
        refetch: getServiceDataRefetch,
        error: getServiceDataError,
        loading: getServiceDataLoading
    } = useGetServiceDataQuery({ pollInterval: 5000 });

    const [serviceConfig, setServiceConfig] = useState<
        Record<
            string,
            {
                serviceName: string;
                constraints: Array<{
                    key: string;
                    name: string;
                    description: string;
                    configType: string;
                    isRequired: boolean;
                    currentValue: unknown;
                    options?: Record<string, unknown>;
                }>;
            }
        >
    >({});
    useEffect(() => {
        if (getServiceDataData?.getServiceConfigConstraint) {
            const configMap: typeof serviceConfig = {};
            getServiceDataData.getServiceConfigConstraint.forEach((service) => {
                const validConstraints = service.constraints.filter(
                    (constraint): constraint is NonNullable<typeof constraint> =>
                        constraint !== null
                );
                configMap[service.serviceId] = {
                    serviceName: service.serviceName,
                    constraints: validConstraints.map((constraint) => ({
                        key: constraint.key,
                        name: constraint.name,
                        description: constraint.description,
                        configType: constraint.configType,
                        isRequired: constraint.isRequired,
                        currentValue: constraint.currentValue,
                        options: constraint.options
                    }))
                };
            });
            setServiceConfig(configMap);
        }
    }, [getServiceDataData?.getServiceConfigConstraint]);

    const [serviceStatus, setServiceStatus] = useState<
        Record<
            string,
            {
                description: string;
                serviceName: string;
                restarts: number;
                startedAt: number;
                stoppedAt: number;
                updatedAt: number;
                isRunning: boolean;
            }
        >
    >({});
    useEffect(() => {
        if (getServiceDataData?.getServiceStatus) {
            const statusMap: typeof serviceStatus = {};
            getServiceDataData.getServiceStatus.forEach((service) => {
                statusMap[service.serviceId] = {
                    description: service.description,
                    serviceName: service.name,
                    restarts: service.restarts,
                    startedAt: service.startedAt,
                    stoppedAt: service.stoppedAt,
                    updatedAt: service.updatedAt,
                    isRunning: service.isRunning
                };
            });
            setServiceStatus(statusMap);
        }
    }, [getServiceDataData?.getServiceStatus]);

    const [restartService] = useRestartServiceMutation();
    const handleRestartService = useCallback(
        (serviceId: string) => {
            sendUserConfirm(t('views.Settings.Service.restart_service.confirm_message'), {
                title: t('views.Settings.Service.restart_service.confirm_title'),
                cancelBtnText: t('views.Settings.Service.restart_service.cancel_button'),
                confirmBtnText: t('views.Settings.Service.restart_service.submit_button'),
                onConfirmed: async () => {
                    await sendPromiseAlert(
                        restartService({ variables: { serviceId } }),
                        t('views.Settings.Service.restart_service.restarting'),
                        t('views.Settings.Service.restart_service.success'),
                        (error) => t('views.Settings.Service.restart_service.error', { error })
                    );
                    await getServiceDataRefetch();
                }
            });
        },
        [t, getServiceDataRefetch, restartService]
    );

    const [[startService], [stopService]] = [useStartServiceMutation(), useStopServiceMutation()];
    const handleStartStopService = useCallback(
        (serviceId: string, action: 'start' | 'stop') => {
            switch (action) {
                case 'start':
                    sendUserConfirm(t('views.Settings.Service.start_service.confirm_message'), {
                        title: t('views.Settings.Service.start_service.confirm_title'),
                        cancelBtnText: t('views.Settings.Service.start_service.cancel_button'),
                        confirmBtnText: t('views.Settings.Service.start_service.submit_button'),
                        onConfirmed: async () => {
                            await sendPromiseAlert(
                                startService({ variables: { serviceId } }),
                                t('views.Settings.Service.start_service.starting'),
                                t('views.Settings.Service.start_service.success'),
                                (error) =>
                                    t('views.Settings.Service.start_service.error', { error })
                            );
                            await getServiceDataRefetch();
                        }
                    });
                    break;
                case 'stop':
                    sendUserConfirm(t('views.Settings.Service.stop_service.confirm_message'), {
                        title: t('views.Settings.Service.stop_service.confirm_title'),
                        cancelBtnText: t('views.Settings.Service.stop_service.cancel_button'),
                        confirmBtnText: t('views.Settings.Service.stop_service.submit_button'),
                        onConfirmed: async () => {
                            await sendPromiseAlert(
                                stopService({ variables: { serviceId } }),
                                t('views.Settings.Service.stop_service.stopping'),
                                t('views.Settings.Service.stop_service.success'),
                                (error) => t('views.Settings.Service.stop_service.error', { error })
                            );
                            await getServiceDataRefetch();
                        }
                    });
                    break;
            }
        },
        [t, getServiceDataRefetch, startService, stopService]
    );

    const [isConfigureModalOpen, setIsConfigureModalOpen] = useState(false);
    const [configureModalStatus, setConfigureModalStatus] = useState<{
        message: string;
        status?: 'updating' | 'success' | 'error';
    }>({ message: '' });
    const [selectedServiceId, setSelectedServiceId] = useState('');
    const [updateServiceConfigMutation] = useUpdateServiceConfigMutation();
    const handleUpdateServiceConfig = useCallback(
        async (serviceId: string, key: string, val: unknown, isRequired: boolean) => {
            if (isRequired && !val && typeof val !== 'boolean') {
                if (val !== 0 && typeof val === 'number') {
                    setConfigureModalStatus({ status: 'error', message: 'This field is required' });
                    return;
                }
            }
            try {
                setConfigureModalStatus({
                    status: 'updating',
                    message: t('views.Settings.Service.edit_service.updating')
                });
                await updateServiceConfigMutation({ variables: { serviceId, key, val } });
                await getServiceDataRefetch();
                setConfigureModalStatus({
                    status: 'success',
                    message: t('views.Settings.Service.edit_service.success')
                });
            } catch (error) {
                setConfigureModalStatus({
                    status: 'error',
                    message: t('views.Settings.Service.edit_service.error', { error })
                });
            }
            await new Promise((resolve) => setTimeout(resolve, 5000));
            setConfigureModalStatus({ message: '' });
        },
        [t, getServiceDataRefetch, updateServiceConfigMutation]
    );

    const [restoreServiceConfig] = useRestoreServiceConfigMutation();
    const handleResetServiceConfig = useCallback(
        async (serviceId: string) => {
            const requestFn = async () => {
                await restoreServiceConfig({ variables: { serviceId } });
                await getServiceDataRefetch();
            };
            sendUserConfirm(t('views.Settings.Service.reset_service.confirm_message'), {
                title: t('views.Settings.Service.reset_service.confirm_title'),
                cancelBtnText: t('views.Settings.Service.reset_service.cancel_button'),
                confirmBtnText: t('views.Settings.Service.reset_service.submit_button'),
                onConfirmed: async () =>
                    await sendPromiseAlert(
                        requestFn(),
                        t('views.Settings.Service.reset_service.resetting'),
                        t('views.Settings.Service.reset_service.success'),
                        (error) => t('views.Settings.Service.reset_service.error', { error })
                    )
            });
        },
        [t, getServiceDataRefetch, restoreServiceConfig]
    );

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
        <div className="flex items-center justify-center">
            <span className="loading loading-spinner text-primary mt-6" />
        </div>
    ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
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
                            <div className="join my-1 flex justify-center">
                                {status?.isRunning && (
                                    <button
                                        className="join-item btn btn-sm flex items-center"
                                        onClick={() => handleRestartService(serviceId)}
                                    >
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiRestart}
                                            size={0.7}
                                        />
                                        <span>
                                            {t(
                                                'views.Settings.Service.restart_service.submit_button'
                                            )}
                                        </span>
                                    </button>
                                )}
                                <button
                                    className="join-item btn btn-sm flex items-center"
                                    onClick={() =>
                                        handleStartStopService(
                                            serviceId,
                                            status?.isRunning ? 'stop' : 'start'
                                        )
                                    }
                                >
                                    {status?.isRunning ? (
                                        <>
                                            <Icon
                                                className="flex-shrink-0"
                                                path={mdiStop}
                                                size={0.7}
                                            />
                                            <span>
                                                {t(
                                                    'views.Settings.Service.stop_service.submit_button'
                                                )}
                                            </span>
                                        </>
                                    ) : (
                                        <>
                                            <Icon
                                                className="flex-shrink-0"
                                                path={mdiPlay}
                                                size={0.7}
                                            />
                                            <span>
                                                {t(
                                                    'views.Settings.Service.start_service.submit_button'
                                                )}
                                            </span>
                                        </>
                                    )}
                                </button>
                            </div>
                            <div className="join my-1 flex justify-center">
                                {config.constraints.length > 0 && (
                                    <button
                                        className="join-item btn btn-sm flex items-center"
                                        onClick={() => {
                                            setSelectedServiceId(serviceId);
                                            setIsConfigureModalOpen(true);
                                        }}
                                    >
                                        <Icon
                                            className="flex-shrink-0"
                                            path={mdiSquareEditOutline}
                                            size={0.7}
                                        />
                                        <span>
                                            {t('views.Settings.Service.edit_service.submit_button')}
                                        </span>
                                    </button>
                                )}
                                <button
                                    className="join-item btn btn-sm flex items-center"
                                    onClick={() => handleResetServiceConfig(serviceId)}
                                >
                                    <Icon
                                        className="flex-shrink-0"
                                        path={mdiBackupRestore}
                                        size={0.7}
                                    />
                                    <span>
                                        {t('views.Settings.Service.reset_service.submit_button')}
                                    </span>
                                </button>
                            </div>
                        </StatusCard>
                    );
                })}

            {serviceConfig[selectedServiceId] && (
                <DialogModal
                    fullScreen
                    heading={
                        <div className="text-gray-800">
                            <h2 className="flex items-center space-x-2 py-4 text-lg font-extrabold">
                                <Icon
                                    className="flex-shrink-0"
                                    path={mdiSquareEditOutline}
                                    size={1}
                                />
                                <span>
                                    {t('views.Settings.Service.edit_service_modal.title', {
                                        name: serviceConfig[selectedServiceId].serviceName
                                    })}
                                </span>
                            </h2>
                            <p className="text-sm">
                                {serviceStatus[selectedServiceId].description}
                            </p>
                        </div>
                    }
                    open={isConfigureModalOpen}
                    onClose={() => {
                        setIsConfigureModalOpen(false);
                        setConfigureModalStatus({ message: '' });
                    }}
                >
                    {configureModalStatus.status === 'updating' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-gray-500">
                            <span className="loading loading-spinner loading-xs" />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    {configureModalStatus.status === 'error' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-red-500">
                            <Icon className="flex-shrink-0" path={mdiAlertCircle} size={0.8} />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    {configureModalStatus.status === 'success' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-green-500">
                            <Icon className="flex-shrink-0" path={mdiCheckCircle} size={0.8} />
                            <span>{configureModalStatus.message}</span>
                        </div>
                    )}
                    <div className="flex flex-wrap justify-center gap-4 lg:justify-start">
                        {serviceConfig[selectedServiceId].constraints.map((item, index) => (
                            <fieldset
                                className="fieldset w-md"
                                key={`${selectedServiceId}-${item.key}-${index}`}
                            >
                                <legend className="fieldset-legend">
                                    <span className="font-bold text-gray-700">{item.name}</span>
                                    {item.isRequired && <span className="text-red-500">*</span>}
                                </legend>
                                <TypedInput
                                    key={`${item.key}-${item.currentValue}`}
                                    dataType={item.configType as InputType}
                                    placeholder={item.name}
                                    defaultValue={item.currentValue as string}
                                    onSubmit={(value) =>
                                        handleUpdateServiceConfig(
                                            selectedServiceId,
                                            item.key,
                                            value,
                                            item.isRequired
                                        )
                                    }
                                    fieldName={item.name}
                                    options={item.options}
                                />
                                <p className="w-[80%] text-xs text-gray-500">{item.description}</p>
                            </fieldset>
                        ))}
                    </div>
                </DialogModal>
            )}
        </div>
    );
};
