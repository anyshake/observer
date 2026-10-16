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
                        <h2 className="text-lg font-extrabold text-gray-800">
                            {cleanupModalTarget.confirmTitle}
                        </h2>
                    }
                >
                    <div className="space-y-5">
                        <p className="text-sm leading-6 text-gray-500">
                            {cleanupModalTarget.description}
                        </p>

                        <div className="flex flex-col items-start space-y-2">
                            <span className="text-sm font-medium text-gray-700">
                                {t('views.Settings.Dangerous.purge_data.mode_label')}
                            </span>
                            <div className="join w-fit">
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
                            <div className="grid gap-3 sm:grid-cols-2">
                                <label className="space-y-1">
                                    <span className="text-xs font-medium text-gray-600">
                                        {t('views.Settings.Dangerous.purge_data.start_date')}
                                    </span>
                                    <TimePicker
                                        dateOnly
                                        currentLocale={currentLocale}
                                        value={cleanupRanges[cleanupModalTargetId].startDate}
                                        placeholder={t(
                                            'views.Settings.Dangerous.purge_data.start_date'
                                        )}
                                        className="w-full rounded-md border border-gray-300 px-3 py-2 text-center text-sm shadow-sm transition-all hover:ring focus:outline-none disabled:bg-gray-100"
                                        onChange={handleCleanupStartDateChange}
                                    />
                                </label>

                                <label className="space-y-1">
                                    <span className="text-xs font-medium text-gray-600">
                                        {t('views.Settings.Dangerous.purge_data.end_date')}
                                    </span>
                                    <TimePicker
                                        dateOnly
                                        currentLocale={currentLocale}
                                        value={cleanupRanges[cleanupModalTargetId].endDate}
                                        placeholder={t(
                                            'views.Settings.Dangerous.purge_data.end_date'
                                        )}
                                        className="w-full rounded-md border border-gray-300 px-3 py-2 text-center text-sm shadow-sm transition-all hover:ring focus:outline-none disabled:bg-gray-100"
                                        onChange={handleCleanupEndDateChange}
                                    />
                                </label>
                            </div>
                        )}

                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-700">
                                {t('views.Settings.Dangerous.purge_data.verification_label')}
                            </label>
                            <p className="text-sm leading-6 text-gray-500">
                                {t('views.Settings.Dangerous.purge_data.verification_help', {
                                    verificationText: cleanupVerificationText
                                })}
                            </p>
                            <input
                                className="input w-full border border-gray-300 shadow-sm transition-all hover:ring focus:outline-none"
                                value={cleanupVerificationInput}
                                placeholder={cleanupVerificationText}
                                onChange={handleCleanupVerificationInputChange}
                            />
                        </div>

                        <div className="flex flex-col-reverse justify-end gap-2 pt-2 sm:flex-row">
                            <button
                                className="btn rounded-lg px-4 py-2 font-medium transition-all hover:bg-gray-300"
                                onClick={handleCleanupModalClose}
                            >
                                {cleanupModalTarget.cancelBtnText}
                            </button>
                            <button
                                className="btn rounded-lg bg-red-500 px-4 py-2 font-medium text-white transition-all hover:bg-red-700 disabled:cursor-not-allowed disabled:bg-gray-400"
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
