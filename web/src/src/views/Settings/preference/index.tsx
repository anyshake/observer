import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';
import { useLayoutStore } from '../../../stores/layout';

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
        <div className="app-settings-content-padded">
            <div className="app-settings-action-bordered">
                <div className="app-stack-sm">
                    <h3 className="app-subheading">
                        {t('views.Settings.Preference.reset_layout_preference.title')}
                    </h3>
                    <p className="app-description">
                        {t('views.Settings.Preference.reset_layout_preference.description')}
                    </p>
                </div>
                <button className="btn app-preference-reset" onClick={handleClearPreference}>
                    {t('views.Settings.Preference.reset_layout_preference.submit_button')}
                </button>
            </div>
        </div>
    );
};
