import { ChangeEvent, FormEvent, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import { RetentionConstraints } from '../../../config/constraints';
import { normalizeRetention } from '../../../helpers/utils/normalizeRetention';

interface IRetentionModal {
    readonly retention: number;
    readonly onClose: () => void;
    readonly onSubmit: (retention: number) => void;
}

const { minValue, maxValue, step } = RetentionConstraints;

export const RetentionModal = ({ retention, onClose, onSubmit }: IRetentionModal) => {
    const { t } = useTranslation();
    const [draft, setDraft] = useState(String(retention));
    const value = normalizeRetention(draft, retention);

    const handleDraftChange = useCallback(({ currentTarget }: ChangeEvent<HTMLInputElement>) => {
        setDraft(currentTarget.value);
    }, []);

    const handleDraftBlur = useCallback(() => {
        setDraft(String(value));
    }, [value]);

    const handleSubmit = useCallback(
        (event: FormEvent<HTMLFormElement>) => {
            event.preventDefault();
            onSubmit(value);
        },
        [onSubmit, value]
    );

    return (
        <DialogModal
            open={true}
            nativeModal={false}
            onClose={onClose}
            heading={
                <h2 className="app-dialog-title">
                    {t('views.Settings.Preference.realtime_retention.title')}
                </h2>
            }
        >
            <form className="app-stack-lg" noValidate onSubmit={handleSubmit}>
                <p className="app-help-text">
                    {t('views.Settings.Preference.realtime_retention.description')}
                </p>
                <div className="app-retention-fields">
                    <input
                        type="range"
                        className="range range-primary range-sm app-retention-range"
                        min={minValue}
                        max={maxValue}
                        step={step}
                        value={value}
                        onChange={handleDraftChange}
                    />
                    <label className="input app-input app-retention-input">
                        <input
                            type="number"
                            inputMode="numeric"
                            className="app-grow"
                            min={minValue}
                            max={maxValue}
                            step={step}
                            value={draft}
                            onChange={handleDraftChange}
                            onBlur={handleDraftBlur}
                        />
                        <span className="label">
                            {t('views.Settings.Preference.realtime_retention.unit')}
                        </span>
                    </label>
                </div>
                <div className="app-modal-actions">
                    <button type="button" className="btn app-button-cancel" onClick={onClose}>
                        {t('views.Settings.Preference.realtime_retention.cancel_button')}
                    </button>
                    <button type="submit" className="btn app-button-primary app-button-auto">
                        {t('views.Settings.Preference.realtime_retention.save_button')}
                    </button>
                </div>
            </form>
        </DialogModal>
    );
};
