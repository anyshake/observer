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
                <h2 className="app-dialog-heading-colored">
                    <Icon className="app-icon" path={mdiAccountEdit} size={1} />
                    <span>{t('views.Settings.Users.edit_user_modal.title')}</span>
                </h2>
            }
            onClose={() => {
                setEditUserModalOpen(false);
                setEditUserStatus({ message: '' });
                setTargetUserInEditModal({ username: '', userId: '', admin: false });
            }}
        >
            <div className="app-stack">
                {editUserStatus.status === 'editing' && (
                    <div className="app-form-status">
                        <span className="loading loading-spinner loading-xs" />
                        <span>{editUserStatus.message}</span>
                    </div>
                )}
                {editUserStatus.status === 'error' && (
                    <div className="app-form-status-error">
                        <Icon className="app-icon" path={mdiAlertCircle} size={0.8} />
                        <span>{editUserStatus.message}</span>
                    </div>
                )}
                {editUserStatus.status === 'success' && (
                    <div className="app-form-status-success">
                        <Icon className="app-icon" path={mdiCheckCircle} size={0.8} />
                        <span>{editUserStatus.message}</span>
                    </div>
                )}

                <div className="app-stack-xs">
                    <label className="label fieldset-legend app-form-legend">
                        {t('views.Settings.Users.edit_user_modal.username.label')}
                    </label>
                    <input
                        key={targetUserInEditModal.userId}
                        type="text"
                        autoComplete="off"
                        className="input app-input app-full-width"
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

                <div className="app-stack-xs">
                    <label className="label fieldset-legend app-form-legend">
                        {t('views.Settings.Users.edit_user_modal.password.label')}
                    </label>
                    <input
                        key={targetUserInEditModal.userId}
                        type="password"
                        className="input app-input app-full-width"
                        placeholder={t('views.Settings.Users.edit_user_modal.password.placeholder')}
                        onChange={({ currentTarget }) =>
                            setTargetUserInEditModal({
                                ...targetUserInEditModal,
                                password: currentTarget.value
                            })
                        }
                    />
                </div>

                <div className="app-row">
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
                    <label className="label app-checkbox-label">
                        {t('views.Settings.Users.edit_user_modal.admin.label')}
                    </label>
                </div>

                <button
                    className="btn app-button-primary app-full-width"
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
