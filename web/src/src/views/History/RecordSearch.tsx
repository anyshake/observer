import { mdiClock } from '@mdi/js';
import { useTranslation } from 'react-i18next';

import { Card } from '../../components/widget/Card';
import { TimePicker } from '../../components/widget/TimePicker';

export const RecordSearch = ({
    currentLocale,
    startTime,
    endTime,
    setStartTime,
    setEndTime,
    isBusy,
    handleSearchRecords
}: {
    currentLocale: string;
    startTime: number;
    endTime: number;
    setStartTime: (value: number) => void;
    setEndTime: (value: number) => void;
    isBusy: boolean;
    handleSearchRecords: () => Promise<void>;
}) => {
    const { t } = useTranslation();

    return (
        <div className="lg:col-span-3">
            <Card
                className="flex h-56 flex-shrink-0 flex-col space-y-4"
                title={t('views.History.search_records.title')}
                iconPath={mdiClock}
            >
                <label className="font-medium whitespace-nowrap">
                    {t('views.History.search_records.start_time')}
                </label>
                <TimePicker
                    className="w-full rounded-md border border-gray-300 py-1 text-center shadow-sm transition-all hover:ring focus:outline-none disabled:bg-gray-100"
                    currentLocale={currentLocale}
                    defaultValue={startTime}
                    onChange={(v) => setStartTime(v)}
                />
                <label className="font-medium whitespace-nowrap">
                    {t('views.History.search_records.end_time')}
                </label>
                <TimePicker
                    className="w-full rounded-md border border-gray-300 py-1 text-center shadow-sm transition-all hover:ring focus:outline-none disabled:bg-gray-100"
                    currentLocale={currentLocale}
                    defaultValue={endTime}
                    onChange={(v) => setEndTime(v)}
                />
                <button
                    disabled={isBusy}
                    onClick={handleSearchRecords}
                    className="btn btn-sm mt-auto w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                >
                    {t('views.History.search_records.search_button')}
                </button>
            </Card>
        </div>
    );
};
