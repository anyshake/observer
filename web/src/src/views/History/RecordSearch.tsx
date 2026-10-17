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
        <div className="app-history-narrow">
            <Card
                className="app-record-search"
                title={t('views.History.search_records.title')}
                iconPath={mdiClock}
            >
                <label className="app-field-label">
                    {t('views.History.search_records.start_time')}
                </label>
                <TimePicker
                    className="app-time-input"
                    currentLocale={currentLocale}
                    defaultValue={startTime}
                    onChange={(v) => setStartTime(v)}
                />
                <label className="app-field-label">
                    {t('views.History.search_records.end_time')}
                </label>
                <TimePicker
                    className="app-time-input"
                    currentLocale={currentLocale}
                    defaultValue={endTime}
                    onChange={(v) => setEndTime(v)}
                />
                <button
                    disabled={isBusy}
                    onClick={handleSearchRecords}
                    className="btn btn-sm app-button-primary app-button-search-bottom"
                >
                    {t('views.History.search_records.search_button')}
                </button>
            </Card>
        </div>
    );
};
