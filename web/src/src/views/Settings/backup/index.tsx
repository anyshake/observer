import { useCallback, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { CodeBlock } from '../../../components/widget/CodeBlock';
import { useExportGlobalConfigLazyQuery, useImportGlobalConfigMutation } from '../../../graphql';
import { sendPromiseAlert } from '../../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../../helpers/alert/sendUserConfirm';

export const Backup = () => {
    const { t } = useTranslation();
    const [mode, setMode] = useState<'backup' | 'restore'>('backup');

    const [, { refetch: exportGlobalConfigRefetch, data: exportGlobalConfigData }] =
        useExportGlobalConfigLazyQuery();
    const handleExportGlobalConfig = useCallback(
        async () =>
            await sendPromiseAlert(
                exportGlobalConfigRefetch(),
                t('views.Settings.Backup.backup_config.generating'),
                t('views.Settings.Backup.backup_config.success'),
                (error) => t('views.Settings.Backup.backup_config.error', { error })
            ),
        [exportGlobalConfigRefetch, t]
    );

    const fileInputRef = useRef<HTMLInputElement>(null);
    const [importGlobalConfig] = useImportGlobalConfigMutation();
    const handleSelectBackupMode = useCallback(() => {
        setMode('backup');
    }, []);

    const handleSelectRestoreMode = useCallback(() => {
        setMode('restore');
    }, []);

    const handleConfirmedImportGlobalConfig = useCallback(async () => {
        const { current: fileInput } = fileInputRef;
        if (!fileInput) {
            return;
        }
        const file = fileInput.files?.[0];
        if (!file) {
            sendUserAlert(t('views.Settings.Backup.restore_config.no_file_selected'), true);
            return;
        }
        await sendPromiseAlert(
            importGlobalConfig({ variables: { data: await file.text() } }),
            t('views.Settings.Backup.restore_config.restoring'),
            t('views.Settings.Backup.restore_config.success'),
            (error) => t('views.Settings.Backup.restore_config.error', { error })
        );
    }, [importGlobalConfig, t]);

    const handleImportGlobalConfig = useCallback(() => {
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
            onConfirmed: handleConfirmedImportGlobalConfig
        });
    }, [handleConfirmedImportGlobalConfig, t]);

    return (
        <div className="app-settings-content">
            <div className="app-options">
                <div className="app-option">
                    <input
                        id="mode-radio-backup"
                        type="radio"
                        checked={mode === 'backup'}
                        onChange={handleSelectBackupMode}
                        className="radio radio-primary radio-xs"
                    />
                    <label htmlFor="mode-radio-backup" className="label app-pointer">
                        {t('views.Settings.Backup.select_mode.backup')}
                    </label>
                </div>
                <div className="app-option">
                    <input
                        id="mode-radio-restore"
                        type="radio"
                        checked={mode === 'restore'}
                        onChange={handleSelectRestoreMode}
                        className="radio radio-primary radio-xs"
                    />
                    <label htmlFor="mode-radio-restore" className="label app-pointer">
                        {t('views.Settings.Backup.select_mode.restore')}
                    </label>
                </div>
            </div>

            {mode === 'backup' && (
                <div className="app-stack">
                    <div className="app-settings-action">
                        <div className="app-stack-sm">
                            <h3 className="app-subheading">
                                {t('views.Settings.Backup.backup_config.title')}
                            </h3>
                            <p className="app-description">
                                {t('views.Settings.Backup.backup_config.description')}
                            </p>
                        </div>
                        <button
                            className="btn app-button-primary app-button-responsive"
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
                <div className="app-backup-import">
                    <div className="app-stack-sm">
                        <h3 className="app-subheading">
                            {t('views.Settings.Backup.restore_config.title')}
                        </h3>
                        <p className="app-description">
                            {t('views.Settings.Backup.restore_config.description')}
                        </p>
                    </div>
                    <div className="app-backup-import-controls">
                        <input
                            ref={fileInputRef}
                            type="file"
                            accept=".json"
                            className="file-input app-file-input"
                        />
                        <button
                            className="btn app-button-primary app-button-auto"
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
