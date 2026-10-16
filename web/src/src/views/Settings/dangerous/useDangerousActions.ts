import { type MouseEvent, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import {
    useRestartApplicationMutation,
    useRestoreServiceConfigMutation,
    useRestoreStationConfigMutation
} from '../../../graphql';
import { sendPromiseAlert } from '../../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';

export const useDangerousActions = () => {
    const { t } = useTranslation();
    const [[restartApplication], [resetStationConfig], [resetServiceConfig]] = [
        useRestartApplicationMutation(),
        useRestoreStationConfigMutation(),
        useRestoreServiceConfigMutation()
    ];

    const handleResetStationConfig = useCallback(
        async () =>
            await sendPromiseAlert(
                resetStationConfig(),
                t('views.Settings.Dangerous.reset_station_config.resetting'),
                t('views.Settings.Dangerous.reset_station_config.success'),
                (error) => t('views.Settings.Dangerous.reset_station_config.error', { error })
            ),
        [resetStationConfig, t]
    );

    const handleResetServiceConfig = useCallback(
        async () =>
            await sendPromiseAlert(
                resetServiceConfig(),
                t('views.Settings.Dangerous.reset_service_config.resetting'),
                t('views.Settings.Dangerous.reset_service_config.success'),
                (error) => t('views.Settings.Dangerous.reset_service_config.error', { error })
            ),
        [resetServiceConfig, t]
    );

    const handleRestartApplication = useCallback(
        async () =>
            await sendPromiseAlert(
                restartApplication(),
                t('views.Settings.Dangerous.restart_application.restarting'),
                t('views.Settings.Dangerous.restart_application.success'),
                (error) =>
                    t('views.Settings.Dangerous.restart_application.error', {
                        error
                    })
            ),
        [restartApplication, t]
    );

    const otherActions = useMemo(
        () => [
            {
                title: t('views.Settings.Dangerous.reset_station_config.title'),
                description: t('views.Settings.Dangerous.reset_station_config.description'),
                buttonText: t('views.Settings.Dangerous.reset_station_config.submit_button'),
                confirmTitle: t('views.Settings.Dangerous.reset_station_config.confirm_title'),
                confirmMessage: t('views.Settings.Dangerous.reset_station_config.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.reset_station_config.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.reset_station_config.cancel_button'),
                onConfirmed: handleResetStationConfig
            },
            {
                title: t('views.Settings.Dangerous.reset_service_config.title'),
                description: t('views.Settings.Dangerous.reset_service_config.description'),
                buttonText: t('views.Settings.Dangerous.reset_service_config.submit_button'),
                confirmTitle: t('views.Settings.Dangerous.reset_service_config.confirm_title'),
                confirmMessage: t('views.Settings.Dangerous.reset_service_config.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.reset_service_config.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.reset_service_config.cancel_button'),
                onConfirmed: handleResetServiceConfig
            },
            {
                title: t('views.Settings.Dangerous.restart_application.title'),
                description: t('views.Settings.Dangerous.restart_application.description'),
                buttonText: t('views.Settings.Dangerous.restart_application.submit_button'),
                confirmTitle: t('views.Settings.Dangerous.restart_application.confirm_title'),
                confirmMessage: t('views.Settings.Dangerous.restart_application.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.restart_application.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.restart_application.cancel_button'),
                onConfirmed: handleRestartApplication
            }
        ],
        [t, handleRestartApplication, handleResetServiceConfig, handleResetStationConfig]
    );

    const confirmOtherAction = useCallback(
        (index: number) => {
            const action = otherActions[index];
            if (action === undefined) {
                return;
            }
            sendUserConfirm(action.confirmMessage, {
                title: action.confirmTitle,
                cancelBtnText: action.cancelBtnText,
                confirmBtnText: action.confirmBtnText,
                onConfirmed: action.onConfirmed
            });
        },
        [otherActions]
    );

    const handleOtherActionClick = useCallback(
        ({ currentTarget }: MouseEvent<HTMLButtonElement>) => {
            confirmOtherAction(Number(currentTarget.value));
        },
        [confirmOtherAction]
    );

    return { otherActions, handleOtherActionClick };
};
