import { mdiImageOutline } from '@mdi/js';
import Icon from '@mdi/react';
import { RefObject, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { exportWorkspace } from '../../helpers/utils/exportWorkspace';

interface IWorkspaceExportButton {
    readonly workspaceRef: RefObject<HTMLDivElement | null>;
    readonly fileName: string;
    readonly disabled?: boolean;
}

export const WorkspaceExportButton = ({
    workspaceRef,
    fileName,
    disabled
}: IWorkspaceExportButton) => {
    const { t } = useTranslation();
    const [exporting, setExporting] = useState(false);

    const handleExport = useCallback(async () => {
        if (!workspaceRef.current) {
            return;
        }
        setExporting(true);
        try {
            await exportWorkspace(workspaceRef.current, '.app-draggable', fileName);
        } catch (error) {
            sendUserAlert(
                t('components.WorkspaceExportButton.export_failed', {
                    message: (error as Error).message
                }),
                true
            );
        } finally {
            setExporting(false);
        }
    }, [fileName, t, workspaceRef]);

    return (
        <button
            className="btn btn-sm app-button-toolbar"
            disabled={disabled || exporting}
            aria-busy={exporting}
            onClick={handleExport}
        >
            {exporting ? (
                <span className="loading loading-spinner loading-xs" />
            ) : (
                <Icon className="app-icon" path={mdiImageOutline} size={0.7} />
            )}
            <span>{t('components.WorkspaceExportButton.export_button')}</span>
        </button>
    );
};
