import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';
import { useLayoutStore } from '../../../stores/layout';
import { useRetentionStore } from '../../../stores/retention';
import { RetentionModal } from './RetentionModal';

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

    const { retention, setRetention } = useRetentionStore();
    const [retentionModalOpen, setRetentionModalOpen] = useState(false);
    const handleRetentionModalOpen = useCallback(() => setRetentionModalOpen(true), []);
    const handleRetentionModalClose = useCallback(() => setRetentionModalOpen(false), []);
    const handleRetentionSubmit = useCallback(
        (value: number) => {
            setRetention(value);
            setRetentionModalOpen(false);
            sendUserAlert(t('views.Settings.Preference.realtime_retention.success'));
        },
        [t, setRetention]
    );

    return (
        <div className="app-settings-content-padded">
            <div className="app-settings-action-bordered">
                <div className="app-stack-sm">
                    <h3 className="app-subheading">
                        {t('views.Settings.Preference.realtime_retention.title')}
                    </h3>
                    <p className="app-description">
                        {t('views.Settings.Preference.realtime_retention.description')}
                    </p>
                </div>
                <button
                    className="btn app-button-primary app-button-responsive"
                    onClick={handleRetentionModalOpen}
                >
                    {t('views.Settings.Preference.realtime_retention.edit_button')}
                </button>
            </div>
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
            {retentionModalOpen && (
                <RetentionModal
                    retention={retention}
                    onClose={handleRetentionModalClose}
                    onSubmit={handleRetentionSubmit}
                />
            )}
        </div>
    );
};
