import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import {
    usePurgeHelicorderFilesMutation,
    usePurgeMiniSeedFilesMutation,
    usePurgeSeisRecordsMutation,
    useRestoreServiceConfigMutation,
    useRestoreStationConfigMutation
} from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';

export const Dangerous = () => {
    const { t } = useTranslation();

    const [[purgeSeisRecords], [purgeMiniSeedFiles], [purgeHelicorderFiles]] = [
        usePurgeSeisRecordsMutation(),
        usePurgeMiniSeedFilesMutation(),
        usePurgeHelicorderFilesMutation()
    ];
    const [[resetStationConfig], [resetServiceConfig]] = [
        useRestoreStationConfigMutation(),
        useRestoreServiceConfigMutation()
    ];

    const actions = useMemo(
        () => [
            {
                title: t('views.Settings.Dangerous.purge_waveform_records.title'),
                description: t('views.Settings.Dangerous.purge_waveform_records.description'),
                buttonText: t('views.Settings.Dangerous.purge_waveform_records.submit_button'),
                confirmMessage: t(
                    'views.Settings.Dangerous.purge_waveform_records.confirm_message'
                ),
                confirmBtnText: t('views.Settings.Dangerous.purge_waveform_records.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.purge_waveform_records.cancel_button'),
                onConfirmed: async () =>
                    await sendPromiseAlert(
                        purgeSeisRecords(),
                        t('views.Settings.Dangerous.purge_waveform_records.purging'),
                        t('views.Settings.Dangerous.purge_waveform_records.success'),
                        (error) =>
                            t('views.Settings.Dangerous.purge_waveform_records.error', { error })
                    )
            },
            {
                title: t('views.Settings.Dangerous.purge_miniseed_files.title'),
                description: t('views.Settings.Dangerous.purge_miniseed_files.description'),
                buttonText: t('views.Settings.Dangerous.purge_miniseed_files.submit_button'),
                confirmMessage: t('views.Settings.Dangerous.purge_miniseed_files.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.purge_miniseed_files.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.purge_miniseed_files.cancel_button'),
                onConfirmed: async () =>
                    await sendPromiseAlert(
                        purgeMiniSeedFiles(),
                        t('views.Settings.Dangerous.purge_miniseed_files.purging'),
                        t('views.Settings.Dangerous.purge_miniseed_files.success'),
                        (error) =>
                            t('views.Settings.Dangerous.purge_miniseed_files.error', { error })
                    )
            },
            {
                title: t('views.Settings.Dangerous.purge_helicorder_files.title'),
                description: t('views.Settings.Dangerous.purge_helicorder_files.description'),
                buttonText: t('views.Settings.Dangerous.purge_helicorder_files.submit_button'),
                confirmMessage: t(
                    'views.Settings.Dangerous.purge_helicorder_files.confirm_message'
                ),
                confirmBtnText: t('views.Settings.Dangerous.purge_helicorder_files.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.purge_helicorder_files.cancel_button'),
                onConfirmed: async () =>
                    await sendPromiseAlert(
                        purgeHelicorderFiles(),
                        t('views.Settings.Dangerous.purge_helicorder_files.purging'),
                        t('views.Settings.Dangerous.purge_helicorder_files.success'),
                        (error) =>
                            t('views.Settings.Dangerous.purge_helicorder_files.error', { error })
                    )
            },
            {
                title: t('views.Settings.Dangerous.reset_station_config.title'),
                description: t('views.Settings.Dangerous.reset_station_config.description'),
                buttonText: t('views.Settings.Dangerous.reset_station_config.submit_button'),
                confirmMessage: t('views.Settings.Dangerous.reset_station_config.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.reset_station_config.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.reset_station_config.cancel_button'),
                onConfirmed: async () =>
                    await sendPromiseAlert(
                        resetStationConfig(),
                        t('views.Settings.Dangerous.reset_station_config.resetting'),
                        t('views.Settings.Dangerous.reset_station_config.success'),
                        (error) =>
                            t('views.Settings.Dangerous.reset_station_config.error', { error })
                    )
            },
            {
                title: t('views.Settings.Dangerous.reset_service_config.title'),
                description: t('views.Settings.Dangerous.reset_service_config.description'),
                buttonText: t('views.Settings.Dangerous.reset_service_config.submit_button'),
                confirmMessage: t('views.Settings.Dangerous.reset_service_config.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.reset_service_config.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.reset_service_config.cancel_button'),
                onConfirmed: async () =>
                    await sendPromiseAlert(
                        resetServiceConfig(),
                        t('views.Settings.Dangerous.reset_service_config.resetting'),
                        t('views.Settings.Dangerous.reset_service_config.success'),
                        (error) =>
                            t('views.Settings.Dangerous.reset_service_config.error', { error })
                    )
            }
        ],
        [
            t,
            purgeHelicorderFiles,
            purgeMiniSeedFiles,
            purgeSeisRecords,
            resetServiceConfig,
            resetStationConfig
        ]
    );

    return (
        <div className="mx-auto max-w-3xl space-y-4 p-6">
            {actions.map(
                (
                    {
                        title,
                        description,
                        confirmBtnText,
                        cancelBtnText,
                        confirmMessage,
                        buttonText,
                        onConfirmed
                    },
                    index
                ) => (
                    <div
                        key={`${index}-${title}`}
                        className="flex flex-col justify-between gap-4 border-b border-gray-300 p-4 sm:flex-row sm:items-center"
                    >
                        <div className="space-y-2">
                            <h3 className="text-lg font-semibold text-gray-800">{title}</h3>
                            <p className="text-sm text-gray-500">{description}</p>
                        </div>
                        <button
                            className="btn w-full rounded-lg bg-red-500 px-4 py-2 font-medium text-white shadow-lg transition-all hover:bg-red-700 sm:w-auto"
                            onClick={() =>
                                sendUserConfirm(confirmMessage, {
                                    title: buttonText,
                                    cancelBtnText,
                                    confirmBtnText,
                                    onConfirmed
                                })
                            }
                        >
                            {buttonText}
                        </button>
                    </div>
                )
            )}
        </div>
    );
};
