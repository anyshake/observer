import { type Dispatch, ReactNode, type SetStateAction, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useGetSeismicEventBySourceLazyQuery } from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { getTimeString } from '../../helpers/utils/getTimeString';

export const useSeismicEvents = ({
    setIsBusy,
    setStartTime,
    setEndTime
}: {
    setIsBusy: Dispatch<SetStateAction<boolean>>;
    setStartTime: Dispatch<SetStateAction<number>>;
    setEndTime: Dispatch<SetStateAction<number>>;
}) => {
    const { t } = useTranslation();
    const [getSeismicEventBySource] = useGetSeismicEventBySourceLazyQuery();
    const [selectedSeismicEventSource, setSelectedSeismicEventSource] = useState<{
        id: string;
        countryFlag: string;
        locales: Record<string, string>;
        defaultLocale: string;
    }>();
    const [isSelectModalOpen, setIsSelectModalOpen] = useState(false);
    const [seismicEventsData, setSeismicEventsData] = useState<
        Record<
            string,
            {
                id: string;
                region: string;
                timestamp: number;
                estimation: number[];
                primary: ReactNode | ReactNode[];
                secondary: ReactNode | ReactNode[];
            }
        >
    >({});
    const handleSearchEvents = useCallback(async () => {
        const requestFn = async () => {
            setIsBusy(true);
            const { data, error } = await getSeismicEventBySource({
                variables: { sourceId: selectedSeismicEventSource?.id ?? '' }
            });
            setIsBusy(false);
            if (error) {
                throw error;
            }
            if (!data?.getEventsBySource.length) {
                throw new Error(t('views.History.search_events.no_data'));
            }
            return data?.getEventsBySource;
        };
        await sendPromiseAlert(
            requestFn(),
            t('views.History.search_events.searching'),
            (data) => {
                setSeismicEventsData(
                    (data ? [...data] : []).reduce(
                        (
                            acc,
                            { eventId, region, timestamp, estimation, depth, distance, magnitude }
                        ) => {
                            acc[eventId] = {
                                id: eventId,
                                region,
                                timestamp,
                                estimation,
                                primary: (
                                    <div className="app-stack-sm">
                                        <h3 className="app-event-region">{region}</h3>
                                        <ul className="app-event-details">
                                            <li>
                                                {t('views.History.event_list_modal.earthquake_at', {
                                                    time: getTimeString(timestamp)
                                                })}
                                            </li>
                                            <li>
                                                {t(
                                                    'views.History.event_list_modal.p_wave_arrival',
                                                    {
                                                        time: getTimeString(
                                                            timestamp + estimation[0] * 1000
                                                        ),
                                                        seconds: estimation[0].toFixed(1)
                                                    }
                                                )}
                                            </li>
                                            <li>
                                                {t(
                                                    'views.History.event_list_modal.s_wave_arrival',
                                                    {
                                                        time: getTimeString(
                                                            timestamp + estimation[1] * 1000
                                                        ),
                                                        seconds: estimation[1].toFixed(1)
                                                    }
                                                )}
                                            </li>
                                            <li>
                                                {t(
                                                    'views.History.event_list_modal.epicenter_distance',
                                                    { distance: distance.toFixed(1) }
                                                )}
                                            </li>
                                            {depth !== -1 && (
                                                <li>
                                                    {t(
                                                        'views.History.event_list_modal.earthquake_depth',
                                                        { depth: depth.toFixed(1) }
                                                    )}
                                                </li>
                                            )}
                                        </ul>
                                    </div>
                                ),
                                secondary: (
                                    <div className="app-event-tags">
                                        {Object.keys(magnitude).map((key, index) => (
                                            <div
                                                key={`${index}-${key}`}
                                                className="badge badge-soft badge-secondary app-nowrap"
                                            >
                                                {`${key} ${magnitude[key].toFixed(1)}`}
                                            </div>
                                        ))}
                                    </div>
                                )
                            };
                            return acc;
                        },
                        {} as Record<string, (typeof seismicEventsData)[string]>
                    )
                );
                setIsSelectModalOpen(true);
                return t('views.History.search_events.success', { count: data.length });
            },
            (error) => t('views.History.search_events.error', { error })
        );
    }, [t, getSeismicEventBySource, selectedSeismicEventSource, setIsBusy]);
    const handleChooseEvent = useCallback(
        (eventId: string) => {
            const selectedEvent = seismicEventsData[eventId];
            if (selectedEvent) {
                const { timestamp, estimation, region } = selectedEvent;
                const [pWave, sWave] = estimation;
                setStartTime(Math.floor(timestamp + pWave * 1000 - 1000 * 60 * 5));
                setEndTime(Math.floor(timestamp + sWave * 1000 + 1000 * 60 * 5));
                sendUserAlert(
                    t('views.History.event_list_modal.selected_an_event', {
                        region,
                        time: getTimeString(timestamp)
                    })
                );
                setIsSelectModalOpen(false);
            }
        },
        [t, seismicEventsData, setStartTime, setEndTime]
    );

    return {
        selectedSeismicEventSource,
        setSelectedSeismicEventSource,
        isSelectModalOpen,
        setIsSelectModalOpen,
        seismicEventsData,
        handleSearchEvents,
        handleChooseEvent
    };
};
