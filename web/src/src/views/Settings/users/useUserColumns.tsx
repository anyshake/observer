import { GridValidRowModel } from '@mui/x-data-grid';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { getTimeString } from '../../../helpers/utils/getTimeString';
import type { useUsers } from './useUsers';

export const useUserColumns = (
    users: Pick<
        ReturnType<typeof useUsers>,
        'handleRemoveUser' | 'setTargetUserInEditModal' | 'setEditUserModalOpen'
    >
) => {
    const { t } = useTranslation();
    const { handleRemoveUser, setTargetUserInEditModal, setEditUserModalOpen } = users;
    const columns = useMemo(
        () => [
            {
                field: 'id',
                headerName: t('views.Settings.Users.user_list_columns.id'),
                hideable: false,
                sortable: true,
                minWidth: 120
            },
            {
                field: 'userId',
                headerName: t('views.Settings.Users.user_list_columns.user_id'),
                hideable: false,
                sortable: true,
                minWidth: 230
            },
            {
                field: 'username',
                headerName: t('views.Settings.Users.user_list_columns.username'),
                hideable: false,
                sortable: true,
                minWidth: 200
            },
            {
                field: 'admin',
                headerName: t('views.Settings.Users.user_list_columns.admin'),
                hideable: true,
                sortable: true,
                minWidth: 150,
                renderCell: ({ value }: GridValidRowModel) =>
                    value
                        ? t('views.Settings.Users.user_list_columns.yes')
                        : t('views.Settings.Users.user_list_columns.no')
            },
            {
                field: 'createdAt',
                headerName: t('views.Settings.Users.user_list_columns.created_at'),
                hideable: true,
                sortable: true,
                minWidth: 230,
                renderCell: (cell: GridValidRowModel) => getTimeString(cell.value)
            },
            {
                field: 'lastLogin',
                headerName: t('views.Settings.Users.user_list_columns.last_login'),
                hideable: true,
                sortable: true,
                minWidth: 230,
                renderCell: ({ value }) =>
                    value ? getTimeString(value) : t('views.Settings.Users.user_list_columns.never')
            },
            {
                field: 'userIp',
                headerName: t('views.Settings.Users.user_list_columns.user_ip'),
                hideable: true,
                sortable: true,
                minWidth: 180,
                renderCell: ({ value }) => (value ? value : 'N/A')
            },
            {
                field: 'userAgent',
                headerName: t('views.Settings.Users.user_list_columns.user_agent'),
                hideable: true,
                sortable: true,
                minWidth: 500,
                renderCell: ({ value }) => (value ? value : 'N/A')
            },
            {
                field: 'updatedAt',
                headerName: t('views.Settings.Users.user_list_columns.updated_at'),
                hideable: true,
                sortable: true,
                minWidth: 250,
                renderCell: ({ value }) =>
                    value ? getTimeString(value) : t('views.Settings.Users.user_list_columns.never')
            },
            {
                field: 'actions',
                headerName: t('views.Settings.Users.user_list_columns.actions'),
                sortable: false,
                resizable: false,
                minWidth: 150,
                headerAlign: 'center',
                align: 'center',
                renderCell: ({ row: { userId, username, admin } }: GridValidRowModel) => (
                    <div className="app-action-icons">
                        <button
                            className="app-link-action"
                            onClick={() => {
                                setTargetUserInEditModal({ userId, username, admin });
                                setEditUserModalOpen(true);
                            }}
                        >
                            {t('views.Settings.Users.edit_user.submit_button')}
                        </button>
                        <button
                            className="app-delete-action"
                            onClick={() => handleRemoveUser(userId, username)}
                        >
                            {t('views.Settings.Users.remove_user.submit_button')}
                        </button>
                    </div>
                )
            }
        ],
        [t, handleRemoveUser, setTargetUserInEditModal, setEditUserModalOpen]
    );

    return { columns };
};
