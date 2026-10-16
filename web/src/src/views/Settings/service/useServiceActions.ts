import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import {
    useRestartServiceMutation,
    useRestoreServiceConfigMutation,
    useStartServiceMutation,
    useStopServiceMutation
} from '../../../graphql';
import { sendPromiseAlert } from '../../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';

export const useServiceActions = (getServiceDataRefetch: () => Promise<unknown>) => {
    const { t } = useTranslation();
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

    return { handleRestartService, handleStartStopService, handleResetServiceConfig };
};
