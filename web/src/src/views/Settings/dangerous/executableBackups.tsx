import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import {
    RemoveExecutableBackupsMutation,
    useGetExecutableBackupsQuery,
    useRemoveExecutableBackupsMutation
} from '../../../graphql';
import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { getTimeString } from '../../../helpers/utils/getTimeString';
import { DangerousAction } from './action';

export const ExecutableBackups = () => {
    const { t } = useTranslation();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedNames, setSelectedNames] = useState<string[]>([]);
    const [removalResults, setRemovalResults] = useState<
        RemoveExecutableBackupsMutation['removeExecutableBackups'] | null
    >(null);
    const { data, error, loading, refetch } = useGetExecutableBackupsQuery({
        fetchPolicy: 'network-only',
        skip: !isModalOpen,
        notifyOnNetworkStatusChange: true
    });
    const [removeBackups, { loading: removing }] = useRemoveExecutableBackupsMutation();
    const backups = data?.getExecutableBackups ?? [];
    const failedRemovals = removalResults?.filter(({ removed }) => !removed) ?? [];
    const removalSummary = removalResults
        ? t('views.Settings.Dangerous.executable_backups.result_summary', {
              removed: removalResults.length - failedRemovals.length,
              failed: failedRemovals.length
          })
        : '';

    useEffect(() => {
        if (data?.getExecutableBackups) {
            const available = new Set(data.getExecutableBackups.map(({ name }) => name));
            setSelectedNames((names) => names.filter((name) => available.has(name)));
        }
    }, [data]);

    const handleOpen = () => {
        setSelectedNames([]);
        setRemovalResults(null);
        setIsModalOpen(true);
    };

    const handleClose = () => {
        if (removing) {
            return;
        }
        setIsModalOpen(false);
        setSelectedNames([]);
        setRemovalResults(null);
    };

    const handleDelete = async () => {
        const names = [...selectedNames];
        if (names.length === 0 || loading || removing || error) {
            return;
        }
        setRemovalResults(null);
        try {
            const response = await removeBackups({ variables: { names } });
            if (!response.data) {
                throw new Error(t('views.Settings.Dangerous.executable_backups.empty_response'));
            }
            const results = response.data.removeExecutableBackups;
            const failed = results.filter(({ removed }) => !removed);
            setRemovalResults(results);
            setSelectedNames(failed.map(({ name }) => name));
            sendUserAlert(
                t('views.Settings.Dangerous.executable_backups.result_summary', {
                    removed: results.length - failed.length,
                    failed: failed.length
                }),
                failed.length > 0,
                5000
            );
        } catch (error) {
            sendUserAlert(
                t('views.Settings.Dangerous.executable_backups.remove_failed', {
                    error: error instanceof Error ? error.message : String(error)
                }),
                true,
                5000
            );
        } finally {
            void refetch().catch(() => undefined);
        }
    };

    return (
        <>
            <DangerousAction
                title={t('views.Settings.Dangerous.executable_backups.title')}
                description={t('views.Settings.Dangerous.executable_backups.description')}
                buttonText={t('views.Settings.Dangerous.executable_backups.submit_button')}
                disabled={removing}
                onClick={handleOpen}
            />

            {isModalOpen && (
                <DialogModal
                    open={true}
                    nativeModal={false}
                    onClose={handleClose}
                    heading={
                        <h2 className="app-dialog-title">
                            {t('views.Settings.Dangerous.executable_backups.confirm_title')}
                        </h2>
                    }
                >
                    <div className="app-stack-lg">
                        <p className="app-help-text">
                            {t('views.Settings.Dangerous.executable_backups.description')}
                        </p>

                        <div className="app-backup-toolbar">
                            <label className="app-backup-select-all">
                                <input
                                    type="checkbox"
                                    className="checkbox checkbox-sm"
                                    checked={
                                        backups.length > 0 &&
                                        selectedNames.length === backups.length
                                    }
                                    disabled={
                                        loading || removing || !!error || backups.length === 0
                                    }
                                    onChange={({ currentTarget }) =>
                                        setSelectedNames(
                                            currentTarget.checked
                                                ? backups.map(({ name }) => name)
                                                : []
                                        )
                                    }
                                />
                                {t('views.Settings.Dangerous.executable_backups.select_all')}
                            </label>
                            <button
                                type="button"
                                className="btn btn-sm"
                                disabled={loading || removing}
                                onClick={() => void refetch().catch(() => undefined)}
                            >
                                {t('views.Settings.Dangerous.executable_backups.refresh_button')}
                            </button>
                        </div>

                        {error && (
                            <p className="app-backup-error" role="alert">
                                {t('views.Settings.Dangerous.executable_backups.list_failed', {
                                    error: error.message
                                })}
                            </p>
                        )}
                        {loading && (
                            <p className="app-description" role="status">
                                {t('views.Settings.Dangerous.executable_backups.loading')}
                            </p>
                        )}
                        {removalResults && (
                            <div
                                className={`app-backup-removal-status ${failedRemovals.length > 0 ? 'app-text-danger' : 'app-text-body'}`}
                                role="status"
                            >
                                <p>{removalSummary}</p>
                                {failedRemovals.length > 0 && (
                                    <ul className="app-backup-removal-list">
                                        {failedRemovals.map(({ name, error }) => (
                                            <li key={name} className="app-break-all">
                                                <span className="app-text-mono">{name}</span>:{' '}
                                                {error}
                                            </li>
                                        ))}
                                    </ul>
                                )}
                            </div>
                        )}
                        {!loading && !error && backups.length === 0 && (
                            <p className="app-description">
                                {t('views.Settings.Dangerous.executable_backups.empty')}
                            </p>
                        )}
                        {!loading && !error && backups.length > 0 && (
                            <div className="app-backup-list">
                                {backups.map(({ name, size, modifiedAt }) => (
                                    <label key={name} className="app-backup-item">
                                        <input
                                            type="checkbox"
                                            className="checkbox checkbox-sm app-no-shrink"
                                            checked={selectedNames.includes(name)}
                                            disabled={removing}
                                            onChange={({ currentTarget }) =>
                                                setSelectedNames((previous) =>
                                                    currentTarget.checked
                                                        ? [...previous, name]
                                                        : previous.filter(
                                                              (selected) => selected !== name
                                                          )
                                                )
                                            }
                                        />
                                        <span className="app-backup-item-content">
                                            <span className="app-backup-name">{name}</span>
                                            <span className="app-backup-date">
                                                {`${(size / 1024 / 1024).toFixed(2)} MB · ${getTimeString(modifiedAt)}`}
                                            </span>
                                        </span>
                                    </label>
                                ))}
                            </div>
                        )}

                        {selectedNames.length > 0 && (
                            <p className="app-help-text">
                                {t('views.Settings.Dangerous.executable_backups.confirm_message', {
                                    count: selectedNames.length
                                })}
                            </p>
                        )}
                        <div className="app-modal-actions">
                            <button
                                type="button"
                                className="btn app-button-cancel"
                                disabled={removing}
                                onClick={handleClose}
                            >
                                {t('views.Settings.Dangerous.executable_backups.cancel_button')}
                            </button>
                            <button
                                type="button"
                                className="btn app-button-danger"
                                disabled={
                                    selectedNames.length === 0 || loading || removing || !!error
                                }
                                onClick={() => void handleDelete()}
                            >
                                {removing && (
                                    <span className="loading loading-spinner loading-xs" />
                                )}
                                {t('views.Settings.Dangerous.executable_backups.delete_button', {
                                    count: selectedNames.length
                                })}
                            </button>
                        </div>
                    </div>
                </DialogModal>
            )}
        </>
    );
};
