import { mdiAlertCircle, mdiKeyPlus, mdiKeyRemove, mdiRefresh } from '@mdi/js';
import Icon from '@mdi/react';
import { GridColDef, GridValidRowModel } from '@mui/x-data-grid';
import { FormEvent, useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import { ErrorPage } from '../../../components/ui/ErrorPage';
import { CodeBlock } from '../../../components/widget/CodeBlock';
import { TableList } from '../../../components/widget/TableList';
import { TimePicker } from '../../../components/widget/TimePicker';
import {
    FdsnAccessKey,
    useCreateFdsnAccessKeyMutation,
    useGetFdsnAccessKeysQuery,
    useRevokeFdsnAccessKeysMutation
} from '../../../graphql';
import { sendPromiseAlert } from '../../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';
import { getTimeString } from '../../../helpers/utils/getTimeString';

export const FDSN = ({ currentLocale }: { currentLocale: string }) => {
    const { t } = useTranslation();
    const { data, loading, error, refetch } = useGetFdsnAccessKeysQuery({
        fetchPolicy: 'network-only'
    });
    const [createKey, { loading: creating }] = useCreateFdsnAccessKeyMutation({
        fetchPolicy: 'no-cache'
    });
    const [revokeKeys] = useRevokeFdsnAccessKeysMutation();
    const [revoking, setRevoking] = useState(false);
    const [selectedKeyIds, setSelectedKeyIds] = useState<string[]>([]);
    const [hideRevoked, setHideRevoked] = useState(true);
    const [remark, setRemark] = useState('');
    const [expiresAt, setExpiresAt] = useState<number>();
    const [noExpiration, setNoExpiration] = useState(false);
    const [createModalOpen, setCreateModalOpen] = useState(false);
    const [createError, setCreateError] = useState('');
    const [secret, setSecret] = useState<string | null>(null);
    const [refreshError, setRefreshError] = useState('');
    const enabled = data?.getStationConfig?.fdsn_service_enabled === true;
    const visibleKeys = (data?.getFDSNAccessKeys ?? []).filter(
        (key) => !hideRevoked || !key.revokedAt
    );
    const selectedKeys = (data?.getFDSNAccessKeys ?? []).filter(
        (key) => selectedKeyIds.includes(key.keyId) && !key.revokedAt
    );

    const handleRevoke = useCallback(
        (keyIds: string[], batch = false) => {
            if (revoking || keyIds.length === 0) {
                return;
            }
            sendUserConfirm(
                batch
                    ? t('views.Settings.FDSN.revoke_selected_confirm', { count: keyIds.length })
                    : t('views.Settings.FDSN.revoke_confirm'),
                {
                    title: t('views.Settings.FDSN.revoke'),
                    cancelBtnText: t('views.Settings.FDSN.cancel'),
                    confirmBtnText: t('views.Settings.FDSN.revoke'),
                    onConfirmed: async () => {
                        setRevoking(true);
                        try {
                            const request = async () => {
                                const result = await revokeKeys({ variables: { keyIds } });
                                if (!result.data?.revokeFDSNAccessKeys) {
                                    throw new Error(t('views.Settings.FDSN.revoke_failed'));
                                }
                                setSelectedKeyIds((ids) =>
                                    ids.filter((id) => !keyIds.includes(id))
                                );
                                await refetch();
                            };
                            await sendPromiseAlert(
                                request(),
                                t('views.Settings.FDSN.revoking'),
                                batch
                                    ? t('views.Settings.FDSN.revoke_selected_success', {
                                          count: keyIds.length
                                      })
                                    : t('views.Settings.FDSN.revoke_success'),
                                (error) => String(error)
                            );
                        } finally {
                            setRevoking(false);
                        }
                    }
                }
            );
        },
        [t, revoking, revokeKeys, refetch]
    );

    const handleCreate = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (creating) {
            return;
        }
        setCreateError('');
        setRefreshError('');
        if (!enabled) {
            setCreateError(t('views.Settings.FDSN.disabled'));
            return;
        }
        const expiration = noExpiration ? null : (expiresAt ?? 0);
        if (expiration !== null && (!Number.isFinite(expiration) || expiration <= Date.now())) {
            setCreateError(t('views.Settings.FDSN.invalid_expiration'));
            return;
        }
        try {
            const result = await createKey({ variables: { remark, expiresAt: expiration } });
            if (!result.data) {
                throw new Error(t('views.Settings.FDSN.create_failed'));
            }
            setCreateModalOpen(false);
            setSecret(result.data.createFDSNAccessKey.key);
            setRemark('');
            setExpiresAt(undefined);
            setNoExpiration(false);
            try {
                await refetch();
            } catch (error) {
                setRefreshError(
                    t('views.Settings.FDSN.refresh_after_create_failed', {
                        error: error instanceof Error ? error.message : String(error)
                    })
                );
            }
        } catch (error) {
            setCreateError(error instanceof Error ? error.message : String(error));
        }
    };

    const columns = useMemo<GridColDef<FdsnAccessKey>[]>(
        () => [
            {
                field: 'keyPrefix',
                headerName: t('views.Settings.FDSN.key'),
                minWidth: 170
            },
            {
                field: 'createdAt',
                headerName: t('views.Settings.FDSN.created_at'),
                minWidth: 230,
                renderCell: ({ row }) => getTimeString(row.createdAt)
            },
            {
                field: 'expiresAt',
                headerName: t('views.Settings.FDSN.expires_at'),
                minWidth: 230,
                renderCell: ({ row }) =>
                    row.expiresAt
                        ? getTimeString(row.expiresAt)
                        : t('views.Settings.FDSN.no_expiration')
            },
            {
                field: 'remark',
                headerName: t('views.Settings.FDSN.remark'),
                minWidth: 220
            },
            {
                field: 'status',
                headerName: t('views.Settings.FDSN.status'),
                minWidth: 150,
                valueGetter: (_value, row) =>
                    row.revokedAt
                        ? t('views.Settings.FDSN.revoked')
                        : row.expiresAt && row.expiresAt <= Date.now()
                          ? t('views.Settings.FDSN.expired')
                          : t('views.Settings.FDSN.active')
            },
            {
                field: 'revokedAt',
                headerName: t('views.Settings.FDSN.revoked_at'),
                minWidth: 230,
                renderCell: ({ row }) => (row.revokedAt ? getTimeString(row.revokedAt) : 'N/A')
            },
            {
                field: 'lastRequestedAt',
                headerName: t('views.Settings.FDSN.last_requested_at'),
                minWidth: 230,
                renderCell: ({ row }) =>
                    row.lastRequestedAt
                        ? getTimeString(row.lastRequestedAt)
                        : t('views.Settings.FDSN.never')
            },
            {
                field: 'userIp',
                headerName: t('views.Settings.FDSN.user_ip'),
                minWidth: 180
            },
            {
                field: 'userAgent',
                headerName: 'User-Agent',
                minWidth: 400
            },
            {
                field: 'actions',
                headerName: t('views.Settings.FDSN.actions'),
                minWidth: 150,
                sortable: false,
                resizable: false,
                headerAlign: 'center',
                align: 'center',
                renderCell: ({ row }) => (
                    <div className="app-action-icons">
                        <button
                            className="app-delete-action disabled:cursor-not-allowed disabled:opacity-50"
                            disabled={!!row.revokedAt || revoking}
                            onClick={() => handleRevoke([row.keyId])}
                        >
                            {t('views.Settings.FDSN.revoke')}
                        </button>
                    </div>
                )
            }
        ],
        [t, revoking, handleRevoke]
    );

    if (error && !data) {
        return <ErrorPage content={error.message} debug={JSON.stringify(error)} />;
    }
    if (loading && !data) {
        return (
            <div className="app-center">
                <span className="loading loading-spinner app-loading-delayed" />
            </div>
        );
    }

    return (
        <>
            <div className="app-toolbar">
                <button
                    className="btn btn-sm app-button-toolbar"
                    onClick={() => refetch().catch((error) => sendUserAlert(String(error), true))}
                >
                    <Icon className="app-icon" path={mdiRefresh} size={0.7} />
                    <span>{t('views.Settings.FDSN.refresh')}</span>
                </button>
                <button
                    className="btn btn-sm app-button-toolbar"
                    disabled={!enabled || creating}
                    onClick={() => {
                        setCreateError('');
                        setCreateModalOpen(true);
                    }}
                >
                    <Icon className="app-icon" path={mdiKeyPlus} size={0.7} />
                    <span>{t('views.Settings.FDSN.create')}</span>
                </button>
                {selectedKeys.length > 0 && (
                    <button
                        className="btn btn-sm app-button-toolbar"
                        disabled={revoking}
                        onClick={() =>
                            handleRevoke(
                                selectedKeys.map((key) => key.keyId),
                                true
                            )
                        }
                    >
                        <Icon className="app-icon" path={mdiKeyRemove} size={0.7} />
                        <span>{t('views.Settings.FDSN.revoke_selected')}</span>
                    </button>
                )}
            </div>
            <div className="app-row">
                <input
                    id="fdsn-hide-revoked"
                    type="checkbox"
                    checked={hideRevoked}
                    onChange={({ currentTarget }) => setHideRevoked(currentTarget.checked)}
                />
                <label className="label app-checkbox-label" htmlFor="fdsn-hide-revoked">
                    {t('views.Settings.FDSN.hide_revoked')}
                </label>
            </div>
            {!enabled && <p className="app-description">{t('views.Settings.FDSN.disabled')}</p>}
            <p className="app-description">{t('views.Settings.FDSN.usage')}</p>
            <TableList
                currentLocale={currentLocale}
                data={visibleKeys.map((key) => ({ ...key, id: key.keyId }))}
                columns={columns as GridColDef<GridValidRowModel>[]}
                sortField="createdAt"
                sortDirection="desc"
                onSelect={(selected) => setSelectedKeyIds(selected.map((key) => key.keyId))}
                rowSelectionModel={selectedKeys.map((key) => key.keyId)}
                isRowSelectable={({ row }) => !row.revokedAt}
            />
            <DialogModal
                open={createModalOpen}
                onClose={() => {
                    setCreateModalOpen(false);
                    setCreateError('');
                    setRemark('');
                    setExpiresAt(undefined);
                    setNoExpiration(false);
                }}
                heading={
                    <h2 className="app-dialog-heading-colored">
                        <Icon className="app-icon" path={mdiKeyPlus} size={1} />
                        <span>{t('views.Settings.FDSN.create')}</span>
                    </h2>
                }
            >
                <form className="app-stack" onSubmit={handleCreate}>
                    {createError && (
                        <div className="app-form-status-error" role="alert">
                            <Icon className="app-icon" path={mdiAlertCircle} size={0.8} />
                            <span>{createError}</span>
                        </div>
                    )}
                    {!noExpiration && (
                        <div className="app-stack-xs">
                            <label
                                className="label fieldset-legend app-form-legend"
                                htmlFor="fdsn-key-expiration"
                            >
                                {t('views.Settings.FDSN.expires_at')}
                            </label>
                            <TimePicker
                                id="fdsn-key-expiration"
                                className="input app-input app-full-width"
                                currentLocale={currentLocale}
                                placeholder={t('views.Settings.FDSN.expiration_placeholder')}
                                value={expiresAt}
                                onChange={setExpiresAt}
                            />
                            <p className="app-field-help">
                                {t('views.Settings.FDSN.expiration_help')}
                            </p>
                        </div>
                    )}
                    <div className="app-row">
                        <input
                            id="fdsn-key-no-expiration"
                            type="checkbox"
                            checked={noExpiration}
                            onChange={({ currentTarget }) => setNoExpiration(currentTarget.checked)}
                        />
                        <label
                            className="label app-checkbox-label"
                            htmlFor="fdsn-key-no-expiration"
                        >
                            {t('views.Settings.FDSN.no_expiration')}
                        </label>
                    </div>
                    <div className="app-stack-xs">
                        <label
                            className="label fieldset-legend app-form-legend"
                            htmlFor="fdsn-key-remark"
                        >
                            {t('views.Settings.FDSN.remark')}
                        </label>
                        <input
                            id="fdsn-key-remark"
                            className="input app-input app-full-width"
                            placeholder={t('views.Settings.FDSN.remark_placeholder')}
                            value={remark}
                            maxLength={500}
                            onChange={(event) => setRemark(event.currentTarget.value)}
                        />
                    </div>
                    <button
                        className="btn app-button-primary app-full-width"
                        type="submit"
                        disabled={!enabled || creating}
                    >
                        {creating && <span className="loading loading-spinner loading-xs" />}
                        <span>{t('views.Settings.FDSN.create')}</span>
                    </button>
                </form>
            </DialogModal>
            <DialogModal
                open={secret !== null}
                onClose={() => {
                    setSecret(null);
                    setRefreshError('');
                }}
                heading={
                    <h2 className="app-dialog-heading-colored">
                        {t('views.Settings.FDSN.created')}
                    </h2>
                }
            >
                <div className="app-stack">
                    <p className="app-help-text">{t('views.Settings.FDSN.secret_help')}</p>
                    {refreshError && (
                        <div className="app-form-status-error" role="alert">
                            <Icon className="app-icon" path={mdiAlertCircle} size={0.8} />
                            <span>{refreshError}</span>
                        </div>
                    )}
                    <CodeBlock fileName="fdsn-access-key.txt">{secret ?? ''}</CodeBlock>
                </div>
            </DialogModal>
        </>
    );
};
