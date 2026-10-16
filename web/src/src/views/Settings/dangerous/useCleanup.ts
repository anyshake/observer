import { type ChangeEvent, type MouseEvent, useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import {
    PurgeDataJob,
    usePurgeHelicorderFilesByDateMutation,
    usePurgeHelicorderFilesMutation,
    usePurgeMiniSeedFilesByDateMutation,
    usePurgeMiniSeedFilesMutation,
    usePurgeSeisRecordsByDateMutation,
    usePurgeSeisRecordsMutation
} from '../../../graphql';
import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import {
    type CleanupMode,
    type CleanupRange,
    type CleanupTargetId,
    getErrorMessage,
    initialCleanupModes,
    initialCleanupRanges
} from './cleanup';
import { useCleanupStatus } from './useCleanupStatus';

export const useCleanup = () => {
    const { t } = useTranslation();
    const {
        isCleanupRunning,
        runningTarget,
        refetchCleanupStatus,
        getCleanupStatusMeta,
        setTrackedCleanupJobId
    } = useCleanupStatus();
    const [cleanupModes, setCleanupModes] =
        useState<Record<CleanupTargetId, CleanupMode>>(initialCleanupModes);
    const [cleanupRanges, setCleanupRanges] =
        useState<Record<CleanupTargetId, CleanupRange>>(initialCleanupRanges);
    const [cleanupModalTargetId, setCleanupModalTargetId] = useState<CleanupTargetId | null>(null);
    const [cleanupVerificationInput, setCleanupVerificationInput] = useState('');

    const [purgeSeisRecords] = usePurgeSeisRecordsMutation();
    const [purgeSeisRecordsByDate] = usePurgeSeisRecordsByDateMutation();
    const [purgeMiniSeedFiles] = usePurgeMiniSeedFilesMutation();
    const [purgeMiniSeedFilesByDate] = usePurgeMiniSeedFilesByDateMutation();
    const [purgeHelicorderFiles] = usePurgeHelicorderFilesMutation();
    const [purgeHelicorderFilesByDate] = usePurgeHelicorderFilesByDateMutation();

    const cleanupTargets = useMemo(
        () => [
            {
                id: 'seisRecords',
                title: t('views.Settings.Dangerous.purge_waveform_records.title'),
                description: t('views.Settings.Dangerous.purge_waveform_records.description'),
                buttonText: t('views.Settings.Dangerous.purge_waveform_records.submit_button'),
                confirmTitle: t('views.Settings.Dangerous.purge_waveform_records.confirm_title'),
                confirmMessage: t(
                    'views.Settings.Dangerous.purge_waveform_records.confirm_message'
                ),
                confirmBtnText: t('views.Settings.Dangerous.purge_waveform_records.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.purge_waveform_records.cancel_button')
            },
            {
                id: 'miniSeed',
                title: t('views.Settings.Dangerous.purge_miniseed_files.title'),
                description: t('views.Settings.Dangerous.purge_miniseed_files.description'),
                buttonText: t('views.Settings.Dangerous.purge_miniseed_files.submit_button'),
                confirmTitle: t('views.Settings.Dangerous.purge_miniseed_files.confirm_title'),
                confirmMessage: t('views.Settings.Dangerous.purge_miniseed_files.confirm_message'),
                confirmBtnText: t('views.Settings.Dangerous.purge_miniseed_files.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.purge_miniseed_files.cancel_button')
            },
            {
                id: 'helicorder',
                title: t('views.Settings.Dangerous.purge_helicorder_files.title'),
                description: t('views.Settings.Dangerous.purge_helicorder_files.description'),
                buttonText: t('views.Settings.Dangerous.purge_helicorder_files.submit_button'),
                confirmTitle: t('views.Settings.Dangerous.purge_helicorder_files.confirm_title'),
                confirmMessage: t(
                    'views.Settings.Dangerous.purge_helicorder_files.confirm_message'
                ),
                confirmBtnText: t('views.Settings.Dangerous.purge_helicorder_files.confirm_button'),
                cancelBtnText: t('views.Settings.Dangerous.purge_helicorder_files.cancel_button')
            }
        ],
        [t]
    );

    const getTargetLabel = useCallback(
        (targetId: CleanupTargetId) => {
            const target = cleanupTargets.find(({ id }) => id === targetId);
            return target?.title ?? targetId;
        },
        [cleanupTargets]
    );

    const updateCleanupMode = useCallback((targetId: CleanupTargetId, mode: CleanupMode) => {
        setCleanupModes((prev) => ({ ...prev, [targetId]: mode }));
    }, []);

    const updateCleanupRange = useCallback(
        (targetId: CleanupTargetId, key: keyof CleanupRange, value: number) => {
            setCleanupRanges((prev) => ({
                ...prev,
                [targetId]: { ...prev[targetId], [key]: value || undefined }
            }));
        },
        []
    );

    const getCleanupVerificationText = useCallback(
        (targetId: CleanupTargetId) => getTargetLabel(targetId),
        [getTargetLabel]
    );

    const executeCleanup = useCallback(
        async (targetId: CleanupTargetId) => {
            const mode = cleanupModes[targetId];
            const range = cleanupRanges[targetId];
            const startDate = range.startDate;
            const endDate = range.endDate;
            const variables = {
                startDate: startDate ?? 0,
                endDate: endDate ?? 0
            };
            if (mode === 'range') {
                if (startDate === undefined || endDate === undefined) {
                    sendUserAlert(
                        t('views.Settings.Dangerous.purge_data.invalid_range'),
                        true,
                        4000
                    );
                    return;
                }
                if (startDate > endDate) {
                    sendUserAlert(
                        t('views.Settings.Dangerous.purge_data.invalid_range'),
                        true,
                        4000
                    );
                    return;
                }
            }

            try {
                let job: PurgeDataJob | null | undefined;
                if (targetId === 'seisRecords') {
                    job =
                        mode === 'all'
                            ? (await purgeSeisRecords()).data?.purgeSeisRecords
                            : (await purgeSeisRecordsByDate({ variables })).data
                                  ?.purgeSeisRecordsByDate;
                }
                if (targetId === 'miniSeed') {
                    job =
                        mode === 'all'
                            ? (await purgeMiniSeedFiles()).data?.purgeMiniSeedFiles
                            : (await purgeMiniSeedFilesByDate({ variables })).data
                                  ?.purgeMiniSeedFilesByDate;
                }
                if (targetId === 'helicorder') {
                    job =
                        mode === 'all'
                            ? (await purgeHelicorderFiles()).data?.purgeHelicorderFiles
                            : (await purgeHelicorderFilesByDate({ variables })).data
                                  ?.purgeHelicorderFilesByDate;
                }

                if (job === null || job === undefined) {
                    return;
                }

                setTrackedCleanupJobId(job.id);
                await refetchCleanupStatus();
                sendUserAlert(t('views.Settings.Dangerous.purge_data.started'), false, 3000);
            } catch (error) {
                sendUserAlert(
                    t('views.Settings.Dangerous.purge_data.request_failed', {
                        error: getErrorMessage(error)
                    }),
                    true,
                    5000
                );
            }
        },
        [
            cleanupModes,
            cleanupRanges,
            purgeHelicorderFiles,
            purgeHelicorderFilesByDate,
            purgeMiniSeedFiles,
            purgeMiniSeedFilesByDate,
            purgeSeisRecords,
            purgeSeisRecordsByDate,
            refetchCleanupStatus,
            setTrackedCleanupJobId,
            t
        ]
    );

    const confirmCleanup = useCallback(
        (targetId: CleanupTargetId) => {
            const target = cleanupTargets.find(({ id }) => id === targetId);
            if (target === undefined) {
                return;
            }

            const mode = cleanupModes[targetId];
            const range = cleanupRanges[targetId];
            const startDate = range.startDate;
            const endDate = range.endDate;
            if (mode === 'range') {
                if (startDate === undefined || endDate === undefined) {
                    sendUserAlert(
                        t('views.Settings.Dangerous.purge_data.invalid_range'),
                        true,
                        4000
                    );
                    return;
                }
                if (startDate > endDate) {
                    sendUserAlert(
                        t('views.Settings.Dangerous.purge_data.invalid_range'),
                        true,
                        4000
                    );
                    return;
                }
            }

            setCleanupModalTargetId(null);
            setCleanupVerificationInput('');
            executeCleanup(targetId);
        },
        [cleanupModes, cleanupRanges, cleanupTargets, executeCleanup, t]
    );

    const handleCleanupTargetClick = useCallback(
        ({ currentTarget }: MouseEvent<HTMLButtonElement>) => {
            setCleanupVerificationInput('');
            setCleanupModalTargetId(currentTarget.value as CleanupTargetId);
        },
        []
    );

    const handleCleanupModalClose = useCallback(() => {
        setCleanupModalTargetId(null);
    }, []);

    const handleCleanupModeAllClick = useCallback(() => {
        if (cleanupModalTargetId !== null) {
            updateCleanupMode(cleanupModalTargetId, 'all');
        }
    }, [cleanupModalTargetId, updateCleanupMode]);

    const handleCleanupModeRangeClick = useCallback(() => {
        if (cleanupModalTargetId !== null) {
            updateCleanupMode(cleanupModalTargetId, 'range');
        }
    }, [cleanupModalTargetId, updateCleanupMode]);

    const handleCleanupStartDateChange = useCallback(
        (value: number) => {
            if (cleanupModalTargetId !== null) {
                updateCleanupRange(cleanupModalTargetId, 'startDate', value);
            }
        },
        [cleanupModalTargetId, updateCleanupRange]
    );

    const handleCleanupEndDateChange = useCallback(
        (value: number) => {
            if (cleanupModalTargetId !== null) {
                updateCleanupRange(cleanupModalTargetId, 'endDate', value);
            }
        },
        [cleanupModalTargetId, updateCleanupRange]
    );

    const handleCleanupVerificationInputChange = useCallback(
        ({ currentTarget }: ChangeEvent<HTMLInputElement>) => {
            setCleanupVerificationInput(currentTarget.value);
        },
        []
    );

    const handleCleanupSubmit = useCallback(() => {
        if (cleanupModalTargetId !== null) {
            confirmCleanup(cleanupModalTargetId);
        }
    }, [cleanupModalTargetId, confirmCleanup]);

    const cleanupModalTarget =
        cleanupModalTargetId === null
            ? undefined
            : cleanupTargets.find(({ id }) => id === cleanupModalTargetId);
    const cleanupModalMode =
        cleanupModalTargetId === null ? undefined : cleanupModes[cleanupModalTargetId];
    const cleanupStatusMeta = getCleanupStatusMeta();
    const cleanupVerificationText =
        cleanupModalTargetId === null ? '' : getCleanupVerificationText(cleanupModalTargetId);
    const cleanupVerificationPassed = cleanupVerificationInput === cleanupVerificationText;

    return {
        cleanupTargets,
        isCleanupRunning,
        runningTarget,
        cleanupStatusMeta,
        handleCleanupTargetClick,
        cleanupModalTarget,
        cleanupModalTargetId,
        cleanupModalMode,
        cleanupRanges,
        cleanupVerificationText,
        cleanupVerificationInput,
        cleanupVerificationPassed,
        handleCleanupModalClose,
        handleCleanupModeAllClick,
        handleCleanupModeRangeClick,
        handleCleanupStartDateChange,
        handleCleanupEndDateChange,
        handleCleanupVerificationInputChange,
        handleCleanupSubmit
    };
};
