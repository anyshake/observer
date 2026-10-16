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
                                <Icon className="flex-shrink-0" path={mdiAlertCircle} size={0.8} />
                                <span>{createUserStatus.message}</span>
                            </div>
                        )}
                        {createUserStatus.status === 'success' && (
                            <div className="my-2 flex items-center space-x-2 text-sm text-green-500">
                                <Icon className="flex-shrink-0" path={mdiCheckCircle} size={0.8} />
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
    );
};
