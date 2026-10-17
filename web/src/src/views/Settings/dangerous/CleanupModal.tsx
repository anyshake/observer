import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../../components/ui/DialogModal';
import { TimePicker } from '../../../components/widget/TimePicker';
import type { useCleanup } from './useCleanup';

export const CleanupModal = ({
    currentLocale,
    cleanup
}: {
    currentLocale: string;
    cleanup: ReturnType<typeof useCleanup>;
}) => {
    const { t } = useTranslation();
    const {
        cleanupModalTarget,
        cleanupModalTargetId,
        cleanupModalMode,
        cleanupRanges,
        cleanupVerificationText,
        cleanupVerificationInput,
        cleanupVerificationPassed,
        handleCleanupModalClose,
        handleCleanupModeAllClick,
        handleCleanupModeRangeClick,
        handleCleanupStartDateChange,
        handleCleanupEndDateChange,
        handleCleanupVerificationInputChange,
        handleCleanupSubmit
    } = cleanup;

    return (
        <>
            {cleanupModalTarget !== undefined && cleanupModalTargetId !== null && (
                <DialogModal
                    open={true}
                    nativeModal={false}
                    onClose={handleCleanupModalClose}
                    heading={
                        <h2 className="app-dialog-title">{cleanupModalTarget.confirmTitle}</h2>
                    }
                >
                    <div className="app-stack-lg">
                        <p className="app-help-text">{cleanupModalTarget.description}</p>

                        <div className="app-cleanup-mode">
                            <span className="app-form-label">
                                {t('views.Settings.Dangerous.purge_data.mode_label')}
                            </span>
                            <div className="join app-fit-width">
                                <button
                                    className={`btn join-item btn-sm ${
                                        cleanupModalMode === 'all' ? 'btn-active' : ''
                                    }`}
                                    onClick={handleCleanupModeAllClick}
                                >
                                    {t('views.Settings.Dangerous.purge_data.mode_all')}
                                </button>
                                <button
                                    className={`btn join-item btn-sm ${
                                        cleanupModalMode === 'range' ? 'btn-active' : ''
                                    }`}
                                    onClick={handleCleanupModeRangeClick}
                                >
                                    {t('views.Settings.Dangerous.purge_data.mode_range')}
                                </button>
                            </div>
                        </div>

                        {cleanupModalMode === 'range' && (
                            <div className="app-cleanup-range">
                                <label className="app-stack-xs">
                                    <span className="app-cleanup-date-label">
                                        {t('views.Settings.Dangerous.purge_data.start_date')}
                                    </span>
                                    <TimePicker
                                        dateOnly
                                        currentLocale={currentLocale}
                                        value={cleanupRanges[cleanupModalTargetId].startDate}
                                        placeholder={t(
                                            'views.Settings.Dangerous.purge_data.start_date'
                                        )}
                                        className="app-time-input-padded"
                                        onChange={handleCleanupStartDateChange}
                                    />
                                </label>

                                <label className="app-stack-xs">
                                    <span className="app-cleanup-date-label">
                                        {t('views.Settings.Dangerous.purge_data.end_date')}
                                    </span>
                                    <TimePicker
                                        dateOnly
                                        currentLocale={currentLocale}
                                        value={cleanupRanges[cleanupModalTargetId].endDate}
                                        placeholder={t(
                                            'views.Settings.Dangerous.purge_data.end_date'
                                        )}
                                        className="app-time-input-padded"
                                        onChange={handleCleanupEndDateChange}
                                    />
                                </label>
                            </div>
                        )}

                        <div className="app-stack-sm">
                            <label className="app-form-label">
                                {t('views.Settings.Dangerous.purge_data.verification_label')}
                            </label>
                            <p className="app-help-text">
                                {t('views.Settings.Dangerous.purge_data.verification_help', {
                                    verificationText: cleanupVerificationText
                                })}
                            </p>
                            <input
                                className="input app-input app-full-width"
                                value={cleanupVerificationInput}
                                placeholder={cleanupVerificationText}
                                onChange={handleCleanupVerificationInputChange}
                            />
                        </div>

                        <div className="app-modal-actions">
                            <button
                                className="btn app-button-cancel"
                                onClick={handleCleanupModalClose}
                            >
                                {cleanupModalTarget.cancelBtnText}
                            </button>
                            <button
                                className="btn app-button-danger"
                                disabled={!cleanupVerificationPassed}
                                onClick={handleCleanupSubmit}
                            >
                                {cleanupModalTarget.confirmBtnText}
                            </button>
                        </div>
                    </div>
                </DialogModal>
            )}
        </>
    );
};
