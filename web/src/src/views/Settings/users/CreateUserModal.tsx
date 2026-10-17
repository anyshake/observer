import { mdiAccountPlus, mdiAlertCircle, mdiCheckCircle } from '@mdi/js';
import Icon from '@mdi/react';
import { Field, Form, Formik } from 'formik';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import type { useUsers } from './useUsers';

export const CreateUserModal = ({ users }: { users: ReturnType<typeof useUsers> }) => {
    const { t } = useTranslation();
    const {
        createUserModalOpen,
        setCreateUserModalOpen,
        createUserStatus,
        setCreateUserStatus,
        handleCreateUser
    } = users;

    return (
        <DialogModal
            open={createUserModalOpen}
            heading={
                <h2 className="app-dialog-heading-colored">
                    <Icon className="app-icon" path={mdiAccountPlus} size={1} />
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
                    <Form className="app-stack">
                        {createUserStatus.status === 'creating' && (
                            <div className="app-form-status">
                                <span className="loading loading-spinner loading-xs" />
                                <span>{createUserStatus.message}</span>
                            </div>
                        )}
                        {createUserStatus.status === 'error' && (
                            <div className="app-form-status-error">
                                <Icon className="app-icon" path={mdiAlertCircle} size={0.8} />
                                <span>{createUserStatus.message}</span>
                            </div>
                        )}
                        {createUserStatus.status === 'success' && (
                            <div className="app-form-status-success">
                                <Icon className="app-icon" path={mdiCheckCircle} size={0.8} />
                                <span>{createUserStatus.message}</span>
                            </div>
                        )}

                        <div className="app-stack-xs">
                            <label className="label fieldset-legend app-form-legend">
                                {t('views.Settings.Users.create_user_modal.username.label')}
                            </label>
                            <Field
                                required
                                name="username"
                                type="text"
                                autoComplete="off"
                                className="input app-input app-full-width"
                                placeholder={t(
                                    'views.Settings.Users.create_user_modal.username.placeholder'
                                )}
                            />
                        </div>

                        <div className="app-stack-xs">
                            <label className="label fieldset-legend app-form-legend">
                                {t('views.Settings.Users.create_user_modal.password.label')}
                            </label>
                            <Field
                                required
                                name="password"
                                type="password"
                                className="input app-input app-full-width"
                                placeholder={t(
                                    'views.Settings.Users.create_user_modal.password.placeholder'
                                )}
                            />
                        </div>

                        <div className="app-row">
                            <Field name="admin" type="checkbox" />
                            <label className="label app-checkbox-label">
                                {t('views.Settings.Users.create_user_modal.admin.label')}
                            </label>
                        </div>

                        <button
                            className="btn app-button-primary app-full-width"
                            type="submit"
                            disabled={isSubmitting}
                        >
                            {t('views.Settings.Users.create_user_modal.create_button')}
                        </button>
                    </Form>
                )}
            </Formik>
        </DialogModal>
    );
};
