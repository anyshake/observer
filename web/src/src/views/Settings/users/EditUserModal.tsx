import { mdiAccountEdit, mdiAlertCircle, mdiCheckCircle } from '@mdi/js';
import Icon from '@mdi/react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import type { useUsers } from './useUsers';

export const EditUserModal = ({ users }: { users: ReturnType<typeof useUsers> }) => {
    const { t } = useTranslation();
    const {
        editUserModalOpen,
        setEditUserModalOpen,
        editUserStatus,
        setEditUserStatus,
        targetUserInEditModal,
        setTargetUserInEditModal,
        handleUpdateUser
    } = users;

    return (
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
                        placeholder={t('views.Settings.Users.edit_user_modal.username.placeholder')}
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
                        placeholder={t('views.Settings.Users.edit_user_modal.password.placeholder')}
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
    );
};
