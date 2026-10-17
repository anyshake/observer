import { mdiCalendarSearch } from '@mdi/js';
import { Buffer } from 'buffer';
import * as _countryFlags from 'country-flag-icons/string/3x2';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { Card } from '../../components/widget/Card';
import type { GetSeismicEventSourceListQuery } from '../../graphql';
import type { useSeismicEvents } from './useSeismicEvents';

export const EventSearch = ({
    currentLocale,
    isBusy,
    getSeismicEventSourceListData,
    getSeismicEventSourceListLoading,
    events
}: {
    currentLocale: string;
    isBusy: boolean;
    getSeismicEventSourceListData: GetSeismicEventSourceListQuery | undefined;
    getSeismicEventSourceListLoading: boolean;
    events: ReturnType<typeof useSeismicEvents>;
}) => {
    const { t } = useTranslation();
    const countryFlags = useMemo(
        () =>
            Object.fromEntries(
                Object.entries(_countryFlags).map(([key, value]) => [
                    key,
                    `data:image/svg+xml;base64,${Buffer.from(value).toString('base64')}`
                ])
            ),
        []
    );
    const { selectedSeismicEventSource, setSelectedSeismicEventSource, handleSearchEvents } =
        events;

    return (
        <div className="app-history-wide">
            <Card
                className="app-history-search"
                title={t('views.History.search_events.title')}
                iconPath={mdiCalendarSearch}
            >
                <label className="app-field-label">
                    {t('views.History.search_events.data_agency')}
                </label>

                <div className="app-search-options">
                    {getSeismicEventSourceListLoading ? (
                        <div className="app-center-fill">
                            <span className="loading loading-dots loading-md app-loading-muted" />
                        </div>
                    ) : (
                        (getSeismicEventSourceListData?.getEventSource
                            ? [...getSeismicEventSourceListData.getEventSource]
                            : []
                        )
                            .sort((a, b) => {
                                const nameA =
                                    a?.locales[currentLocale] ??
                                    a?.locales[a!.defaultLocale] ??
                                    a?.id ??
                                    '';
                                const nameB =
                                    b?.locales[currentLocale] ??
                                    b?.locales[b!.defaultLocale] ??
                                    b?.id ??
                                    '';
                                return nameA.localeCompare(nameB);
                            })
                            .map((item, index) => {
                                const countryFlagImg = countryFlags[item!.country];
                                return (
                                    <div
                                        key={`${index}-${item?.id}`}
                                        onClick={() => {
                                            if (!isBusy) {
                                                setSelectedSeismicEventSource({
                                                    id: item!.id,
                                                    locales: item!.locales,
                                                    countryFlag: countryFlagImg,
                                                    defaultLocale: item!.defaultLocale
                                                });
                                            }
                                        }}
                                        className={`app-search-option ${
                                            item?.id === selectedSeismicEventSource?.id
                                                ? 'app-search-option-selected'
                                                : ''
                                        }`}
                                    >
                                        <img
                                            src={countryFlagImg}
                                            alt={item?.country}
                                            className="app-event-source-flag"
                                        />
                                        <span>
                                            {item?.locales[currentLocale] ??
                                                item?.locales[item.defaultLocale]}
                                        </span>
                                    </div>
                                );
                            })
                    )}
                </div>

                <button
                    disabled={
                        isBusy || !selectedSeismicEventSource || getSeismicEventSourceListLoading
                    }
                    onClick={handleSearchEvents}
                    className="btn btn-sm app-button-primary app-button-search"
                >
                    {t('views.History.search_events.search_button')}
                </button>
            </Card>
        </div>
    );
};
