import { mdiViewDashboard } from '@mdi/js';
import { useTranslation } from 'react-i18next';

import { Card } from '../../components/widget/Card';
import type { useDataExport } from './useDataExport';

export const DataExport = ({
    isBusy,
    dataExport
}: {
    isBusy: boolean;
    dataExport: ReturnType<typeof useDataExport>;
}) => {
    const { t } = useTranslation();
    const { availableExports, selectedExport, setSelectedExport, handleGetData } = dataExport;

    return (
        <div className="app-history-narrow">
            <Card
                className="app-history-search"
                title={t('views.History.data_export.title')}
                iconPath={mdiViewDashboard}
            >
                <div className="app-export-fields">
                    <div className="app-export-field">
                        <label className="app-field-label">
                            {t('views.History.data_export.file_format')}
                        </label>
                        <div className="app-search-options-auto">
                            {availableExports ? (
                                Object.keys(availableExports.formats)
                                    .sort((a, b) =>
                                        availableExports.formats[a].localeCompare(
                                            availableExports.formats[b]
                                        )
                                    )
                                    .map((formatKey, index) => (
                                        <div
                                            key={`${index}-${formatKey}`}
                                            onClick={() => {
                                                if (!isBusy) {
                                                    setSelectedExport({
                                                        ...selectedExport,
                                                        format: formatKey
                                                    });
                                                }
                                            }}
                                            className={`app-search-option ${
                                                formatKey === selectedExport?.format
                                                    ? 'app-search-option-selected'
                                                    : ''
                                            }`}
                                        >
                                            <span>{availableExports.formats[formatKey]}</span>
                                        </div>
                                    ))
                            ) : (
                                <div className="app-center-fill">
                                    <span className="loading loading-dots loading-md app-loading-muted" />
                                </div>
                            )}
                        </div>
                    </div>
                    <div className="app-export-field">
                        <label className="app-field-label">
                            {t('views.History.data_export.channel_code')}
                        </label>
                        <div className="app-search-options">
                            {availableExports && availableExports.channels ? (
                                availableExports.channels
                                    .sort((a, b) => a.localeCompare(b))
                                    .map((channel, index) => (
                                        <div
                                            key={`${index}-${channel}`}
                                            onClick={() => {
                                                if (!isBusy) {
                                                    setSelectedExport({
                                                        ...selectedExport,
                                                        channel
                                                    });
                                                }
                                            }}
                                            className={`app-search-option ${
                                                channel === selectedExport?.channel
                                                    ? 'app-search-option-selected'
                                                    : ''
                                            }`}
                                        >
                                            <span>{channel}</span>
                                        </div>
                                    ))
                            ) : (
                                <div className="app-center-fill">
                                    <span className="loading loading-dots loading-md app-loading-muted" />
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                <button
                    disabled={
                        isBusy || !selectedExport.channel.length || !selectedExport.format.length
                    }
                    onClick={handleGetData}
                    className="btn btn-sm app-button-primary app-button-search"
                >
                    {t('views.History.data_export.export_button')}
                </button>
            </Card>
        </div>
    );
};
