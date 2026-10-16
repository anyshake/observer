import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import {
    RemoveExecutableBackupsMutation,
    useGetExecutableBackupsQuery,
    useRemoveExecutableBackupsMutation
} from '../../../graphql';
import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { DangerousAction } from './action';

export const ExecutableBackups = ({ currentLocale }: { currentLocale: string }) => {
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
                        <h2 className="text-lg font-extrabold text-gray-800">
                            {t('views.Settings.Dangerous.executable_backups.confirm_title')}
                        </h2>
                    }
                >
                    <div className="space-y-5">
                        <p className="text-sm leading-6 text-gray-500">
                            {t('views.Settings.Dangerous.executable_backups.description')}
                        </p>

                        <div className="flex flex-wrap items-center justify-between gap-3">
                            <label className="flex items-center gap-3 text-sm font-medium text-gray-700">
                                <input
                                    type="checkbox"
                                    className="checkbox checkbox-sm"
                                    checked={
                                        backups.length > 0 &&
                                        selectedNames.length === backups.length
                                    }
                                    disabled={loading || removing || !!error || backups.length === 0}
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
                            <p className="text-sm text-red-600" role="alert">
                                {t('views.Settings.Dangerous.executable_backups.list_failed', {
                                    error: error.message
                                })}
                            </p>
                        )}
                        {loading && (
                            <p className="text-sm text-gray-500" role="status">
                                {t('views.Settings.Dangerous.executable_backups.loading')}
                            </p>
                        )}
                        {removalResults && (
                            <div
                                className={`space-y-2 text-sm ${failedRemovals.length > 0 ? 'text-red-600' : 'text-gray-700'}`}
                                role="status"
                            >
                                <p>{removalSummary}</p>
                                {failedRemovals.length > 0 && (
                                    <ul className="max-h-40 space-y-1 overflow-y-auto">
                                        {failedRemovals.map(({ name, error }) => (
                                            <li key={name} className="break-all">
                                                <span className="font-mono">{name}</span>: {error}
                                            </li>
                                        ))}
                                    </ul>
                                )}
                            </div>
                        )}
                        {!loading && !error && backups.length === 0 && (
                            <p className="text-sm text-gray-500">
                                {t('views.Settings.Dangerous.executable_backups.empty')}
                            </p>
                        )}
                        {!loading && !error && backups.length > 0 && (
                            <div className="max-h-72 space-y-2 overflow-y-auto">
                                {backups.map(({ name, size, modifiedAt }) => (
                                    <label
                                        key={name}
                                        className="flex items-start gap-3 rounded-lg border border-gray-200 p-3 text-sm"
                                    >
                                        <input
                                            type="checkbox"
                                            className="checkbox checkbox-sm shrink-0"
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
                                        <span className="min-w-0 space-y-1">
                                            <span className="block break-all font-mono">{name}</span>
                                            <span className="block text-xs text-gray-500">
                                                {new Intl.NumberFormat(currentLocale).format(size)} B
                                                {' · '}
                                                {new Date(modifiedAt).toLocaleString(currentLocale)}
                                            </span>
                                        </span>
                                    </label>
                                ))}
                            </div>
                        )}

                        {selectedNames.length > 0 && (
                            <p className="text-sm leading-6 text-gray-500">
                                {t('views.Settings.Dangerous.executable_backups.confirm_message', {
                                    count: selectedNames.length
                                })}
                            </p>
                        )}
                        <div className="flex flex-col-reverse justify-end gap-2 pt-2 sm:flex-row">
                            <button
                                type="button"
                                className="btn rounded-lg px-4 py-2 font-medium transition-all hover:bg-gray-300"
                                disabled={removing}
                                onClick={handleClose}
                            >
                                {t('views.Settings.Dangerous.executable_backups.cancel_button')}
                            </button>
                            <button
                                type="button"
                                className="btn rounded-lg bg-red-500 px-4 py-2 font-medium text-white transition-all hover:bg-red-700 disabled:cursor-not-allowed disabled:bg-gray-400"
                                disabled={
                                    selectedNames.length === 0 || loading || removing || !!error
                                }
                                onClick={() => void handleDelete()}
                            >
                                {removing && <span className="loading loading-spinner loading-xs" />}
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
