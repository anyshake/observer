import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { useLayoutStore } from '../../stores/layout';

export const Preference = () => {
    const { t } = useTranslation();

    const { resetLayoutConfig } = useLayoutStore();
    const handleClearPreference = useCallback(() => {
        sendUserConfirm(t('views.Settings.Preference.reset_layout_preference.confirm_message'), {
            title: t('views.Settings.Preference.reset_layout_preference.confirm_title'),
            cancelBtnText: t('views.Settings.Preference.reset_layout_preference.cancel_button'),
            confirmBtnText: t('views.Settings.Preference.reset_layout_preference.confirm_button'),
            onConfirmed: () => {
                resetLayoutConfig();
                sendUserAlert(t('views.Settings.Preference.reset_layout_preference.success'));
            }
        });
    }, [t, resetLayoutConfig]);

    return (
        <div className="mx-auto max-w-3xl space-y-4 p-6">
            <div className="flex flex-col justify-between gap-4 border-b border-gray-300 p-4 sm:flex-row sm:items-center">
                <div className="space-y-2">
                    <h3 className="text-lg font-semibold text-gray-800">
                        {t('views.Settings.Preference.reset_layout_preference.title')}
                    </h3>
                    <p className="text-sm text-gray-500">
                        {t('views.Settings.Preference.reset_layout_preference.description')}
                    </p>
                </div>
                <button
                    className="btn w-full rounded-lg bg-red-500 px-4 py-2 font-medium text-white shadow-lg transition-all hover:bg-red-700 sm:w-auto"
                    onClick={handleClearPreference}
                >
                    {t('views.Settings.Preference.reset_layout_preference.submit_button')}
                </button>
            </div>
        </div>
    );
};
