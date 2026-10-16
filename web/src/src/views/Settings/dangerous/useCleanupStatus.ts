import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { JobStatus, useGetCleanupStatusQuery } from '../../../graphql';
import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { cleanupTargetByKind } from './cleanup';

export const useCleanupStatus = () => {
    const { t } = useTranslation();
    const [trackedCleanupJobId, setTrackedCleanupJobId] = useState<string | null>(null);
    const {
        data: cleanupStatusData,
        error: cleanupStatusError,
        refetch: refetchCleanupStatus
    } = useGetCleanupStatusQuery({ fetchPolicy: 'network-only', pollInterval: 5000 });
    const cleanupJob = cleanupStatusData?.getCleanupStatus;
    const isCleanupRunning = cleanupJob?.status === JobStatus.Running;
    const runningTarget = cleanupTargetByKind(cleanupJob?.kind);

    const getCleanupKindLabel = useCallback(
        (kind?: string) => {
            switch (kind) {
                case 'seis_records_full':
                    return t('views.Settings.Dangerous.purge_data.kind_seis_records_full');
                case 'seis_records_by_date':
                    return t('views.Settings.Dangerous.purge_data.kind_seis_records_by_date');
                case 'miniseed_full':
                    return t('views.Settings.Dangerous.purge_data.kind_miniseed_full');
                case 'miniseed_by_date':
                    return t('views.Settings.Dangerous.purge_data.kind_miniseed_by_date');
                case 'helicorder_full':
                    return t('views.Settings.Dangerous.purge_data.kind_helicorder_full');
                case 'helicorder_by_date':
                    return t('views.Settings.Dangerous.purge_data.kind_helicorder_by_date');
                default:
                    return t('views.Settings.Dangerous.purge_data.kind_unknown');
            }
        },
        [t]
    );

    const getCleanupStatusMeta = useCallback(() => {
        if (cleanupStatusError !== undefined) {
            return {
                status: 'error' as const,
                message: t('views.Settings.Dangerous.purge_data.status_query_failed', {
                    error: cleanupStatusError.message
                })
            };
        }

        if (
            cleanupJob === null ||
            cleanupJob === undefined ||
            cleanupJob.status === JobStatus.Idle
        ) {
            return {
                status: 'ok' as const,
                message: t('views.Settings.Dangerous.purge_data.status_idle')
            };
        }

        const kind = getCleanupKindLabel(cleanupJob.kind);
        if (cleanupJob.status === JobStatus.Running) {
            return {
                status: 'warning' as const,
                message: `${t('views.Settings.Dangerous.purge_data.status_running', {
                    kind
                })} ${t('views.Settings.Dangerous.purge_data.running_hint')}`
            };
        }
        if (cleanupJob.status === JobStatus.Succeeded) {
            return {
                status: 'ok' as const,
                message: t('views.Settings.Dangerous.purge_data.status_success', { kind })
            };
        }
        return {
            status: 'error' as const,
            message: t('views.Settings.Dangerous.purge_data.status_failed', {
                kind,
                error: cleanupJob.error ?? 'unknown error'
            })
        };
    }, [cleanupJob, cleanupStatusError, getCleanupKindLabel, t]);

    useEffect(() => {
        if (
            trackedCleanupJobId === null ||
            cleanupJob === null ||
            cleanupJob === undefined ||
            cleanupJob.id !== trackedCleanupJobId
        ) {
            return;
        }

        const kind = getCleanupKindLabel(cleanupJob.kind);
        if (cleanupJob.status === JobStatus.Succeeded) {
            sendUserAlert(t('views.Settings.Dangerous.purge_data.status_success', { kind }));
            setTrackedCleanupJobId(null);
        }
        if (cleanupJob.status === JobStatus.Failed) {
            sendUserAlert(
                t('views.Settings.Dangerous.purge_data.status_failed', {
                    kind,
                    error: cleanupJob.error ?? 'unknown error'
                }),
                true,
                5000
            );
            setTrackedCleanupJobId(null);
        }
    }, [cleanupJob, getCleanupKindLabel, t, trackedCleanupJobId]);

    return {
        isCleanupRunning,
        runningTarget,
        refetchCleanupStatus,
        getCleanupStatusMeta,
        setTrackedCleanupJobId
    };
};
