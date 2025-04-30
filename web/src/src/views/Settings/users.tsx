import {
    mdiAccountEdit,
    mdiAccountPlus,
    mdiAccountRemove,
    mdiAlertCircle,
    mdiCheckCircle,
    mdiClose,
    mdiMagnify,
    mdiRefresh
} from '@mdi/js';
import Icon from '@mdi/react';
import { GridColDef, GridValidRowModel } from '@mui/x-data-grid';
import { Field, Form, Formik } from 'formik';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../components/DialogModal';
import { ErrorPage } from '../../components/ErrorPage';
import { TableList } from '../../components/TableList';
import {
    useCreateUserMutation,
    useGetUserListQuery,
    useRemoveUserMutation,
    useUpdateUserMutation
} from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { getTimeString } from '../../helpers/utils/getTimeString';

interface IUsers {
    currentLocale: string;
}

export const Users = ({ currentLocale }: IUsers) => {
    const { t } = useTranslation();

    const {
        data: getUserListData,
        refetch: getUserListDataRefetch,
        error: getUserListDataError,
        loading: getUserListDataLoading
    } = useGetUserListQuery();
    const [userList, setUserList] = useState<
        Array<{
            id: number;
            userId: string;
            username: string;
            createdAt: number;
            updatedAt: number;
            lastLogin: number;
            userIp: string;
            userAgent: string;
            admin: boolean;
        }>
    >([]);
    useEffect(() => {
        if (getUserListData?.getSysUsers) {
            setUserList(
                [...getUserListData.getSysUsers]
                    .sort((a, b) => a.createdAt - b.createdAt)
                    .map((user, index) => ({
                        id: index + 1,
                        userId: user.userId,
                        username: user.username,
                        createdAt: user.createdAt,
                        updatedAt: user.updatedAt,
                        lastLogin: user.lastLogin,
                        userIp: user.userIp,
                        userAgent: user.userAgent,
                        admin: user.admin
                    }))
            );
        }
    }, [getUserListData?.getSysUsers]);

    const [searchUser, setSearchUser] = useState('');
    const filteredUserList = useMemo(() => {
        return userList.filter(
            (user) =>
                user.username.toLowerCase().includes(searchUser.toLowerCase()) ||
                user.userId.toLowerCase().includes(searchUser.toLowerCase())
        );
    }, [searchUser, userList]);

    const handleRefreshUserList = useCallback(async () => {
        await sendPromiseAlert(
            getUserListDataRefetch(),
            t('views.Settings.Users.refresh_list.refreshing'),
            t('views.Settings.Users.refresh_list.success'),
            (error) => t('views.Settings.Users.refresh_list.error', { error })
        );
    }, [t, getUserListDataRefetch]);

    const [createUserModalOpen, setCreateUserModalOpen] = useState(false);
    const [createUserStatus, setCreateUserStatus] = useState<{
        message: string;
        status?: 'creating' | 'success' | 'error';
    }>({ message: '' });
    const [createUser] = useCreateUserMutation();
    const handleCreateUser = useCallback(
        async (username: string, password: string, admin: boolean) => {
            try {
                setCreateUserStatus({
                    status: 'creating',
                    message: t('views.Settings.Users.create_user.creating')
                });
                await createUser({ variables: { username, password, admin } });
                await getUserListDataRefetch();
                setCreateUserStatus({
                    status: 'success',
                    message: t('views.Settings.Users.create_user.success', { username })
                });
            } catch (error) {
                setCreateUserStatus({
                    status: 'error',
                    message: t('views.Settings.Users.create_user.error', { error })
                });
            }
        },
        [t, createUser, getUserListDataRefetch]
    );

    const [editUserModalOpen, setEditUserModalOpen] = useState(false);
    const [editUserStatus, setEditUserStatus] = useState<{
        message: string;
        status?: 'editing' | 'success' | 'error';
    }>({ message: '' });
    const [targetUserInEditModal, setTargetUserInEditModal] = useState<{
        username: string;
        userId: string;
        admin: boolean;
        password?: string;
    }>({ username: '', userId: '', admin: false });
    const [updateUser] = useUpdateUserMutation();
    const handleUpdateUser = useCallback(async () => {
        const { userId, username, password, admin } = targetUserInEditModal;
        try {
            setEditUserStatus({
                status: 'editing',
                message: t('views.Settings.Users.edit_user.updating')
            });
            await updateUser({
                variables: { userId, username, password, admin }
            });
            await getUserListDataRefetch();
            setEditUserStatus({
                status: 'success',
                message: t('views.Settings.Users.edit_user.success', { username })
            });
        } catch (error) {
            setEditUserStatus({
                status: 'error',
                message: t('views.Settings.Users.edit_user.error', { error })
            });
        }
    }, [t, getUserListDataRefetch, targetUserInEditModal, updateUser]);

    const [selectedUsers, setSelectedUsers] = useState<string[]>([]);
    const handleSelectUsers = useCallback(
        (selected: GridValidRowModel[]) => setSelectedUsers(selected.map((s) => s.userId)),
        []
    );
    const [removeUser] = useRemoveUserMutation();
    const handleRemoveUser = useCallback(
        (userId: string, username: string) => {
            const requestFn = async () => {
                await removeUser({ variables: { userId } });
                await getUserListDataRefetch();
            };
            sendUserConfirm(t('views.Settings.Users.remove_user.confirm_message', { username }), {
                title: t('views.Settings.Users.remove_user.confirm_title'),
                cancelBtnText: t('views.Settings.Users.remove_user.cancel_button'),
                confirmBtnText: t('views.Settings.Users.remove_user.confirm_button'),
                onConfirmed: async () => {
                    await sendPromiseAlert(
                        requestFn(),
                        t('views.Settings.Users.remove_user.removing'),
                        t('views.Settings.Users.remove_user.success', { username }),
                        (error) => t('views.Settings.Users.remove_user.error', { error })
                    );
                }
            });
        },
        [t, removeUser, getUserListDataRefetch]
    );
    const handleRemoveUsers = useCallback(async () => {
        const requestFn = async () => {
            for (const userId of selectedUsers) {
                await removeUser({ variables: { userId } });
            }
        };
        sendUserConfirm(t('views.Settings.Users.remove_users.confirm_message'), {
            title: t('views.Settings.Users.remove_users.confirm_title'),
            cancelBtnText: t('views.Settings.Users.remove_users.cancel_button'),
            confirmBtnText: t('views.Settings.Users.remove_users.confirm_button'),
            onConfirmed: async () => {
                await sendPromiseAlert(
                    requestFn(),
                    t('views.Settings.Users.remove_users.removing'),
                    t('views.Settings.Users.remove_users.success', { num: selectedUsers.length }),
                    (error) => t('views.Settings.Users.remove_users.error', { error })
                );
                await getUserListDataRefetch();
            }
        });
    }, [t, getUserListDataRefetch, removeUser, selectedUsers]);

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
                    <div className="space-x-4">
                        <button
                            className="cursor-pointer text-blue-700 hover:opacity-50"
                            onClick={() => {
                                setTargetUserInEditModal({ userId, username, admin });
                                setEditUserModalOpen(true);
                            }}
                        >
                            {t('views.Settings.Users.edit_user.submit_button')}
                        </button>
                        <button
                            className="cursor-pointer text-red-700 hover:opacity-50"
                            onClick={() => handleRemoveUser(userId, username)}
                        >
                            {t('views.Settings.Users.remove_user.submit_button')}
                        </button>
                    </div>
                )
            }
        ],
        [t, handleRemoveUser]
    );

    return getUserListDataError ? (
        <ErrorPage
            content={getUserListDataError.message}
            debug={JSON.stringify(getUserListDataError)}
        />
    ) : getUserListDataLoading ? (
        <div className="flex items-center justify-center">
            <span className="loading loading-spinner text-primary mt-6" />
        </div>
    ) : (
        <>
            <div className="flex flex-wrap items-center gap-2">
                <button
                    className="btn btn-sm flex items-center"
                    onClick={() => handleRefreshUserList()}
                >
                    <Icon className="flex-shrink-0" path={mdiRefresh} size={0.7} />
                    <span>{t('views.Settings.Users.refresh_list.submit_button')}</span>
                </button>
                <button
                    className="btn btn-sm flex items-center"
                    onClick={() => setCreateUserModalOpen(true)}
                >
                    <Icon className="flex-shrink-0" path={mdiAccountPlus} size={0.7} />
                    <span>{t('views.Settings.Users.create_user.submit_button')}</span>
                </button>
                {selectedUsers.length > 1 && (
                    <button
                        className="btn btn-sm flex items-center"
                        onClick={() => handleRemoveUsers()}
                    >
                        <Icon className="flex-shrink-0" path={mdiAccountRemove} size={0.7} />
                        <span>{t('views.Settings.Users.remove_users.submit_button')}</span>
                    </button>
                )}
            </div>
            <label className="input">
                <Icon className="flex-shrink-0" path={mdiMagnify} size={0.8} />
                <input
                    type="search"
                    className="grow"
                    value={searchUser}
                    onChange={({ currentTarget }) => setSearchUser(currentTarget.value)}
                    placeholder={t('views.Settings.Users.search_users.placeholder')}
                />
                {searchUser.length > 0 && (
                    <button
                        className="cursor-pointer opacity-60 transition-all hover:opacity-100"
                        onClick={() => setSearchUser('')}
                    >
                        <Icon className="flex-shrink-0" path={mdiClose} size={0.8} />
                    </button>
                )}
            </label>

            <TableList
                currentLocale={currentLocale}
                data={filteredUserList}
                columns={columns as GridColDef<GridValidRowModel>[]}
                onSelect={(selected) => handleSelectUsers(selected)}
            />

            <DialogModal
                open={createUserModalOpen}
                heading={
                    <h2 className="flex items-center space-x-2 text-lg font-extrabold text-gray-800">
                        <Icon className="flex-shrink-0" path={mdiAccountPlus} size={1} />
                        <span>{t('views.Settings.Users.create_user_modal.title')}</span>
                    </h2>
                }
                onClose={() => {
                    setCreateUserModalOpen(false);
                    setCreateUserStatus({ message: '' });
                }}
            >
                <Formik
                    onSubmit={async ({ username, password, admin }, { setSubmitting }) => {
                        await handleCreateUser(username, password, admin);
                        setSubmitting(false);
                    }}
                    initialValues={{ username: '', password: '', admin: false }}
                >
                    {({ isSubmitting }) => (
                        <Form className="space-y-4">
                            {createUserStatus.status === 'creating' && (
                                <div className="my-2 flex items-center space-x-2 text-sm text-gray-700">
                                    <span className="loading loading-spinner loading-xs" />
                                    <span>{createUserStatus.message}</span>
                                </div>
                            )}
                            {createUserStatus.status === 'error' && (
                                <div className="my-2 flex items-center space-x-2 text-sm text-red-500">
                                    <Icon
                                        className="flex-shrink-0"
                                        path={mdiAlertCircle}
                                        size={0.8}
                                    />
                                    <span>{createUserStatus.message}</span>
                                </div>
                            )}
                            {createUserStatus.status === 'success' && (
                                <div className="my-2 flex items-center space-x-2 text-sm text-green-500">
                                    <Icon
                                        className="flex-shrink-0"
                                        path={mdiCheckCircle}
                                        size={0.8}
                                    />
                                    <span>{createUserStatus.message}</span>
                                </div>
                            )}

                            <div className="space-y-1">
                                <label className="label fieldset-legend text-sm font-bold text-gray-700">
                                    {t('views.Settings.Users.create_user_modal.username.label')}
                                </label>
                                <Field
                                    required
                                    name="username"
                                    type="text"
                                    autoComplete="off"
                                    className="input w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                                    placeholder={t(
                                        'views.Settings.Users.create_user_modal.username.placeholder'
                                    )}
                                />
                            </div>

                            <div className="space-y-1">
                                <label className="label fieldset-legend text-sm font-bold text-gray-700">
                                    {t('views.Settings.Users.create_user_modal.password.label')}
                                </label>
                                <Field
                                    required
                                    name="password"
                                    type="password"
                                    className="input w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                                    placeholder={t(
                                        'views.Settings.Users.create_user_modal.password.placeholder'
                                    )}
                                />
                            </div>

                            <div className="flex items-center space-x-2">
                                <Field name="admin" type="checkbox" />
                                <label className="label text-sm text-gray-700">
                                    {t('views.Settings.Users.create_user_modal.admin.label')}
                                </label>
                            </div>

                            <button
                                className="btn w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                                type="submit"
                                disabled={isSubmitting}
                            >
                                {t('views.Settings.Users.create_user_modal.create_button')}
                            </button>
                        </Form>
                    )}
                </Formik>
            </DialogModal>

            <DialogModal
                open={editUserModalOpen}
                heading={
                    <h2 className="flex items-center space-x-2 text-lg font-extrabold text-gray-800">
                        <Icon className="flex-shrink-0" path={mdiAccountEdit} size={1} />
                        <span>{t('views.Settings.Users.edit_user_modal.title')}</span>
                    </h2>
                }
                onClose={() => {
                    setEditUserModalOpen(false);
                    setEditUserStatus({ message: '' });
                    setTargetUserInEditModal({ username: '', userId: '', admin: false });
                }}
            >
                <div className="space-y-4">
                    {editUserStatus.status === 'editing' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-gray-700">
                            <span className="loading loading-spinner loading-xs" />
                            <span>{editUserStatus.message}</span>
                        </div>
                    )}
                    {editUserStatus.status === 'error' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-red-500">
                            <Icon className="flex-shrink-0" path={mdiAlertCircle} size={0.8} />
                            <span>{editUserStatus.message}</span>
                        </div>
                    )}
                    {editUserStatus.status === 'success' && (
                        <div className="my-2 flex items-center space-x-2 text-sm text-green-500">
                            <Icon className="flex-shrink-0" path={mdiCheckCircle} size={0.8} />
                            <span>{editUserStatus.message}</span>
                        </div>
                    )}

                    <div className="space-y-1">
                        <label className="label fieldset-legend text-sm font-bold text-gray-700">
                            {t('views.Settings.Users.edit_user_modal.username.label')}
                        </label>
                        <input
                            key={targetUserInEditModal.userId}
                            type="text"
                            autoComplete="off"
                            className="input w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                            placeholder={t(
                                'views.Settings.Users.edit_user_modal.username.placeholder'
                            )}
                            defaultValue={targetUserInEditModal.username}
                            onChange={({ currentTarget }) =>
                                setTargetUserInEditModal({
                                    ...targetUserInEditModal,
                                    username: currentTarget.value
                                })
                            }
                        />
                    </div>

                    <div className="space-y-1">
                        <label className="label fieldset-legend text-sm font-bold text-gray-700">
                            {t('views.Settings.Users.edit_user_modal.password.label')}
                        </label>
                        <input
                            key={targetUserInEditModal.userId}
                            type="password"
                            className="input w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                            placeholder={t(
                                'views.Settings.Users.edit_user_modal.password.placeholder'
                            )}
                            onChange={({ currentTarget }) =>
                                setTargetUserInEditModal({
                                    ...targetUserInEditModal,
                                    password: currentTarget.value
                                })
                            }
                        />
                    </div>

                    <div className="flex items-center space-x-2">
                        <input
                            key={targetUserInEditModal.userId}
                            name="admin"
                            type="checkbox"
                            defaultChecked={targetUserInEditModal.admin}
                            onChange={({ currentTarget }) =>
                                setTargetUserInEditModal({
                                    ...targetUserInEditModal,
                                    admin: currentTarget.checked
                                })
                            }
                        />
                        <label className="label text-sm text-gray-700">
                            {t('views.Settings.Users.edit_user_modal.admin.label')}
                        </label>
                    </div>

                    <button
                        className="btn w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                        type="submit"
                        onClick={handleUpdateUser}
                        disabled={editUserStatus.status === 'editing'}
                    >
                        {t('views.Settings.Users.edit_user_modal.update_button')}
                    </button>
                </div>
            </DialogModal>
        </>
    );
};
