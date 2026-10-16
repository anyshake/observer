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
        <div className="lg:col-span-4">
            <Card
                className="flex h-56 flex-shrink-0 flex-col space-y-2"
                title={t('views.History.search_events.title')}
                iconPath={mdiCalendarSearch}
            >
                <label className="font-medium whitespace-nowrap">
                    {t('views.History.search_events.data_agency')}
                </label>

                <div className="flex-grow space-y-1 overflow-y-scroll rounded-md border border-gray-300 p-2">
                    {getSeismicEventSourceListLoading ? (
                        <div className="flex h-full items-center justify-center">
                            <span className="loading loading-dots loading-md text-gray-400" />
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
                                        className={`hover:bg-base-300 flex cursor-pointer items-center space-x-2 rounded-md px-3 py-2 text-sm transition-all ${
                                            item?.id === selectedSeismicEventSource?.id
                                                ? 'bg-base-300 font-medium text-gray-700'
                                                : ''
                                        }`}
                                    >
                                        <img
                                            src={countryFlagImg}
                                            alt={item?.country}
                                            className="size-4 flex-shrink-0 rounded-lg"
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
                    className="btn btn-sm mt-4 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                >
                    {t('views.History.search_events.search_button')}
                </button>
            </Card>
        </div>
    );
};
