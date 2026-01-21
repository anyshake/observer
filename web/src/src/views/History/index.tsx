import {
    mdiBlurLinear,
    mdiCalendarSearch,
    mdiClock,
    mdiLink,
    mdiLock,
    mdiLockOpen,
    mdiLockReset,
    mdiTarget,
    mdiViewDashboard,
    mdiWaveform
} from '@mdi/js';
import Icon from '@mdi/react';
import { Buffer } from 'buffer';
import * as _countryFlags from 'country-flag-icons/string/3x2';
import { ReactNode, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';
import { FFTExecutor } from 'spectrogram-js';

import { LineChart } from '../../components/chart/LineChart';
import { Spectrogram } from '../../components/chart/Spectrogram';
import { DialogModal } from '../../components/ui/DialogModal';
import { DraggableBox } from '../../components/ui/DraggableBox';
import { Card } from '../../components/widget/Card';
import { List } from '../../components/widget/List';
import { TimePicker } from '../../components/widget/TimePicker';
import { HistoryConstraints } from '../../config/constraints';
import { IRouterComponent } from '../../config/router';
import {
    useGetSeismicEventBySourceLazyQuery,
    useGetSeismicEventSourceListQuery,
    useGetSeismicRecordsLazyQuery
} from '../../graphql';
import { sendPromiseAlert } from '../../helpers/alert/sendPromiseAlert';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { getRestfulApiUrl } from '../../helpers/app/getRestfulApiUrl';
import { ApiClient } from '../../helpers/request/ApiClient';
import { useUrlParams } from '../../helpers/request/useUrlParams';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { setClipboardText } from '../../helpers/utils/setClipboardText';
import { useThrottleFnTrailing } from '../../helpers/utils/useThrottleFnTrailing';
import { useLayoutStore } from '../../stores/layout';
import { useRetentionStore } from '../../stores/retention';

const History = ({ currentLocale }: IRouterComponent) => {
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
    const sharedFFTExecutor = useMemo(() => new FFTExecutor(HistoryConstraints.fftSize), []);

    const [displayMode, setDisplayMode] = useState<'waveform' | 'spectrogram'>('waveform');
    const handleToggleDisplayMode = useCallback(() => {
        setDisplayMode((prevMode) => (prevMode === 'waveform' ? 'spectrogram' : 'waveform'));
    }, []);

    const [isBusy, setIsBusy] = useState(true);
    const [searchParams, setSearchParams] = useSearchParams();
    const { start_time: startTimeInSearchParam, end_time: endTimeInSearchParam } = useUrlParams<{
        start_time?: string;
        end_time?: string;
    }>(searchParams);
    const [startTime, setStartTime] = useState(parseInt(startTimeInSearchParam ?? '0'));
    const [endTime, setEndTime] = useState(parseInt(endTimeInSearchParam ?? '0'));
    const { data: getSeismicEventSourceListData, loading: getSeismicEventSourceListLoading } =
        useGetSeismicEventSourceListQuery();
    useEffect(() => {
        if (getSeismicEventSourceListData?.getCurrentTime && !startTime && !endTime) {
            const { getCurrentTime: currentTime } = getSeismicEventSourceListData;
            setStartTime(currentTime - 1000 * 60 * 5);
            setEndTime(currentTime);
            setIsBusy(false);
        } else if (startTime && endTime) {
            setIsBusy(false);
        }
    }, [startTime, endTime, getSeismicEventSourceListData]);

    const [getSeismicEventBySource] = useGetSeismicEventBySourceLazyQuery();
    const [selectedSeismicEventSource, setSelectedSeismicEventSource] = useState<{
        id: string;
        countryFlag: string;
        locales: Record<string, string>;
        defaultLocale: string;
    }>();
    const [isSelectModalOpen, setIsSelectModalOpen] = useState(false);
    const [seismicEventsData, setSeismicEventsData] = useState<
        Array<{
            id: string;
            region: string;
            timestamp: number;
            estimation: number[];
            primary: ReactNode | ReactNode[];
            secondary: ReactNode | ReactNode[];
        }>
    >([]);
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
                    (data ? [...data] : [])
                        .sort((a, b) => b.timestamp - a.timestamp)
                        .map((item) => ({
                            id: item.eventId,
                            region: item.region,
                            timestamp: item.timestamp,
                            estimation: item.estimation,
                            primary: (
                                <div className="space-y-2">
                                    <h3 className="font-semibold text-gray-700">{item.region}</h3>
                                    <ul className="list-inside list-disc whitespace-nowrap text-gray-500">
                                        <li>
                                            {t('views.History.event_list_modal.earthquake_at', {
                                                time: getTimeString(item.timestamp)
                                            })}
                                        </li>
                                        <li>
                                            {t('views.History.event_list_modal.p_wave_arrival', {
                                                time: item.estimation[0].toFixed(1)
                                            })}
                                        </li>
                                        <li>
                                            {t('views.History.event_list_modal.s_wave_arrival', {
                                                time: item.estimation[1].toFixed(1)
                                            })}
                                        </li>
                                        <li>
                                            {t(
                                                'views.History.event_list_modal.epicenter_distance',
                                                {
                                                    distance: item.distance.toFixed(1)
                                                }
                                            )}
                                        </li>
                                        {item.depth !== -1 && (
                                            <li>
                                                {t(
                                                    'views.History.event_list_modal.earthquake_depth',
                                                    {
                                                        depth: item.depth.toFixed(1)
                                                    }
                                                )}
                                            </li>
                                        )}
                                    </ul>
                                </div>
                            ),
                            secondary: (
                                <div className="flex flex-wrap gap-2 font-medium">
                                    {Object.keys(item.magnitude).map((key, index) => (
                                        <div
                                            key={`${index}-${key}`}
                                            className="badge badge-soft badge-secondary whitespace-nowrap"
                                        >
                                            {`${key} ${item.magnitude[key].toFixed(1)}`}
                                        </div>
                                    ))}
                                </div>
                            )
                        })) ?? []
                );
                setIsSelectModalOpen(true);
                return t('views.History.search_events.success', { count: data.length });
            },
            (error) => t('views.History.search_events.error', { error })
        );
    }, [t, getSeismicEventBySource, selectedSeismicEventSource]);
    const handleChooseEvent = useCallback(
        (eventId: string) => {
            const selectedEvent = seismicEventsData.find((item) => item.id === eventId);
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
        [t, seismicEventsData]
    );

    const [availableExports, setAvailableExports] = useState<{
        channels: string[];
        formats: Record<string, string>;
    }>();
    const [selectedExport, setSelectedExport] = useState<{ channel: string; format: string }>({
        channel: '',
        format: ''
    });
    useEffect(() => {
        const requestFn = async () => {
            setIsBusy(true);
            const { data } = await ApiClient.request<{
                channel_code: string[];
                data_format: Record<string, string>;
            }>({
                url: getRestfulApiUrl('/export'),
                method: 'get',
                ignoreErrors: true
            });
            setIsBusy(false);
            if (data) {
                const { channel_code, data_format } = data;
                setAvailableExports({ channels: channel_code, formats: data_format });
            }
        };
        requestFn();
    }, []);
    const handleGetData = useCallback(async () => {
        setIsBusy(true);
        await sendPromiseAlert(
            ApiClient.saveAs({
                url: getRestfulApiUrl('/export'),
                method: 'post',
                data: {
                    start_time: startTime,
                    end_time: endTime,
                    data_format: selectedExport.format,
                    channel_code: selectedExport.channel
                }
            }),
            t('views.History.data_export.exporting', {
                channel: selectedExport.channel,
                format: availableExports?.formats[selectedExport.format]
            }),
            t('views.History.data_export.success', {
                channel: selectedExport.channel,
                format: availableExports?.formats[selectedExport.format]
            }),
            (error) =>
                t('views.History.data_export.error', {
                    channel: selectedExport.channel,
                    format: availableExports?.formats[selectedExport.format],
                    error
                })
        );
        setIsBusy(false);
    }, [
        startTime,
        endTime,
        selectedExport.format,
        selectedExport.channel,
        t,
        availableExports?.formats
    ]);

    const [sampleRate, setSampleRate] = useState(0);
    const [chartData, setChartData] = useState<{ [key: string]: Array<[number, number | null]> }>(
        {}
    );
    const [getSeismicRecords] = useGetSeismicRecordsLazyQuery();
    const handleSearchRecords = useCallback(async () => {
        if (startTime > endTime) {
            sendUserAlert(t('views.History.search_records.invalid_time'), true);
            return;
        }
        const requestFn = async () => {
            setSampleRate(0);
            setChartData({});
            setIsBusy(true);
            const { data, error } = await getSeismicRecords({ variables: { startTime, endTime } });
            setIsBusy(false);
            if (error) {
                throw error;
            }
            if (!data?.getSeisRecordsByTime.length) {
                throw new Error(t('views.History.search_records.no_data'));
            }
            return data;
        };
        await sendPromiseAlert(
            requestFn(),
            t('views.History.search_records.searching'),
            ({ getSeisRecordsByTime }) => {
                const records = [...getSeisRecordsByTime].sort(
                    (a, b) => a!.timestamp - b!.timestamp
                );
                const currentChannels = new Set<string>();
                records.forEach((record) => {
                    record?.channelData.forEach((channel) => {
                        currentChannels.add(channel.channelCode);
                    });
                });

                if (
                    JSON.stringify(Array.from(currentChannels)) !==
                    JSON.stringify(prevChannelsRef.current)
                ) {
                    setActiveChannels((prevChannels) => {
                        const newChannels = { ...prevChannels };
                        Array.from(currentChannels).forEach((channel, index) => {
                            if (!newChannels[channel]) {
                                newChannels[channel] = {
                                    id: `${HistoryConstraints.id}_${channel}`,
                                    index
                                };
                            }
                        });
                        Object.keys(newChannels).forEach((channel) => {
                            if (!currentChannels.has(channel)) {
                                delete newChannels[channel];
                            }
                        });
                        prevChannelsRef.current = Array.from(currentChannels);
                        return newChannels;
                    });
                }

                setChartData((prevChartData) => {
                    const newChartData: typeof prevChartData = {};

                    // Initialize chart data for current channels
                    Array.from(currentChannels).forEach((channel) => {
                        newChartData[channel] = prevChartData[channel] ?? [];
                    });

                    let lastTimestamp: number | null = null;

                    // Fill data with gap checking
                    records.forEach((record) => {
                        const { timestamp, sampleRate, channelData } = record!;
                        const interval = 1000 / sampleRate;

                        // Check gap from previous timestamp
                        if (lastTimestamp !== null && timestamp - lastTimestamp > 2000) {
                            Array.from(currentChannels).forEach((channel) => {
                                newChartData[channel].push([lastTimestamp! + 1, null]);
                            });
                        }

                        // Append current record data
                        channelData.forEach((channel) => {
                            const channelCode = channel.channelCode;
                            const dataArray = newChartData[channelCode];
                            if (dataArray) {
                                for (let i = 0; i < channel.data.length; i++) {
                                    dataArray.push([timestamp + i * interval, channel.data[i]]);
                                }
                            }
                        });

                        lastTimestamp = timestamp;
                    });
                    setSampleRate(records.length > 0 ? records[0]!.sampleRate : 0);

                    return newChartData;
                });

                return t('views.History.search_records.success', { count: records.length });
            },
            (error) => t('views.History.search_records.error', { error })
        );
    }, [t, startTime, endTime, getSeismicRecords]);

    const [activeChannels, setActiveChannels] = useState<
        Record<string, { id: string; index: number }>
    >({});
    const { config, locks, toggleLock, setLayoutConfig, resetLayoutConfig } = useLayoutStore();
    const { retention } = useRetentionStore();
    const prevChannelsRef = useRef<string[]>([]);
    const [activeChart, setActiveChart] = useState<string | null>(null); // Track the active chart

    const getInitialLayout = useCallback(
        (id: string, index: number) => {
            if (config[id]?.position && config[id]?.size) {
                return config[id];
            }

            let x = 20;
            let y = 50;
            for (let i = 0; i < id.length; i++) {
                x = (x + id.charCodeAt(i) * (i + 1)) % 450;
                y = (y + id.charCodeAt(i) * (i + 2)) % 450;
            }
            x = Math.max(20, Math.min(x, 500));
            y = Math.max(50, Math.min(y, 500));

            return {
                position: { x, y },
                size: {
                    width:
                        document.documentElement.clientWidth > 768
                            ? HistoryConstraints.minWidth * 2
                            : HistoryConstraints.minWidth,
                    height:
                        document.documentElement.clientWidth > 768
                            ? HistoryConstraints.minWidth * 2
                            : HistoryConstraints.minHeight
                },
                spectrogram: { ...HistoryConstraints.getDynamicDB(index) }
            };
        },
        [config]
    );

    const handleResetLayout = useCallback(() => {
        sendUserConfirm(t('views.History.reset_layout.confirm_message'), {
            title: t('views.History.reset_layout.confirm_title'),
            cancelBtnText: t('views.History.reset_layout.cancel_button'),
            confirmBtnText: t('views.History.reset_layout.confirm_button'),
            onConfirmed: () => {
                Object.values(activeChannels).forEach(({ id }) => resetLayoutConfig(id));
            }
        });
    }, [t, activeChannels, resetLayoutConfig]);

    const handleShareLink = useCallback(() => {
        const params = new URLSearchParams(searchParams);
        params.set('start_time', startTime.toString());
        params.set('end_time', endTime.toString());
        setSearchParams(params);
        setClipboardText(window.location.href);
        sendUserAlert(t('views.History.share_link.link_copied'));
    }, [t, searchParams, startTime, endTime, setSearchParams]);

    const handleDragStop = useCallback(
        (channel: string, index: number, x: number, y: number) => {
            setLayoutConfig(channel, { ...getInitialLayout(channel, index), position: { x, y } });
        },
        [getInitialLayout, setLayoutConfig]
    );

    const handleResizeStop = useCallback(
        (channel: string, index: number, width: number, height: number) => {
            setLayoutConfig(channel, {
                ...getInitialLayout(channel, index),
                size: { width, height }
            });
        },
        [getInitialLayout, setLayoutConfig]
    );

    const handleSpectrogramUpdate = useThrottleFnTrailing(
        useCallback(
            (channel: string, index: number, minDB: number, maxDB: number) => {
                setLayoutConfig(channel, {
                    ...getInitialLayout(channel, index),
                    spectrogram: { maxDB, minDB }
                });
            },
            [getInitialLayout, setLayoutConfig]
        ),
        500
    );

    return (
        <div className="container mx-auto space-y-6 p-4">
            <div className="grid grid-cols-1 gap-4 lg:grid-cols-10">
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
                                                    <span>
                                                        {availableExports.formats[formatKey]}
                                                    </span>
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
                                isBusy ||
                                !selectedExport.channel.length ||
                                !selectedExport.format.length
                            }
                            onClick={handleGetData}
                            className="btn btn-sm mt-4 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                        >
                            {t('views.History.data_export.export_button')}
                        </button>
                    </Card>
                </div>

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
                                    .sort((a, b) => a?.country.localeCompare(b?.country ?? '') ?? 0)
                                    .sort((a, b) => a?.id.localeCompare(b?.id ?? '') ?? 0)
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
                                isBusy ||
                                !selectedSeismicEventSource ||
                                getSeismicEventSourceListLoading
                            }
                            onClick={handleSearchEvents}
                            className="btn btn-sm mt-4 w-full rounded-lg bg-purple-500 py-2 font-medium text-white shadow-lg transition-all hover:bg-purple-700"
                        >
                            {t('views.History.search_events.search_button')}
                        </button>
                    </Card>
                </div>
            </div>

            <div className="flex flex-wrap items-center gap-2">
                <button
                    className="btn btn-sm flex items-center"
                    onClick={() => toggleLock(HistoryConstraints.id)}
                >
                    {locks[HistoryConstraints.id] ? (
                        <>
                            <Icon className="flex-shrink-0" path={mdiLock} size={0.7} />
                            <span>{t('views.History.layout_locker.unlock_button')}</span>
                        </>
                    ) : (
                        <>
                            <Icon className="flex-shrink-0" path={mdiLockOpen} size={0.7} />
                            <span>{t('views.History.layout_locker.lock_button')}</span>
                        </>
                    )}
                </button>
                <button className="btn btn-sm flex items-center" onClick={handleResetLayout}>
                    <Icon className="flex-shrink-0" path={mdiLockReset} size={0.7} />
                    <span>{t('views.History.reset_layout.reset_button')}</span>
                </button>
                {sampleRate >= HistoryConstraints.freqRange[1] * 2 && (
                    <button
                        className="btn btn-sm flex items-center"
                        onClick={handleToggleDisplayMode}
                    >
                        <Icon
                            className="flex-shrink-0"
                            path={displayMode === 'waveform' ? mdiBlurLinear : mdiWaveform}
                            size={0.7}
                        />
                        <span>
                            {displayMode === 'waveform'
                                ? t('views.History.display_mode.spectrogram_mode')
                                : t('views.History.display_mode.waveform_mode')}
                        </span>
                    </button>
                )}
                <button className="btn btn-sm" onClick={handleShareLink}>
                    <Icon className="flex-shrink-0" path={mdiLink} size={0.7} />
                    <span>{t('views.History.share_link.share_button')}</span>
                </button>
            </div>

            <div className="bg-base-300 relative h-[2000px] w-full overflow-scroll rounded-lg md:h-[1000px] lg:h-[800px] xl:h-screen">
                {Object.keys(activeChannels)
                    .sort((a, b) => {
                        // Move active chart to the end
                        if (a === activeChart) {
                            return 1;
                        }
                        // Move active chart to the end
                        if (b === activeChart) {
                            return -1;
                        }
                        return 0;
                    })
                    .map((channel) => {
                        const initialLayout = getInitialLayout(
                            activeChannels[channel].id,
                            activeChannels[channel].index
                        );
                        return (
                            <DraggableBox
                                key={channel}
                                layout={initialLayout}
                                locked={locks[HistoryConstraints.id]}
                                constraints={HistoryConstraints}
                                onDragStart={() => setActiveChart(channel)}
                                onDragStop={(x, y) =>
                                    handleDragStop(
                                        activeChannels[channel].id,
                                        activeChannels[channel].index,
                                        x,
                                        y
                                    )
                                }
                                onResizeStop={(width, height) =>
                                    handleResizeStop(
                                        activeChannels[channel].id,
                                        activeChannels[channel].index,
                                        width,
                                        height
                                    )
                                }
                            >
                                <div
                                    className={
                                        displayMode === 'waveform'
                                            ? 'block h-full w-full'
                                            : 'hidden'
                                    }
                                >
                                    <LineChart
                                        minSpanValue={HistoryConstraints.minSpanValue}
                                        lineColor={HistoryConstraints.lineColor}
                                        height={'100%'}
                                        yPosition="right"
                                        title={channel}
                                        zoom={true}
                                        animation={false}
                                        data={chartData[channel]}
                                    />
                                </div>
                                <div
                                    className={
                                        displayMode === 'spectrogram'
                                            ? 'block h-full w-full'
                                            : 'hidden'
                                    }
                                >
                                    <Spectrogram
                                        title={channel}
                                        duration={retention}
                                        overlap={HistoryConstraints.overlap}
                                        freqRange={HistoryConstraints.freqRange}
                                        windowSize={HistoryConstraints.windowSize}
                                        maxDB={initialLayout.spectrogram.maxDB}
                                        minDB={initialLayout.spectrogram.minDB}
                                        fftExecutor={sharedFFTExecutor}
                                        data={
                                            chartData[channel]
                                                ? chartData[channel].filter(
                                                      (v): v is [number, number] => v[1] !== null
                                                  )
                                                : []
                                        }
                                        sampleRate={sampleRate}
                                        onSpectrogramUpdate={(minDB, maxDB) =>
                                            handleSpectrogramUpdate(
                                                activeChannels[channel].id,
                                                activeChannels[channel].index,
                                                minDB,
                                                maxDB
                                            )
                                        }
                                    />
                                </div>
                            </DraggableBox>
                        );
                    })}
            </div>

            <DialogModal
                heading={
                    <div className="space-y-4 text-gray-800">
                        <div className="flex items-center space-x-2 text-lg font-extrabold">
                            <Icon path={mdiTarget} size={1.2} className="flex-shrink-0" />
                            <h2>{t('views.History.event_list_modal.title')}</h2>
                        </div>
                        <div className="ml-0.5 flex items-center space-x-2 text-sm">
                            <img
                                src={selectedSeismicEventSource?.countryFlag}
                                className="size-5 flex-shrink-0"
                            />
                            <span>
                                {selectedSeismicEventSource?.locales[currentLocale] ??
                                    selectedSeismicEventSource?.locales[
                                        selectedSeismicEventSource?.defaultLocale
                                    ]}
                            </span>
                        </div>
                    </div>
                }
                open={isSelectModalOpen}
                onClose={() => setIsSelectModalOpen(false)}
            >
                <List
                    className="h-96"
                    onClick={(id) => handleChooseEvent(id)}
                    data={seismicEventsData}
                />
            </DialogModal>
        </div>
    );
};

export default History;
