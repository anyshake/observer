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
        <div className="lg:col-span-3">
            <Card
                className="flex h-56 flex-shrink-0 flex-col space-y-2"
                title={t('views.History.data_export.title')}
                iconPath={mdiViewDashboard}
            >
                <div className="flex h-42 w-full flex-row space-x-2">
                    <div className="flex w-1/2 flex-col space-y-2">
                        <label className="font-medium whitespace-nowrap">
                            {t('views.History.data_export.file_format')}
                        </label>
                        <div className="flex-grow space-y-1 overflow-y-auto rounded-md border border-gray-300 p-2">
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
                                            className={`hover:bg-base-300 flex cursor-pointer items-center space-x-2 rounded-md px-3 py-2 text-sm transition-all ${
                                                formatKey === selectedExport?.format
                                                    ? 'bg-base-300 font-medium text-gray-700'
                                                    : ''
                                            }`}
                                        >
                                            <span>{availableExports.formats[formatKey]}</span>
                                        </div>
                                    ))
                            ) : (
                                <div className="flex h-full items-center justify-center">
                                    <span className="loading loading-dots loading-md text-gray-400" />
                                </div>
                            )}
                        </div>
                    </div>
                    <div className="flex w-1/2 flex-col space-y-2">
                        <label className="font-medium whitespace-nowrap">
                            {t('views.History.data_export.channel_code')}
                        </label>
                        <div className="flex-grow space-y-1 overflow-y-scroll rounded-md border border-gray-300 p-2">
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
                                            className={`hover:bg-base-300 flex cursor-pointer items-center space-x-2 rounded-md px-3 py-2 text-sm transition-all ${
                                                channel === selectedExport?.channel
                                                    ? 'bg-base-300 font-medium text-gray-700'
                                                    : ''
                                            }`}
                                        >
                                            <span>{channel}</span>
                                        </div>
                                    ))
                            ) : (
                                <div className="flex h-full items-center justify-center">
                                    <span className="loading loading-dots loading-md text-gray-400" />
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
                    className="btn btn-sm mt-4 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                >
                    {t('views.History.data_export.export_button')}
                </button>
            </Card>
        </div>
    );
};
