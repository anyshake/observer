import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { CodeBlock } from '../../components/CodeBlock';
import { useExportGlobalConfigLazyQuery, useImportGlobalConfigMutation } from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';

export const Backup = () => {
    const { t } = useTranslation();
    const [mode, setMode] = useState<'backup' | 'restore'>('backup');

    const [, { refetch: exportGlobalConfigRefetch, data: exportGlobalConfigData }] =
        useExportGlobalConfigLazyQuery();
    const handleExportGlobalConfig = async () =>
        await sendPromiseAlert(
            exportGlobalConfigRefetch(),
            t('views.Settings.Backup.backup_config.generating'),
            t('views.Settings.Backup.backup_config.success'),
            (error) => t('views.Settings.Backup.backup_config.error', { error })
        );

    const fileInputRef = useRef<HTMLInputElement>(null);
    const [importGlobalConfig] = useImportGlobalConfigMutation();
    const handleImportGlobalConfig = () => {
        const { current: fileInput } = fileInputRef;
        if (!fileInput) {
            return;
        }
        const file = fileInput.files?.[0];
        if (!file) {
            sendUserAlert(t('views.Settings.Backup.restore_config.no_file_selected'), true);
            return;
        }
        sendUserConfirm(t('views.Settings.Backup.restore_config.confirm_message'), {
            title: t('views.Settings.Backup.restore_config.confirm_title'),
            cancelBtnText: t('views.Settings.Backup.restore_config.cancel_button'),
            confirmBtnText: t('views.Settings.Backup.restore_config.confirm_button'),
            onConfirmed: async () => {
                await sendPromiseAlert(
                    importGlobalConfig({ variables: { data: await file.text() } }),
                    t('views.Settings.Backup.restore_config.restoring'),
                    t('views.Settings.Backup.restore_config.success'),
                    (error) => t('views.Settings.Backup.restore_config.error', { error })
                );
            }
        });
    };

    return (
        <div className="mx-auto max-w-3xl space-y-4">
            <div className="flex flex-wrap gap-4">
                <div className="flex items-center gap-2">
                    <input
                        id="mode-radio-backup"
                        type="radio"
                        checked={mode === 'backup'}
                        onChange={() => setMode('backup')}
                        className="radio radio-primary radio-xs"
                    />
                    <label htmlFor="mode-radio-backup" className="label cursor-pointer">
                        {t('views.Settings.Backup.select_mode.backup')}
                    </label>
                </div>
                <div className="flex items-center gap-2">
                    <input
                        id="mode-radio-restore"
                        type="radio"
                        checked={mode === 'restore'}
                        onChange={() => setMode('restore')}
                        className="radio radio-primary radio-xs"
                    />
                    <label htmlFor="mode-radio-restore" className="label cursor-pointer">
                        {t('views.Settings.Backup.select_mode.restore')}
                    </label>
                </div>
            </div>

            {mode === 'backup' && (
                <div className="space-y-4">
                    <div className="flex flex-col justify-between gap-4 p-4 sm:flex-row sm:items-center">
                        <div className="space-y-2">
                            <h3 className="text-lg font-semibold text-gray-800">
                                {t('views.Settings.Backup.backup_config.title')}
                            </h3>
                            <p className="text-sm text-gray-500">
                                {t('views.Settings.Backup.backup_config.description')}
                            </p>
                        </div>
                        <button
                            className="btn w-full rounded-lg bg-purple-500 px-4 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700 sm:w-auto"
                            onClick={handleExportGlobalConfig}
                        >
                            {t('views.Settings.Backup.backup_config.submit_button')}
                        </button>
                    </div>
                    {exportGlobalConfigData?.exportGlobalConfig && (
                        <CodeBlock
                            fileName={`backup-${Date.now()}.json`}
                            language="json"
                            showLineNumbers
                        >
                            {JSON.stringify(
                                JSON.parse(exportGlobalConfigData.exportGlobalConfig),
                                null,
                                2
                            )}
                        </CodeBlock>
                    )}
                </div>
            )}

            {mode === 'restore' && (
                <div className="flex flex-col gap-4 p-4">
                    <div className="space-y-2">
                        <h3 className="text-lg font-semibold text-gray-800">
                            {t('views.Settings.Backup.restore_config.title')}
                        </h3>
                        <p className="text-sm text-gray-500">
                            {t('views.Settings.Backup.restore_config.description')}
                        </p>
                    </div>
                    <div className="flex flex-col gap-4 sm:flex-row sm:items-center md:justify-between">
                        <input
                            ref={fileInputRef}
                            type="file"
                            accept=".json"
                            className="file-input w-auto border-gray-300 transition-all hover:border-gray-400 hover:shadow-md focus:outline-none md:w-md"
                        />
                        <button
                            className="btn w-auto rounded-lg bg-purple-500 px-4 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                            onClick={handleImportGlobalConfig}
                        >
                            {t('views.Settings.Backup.restore_config.submit_button')}
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};
