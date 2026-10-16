import { GridValidRowModel } from '@mui/x-data-grid';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import {
    useCreateUserMutation,
    useGetUserListQuery,
    useRemoveUserMutation,
    useUpdateUserMutation
} from '../../../graphql';
import { sendPromiseAlert } from '../../../helpers/alert/sendPromiseAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';

export const useUsers = () => {
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

    return {
        getUserListDataError,
        getUserListDataLoading,
        searchUser,
        setSearchUser,
        filteredUserList,
        handleRefreshUserList,
        selectedUsers,
        handleSelectUsers,
        handleRemoveUser,
        handleRemoveUsers,
        createUserModalOpen,
        setCreateUserModalOpen,
        createUserStatus,
        setCreateUserStatus,
        handleCreateUser,
        editUserModalOpen,
        setEditUserModalOpen,
        editUserStatus,
        setEditUserStatus,
        targetUserInEditModal,
        setTargetUserInEditModal,
        handleUpdateUser
    };
};
