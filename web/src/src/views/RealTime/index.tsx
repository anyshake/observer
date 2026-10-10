import {
    mdiBlurLinear,
    mdiClose,
    mdiLock,
    mdiLockOpen,
    mdiLockReset,
    mdiRecordCircle,
    mdiRecordRec,
    mdiWaveform
} from '@mdi/js';
import Icon from '@mdi/react';
import { createRef, RefObject, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import type { ColorMapName } from 'spectrogram-js';
import { FFTExecutor } from 'spectrogram-js';

import { DequeChartHandle } from '../../components/chart/DequeChart';
import { DequeSpectrogramHandle } from '../../components/chart/DequeSpectrogram';
import { DEFAULT_SPECTROGRAM_COLOR_MAP } from '../../components/chart/spectrogramColorMaps';
import { Banner } from '../../components/widget/Banner';
import { RealTimeConstraints } from '../../config/constraints';
import { routerConfig } from '../../config/router';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { getSocketApiUrl } from '../../helpers/app/getSocketApiUrl';
import { useSocket } from '../../helpers/request/useSocket';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { sameStrings } from '../../helpers/utils/sameStrings';
import { useThrottleFnTrailing } from '../../helpers/utils/useThrottleFnTrailing';
import { useCredentialStore } from '../../stores/credential';
import { type LayoutConfig, useLayoutStore } from '../../stores/layout';
import { useRetentionStore } from '../../stores/retention';
import { RealtimeChannel } from './RealtimeChannel';

const RealTime = () => {
    const { t } = useTranslation();

    const { retention } = useRetentionStore();
    const { getCredential } = useCredentialStore();
    const { config, locks, toggleLock, setLayoutConfig, resetLayoutConfig } = useLayoutStore();

    const [sampleRate, setSampleRate] = useState(0);
    const [updatedAt, setUpdatedAt] = useState(0);
    const [displayMode, setDisplayMode] = useState<'waveform' | 'spectrogram'>('waveform');
    useEffect(() => {
        // Make sure sample rate supports spectrogram display
        if (sampleRate < RealTimeConstraints.freqRange[1] * 2) {
            setDisplayMode('waveform');
        }
    }, [sampleRate]);

    const [recordingState, setRecordingState] = useState({
        isRecording: false,
        startTime: 0,
        endTime: 0
    });
    const [recordList, setRecordList] = useState<[number, number, string][]>([]);
    const [activeChannels, setActiveChannels] = useState<
        Record<string, { id: string; index: number }>
    >({});

    const sharedFFTExecutor = useMemo(() => new FFTExecutor(RealTimeConstraints.fftSize), []);
    const waveformRefs = useRef<{ [key: string]: RefObject<DequeChartHandle> }>({});
    const spectrogramRefs = useRef<{ [key: string]: RefObject<DequeSpectrogramHandle> }>({});
    const prevChannelsRef = useRef<string[]>([]);
    const prevChannelIndicesRef = useRef<number[]>([]);
    const layoutCacheRef = useRef(new Map<string, LayoutConfig>());
    const [activeChart, setActiveChart] = useState<string | null>(null); // Track the active chart

    const updateChannels = useCallback((channelData: Record<string, { channel_id: number }>) => {
        const currentChannels = Object.keys(channelData);
        // Backend channel_id starts at 1; object key order does not reflect physical channel order.
        const channelIndices = currentChannels.map((channel) => channelData[channel].channel_id - 1);

        if (
            sameStrings(currentChannels, prevChannelsRef.current) &&
            channelIndices.every((index, position) => index === prevChannelIndicesRef.current[position])
        ) {
            return;
        }
        setActiveChannels((prevChannels) => {
            const newChannels = { ...prevChannels };
            currentChannels.forEach((channel, position) => {
                if (!newChannels[channel]) {
                    waveformRefs.current[channel] =
                        createRef<DequeChartHandle>() as RefObject<DequeChartHandle>;
                    spectrogramRefs.current[channel] =
                        createRef<DequeSpectrogramHandle>() as RefObject<DequeSpectrogramHandle>;
                }
                newChannels[channel] = {
                    id: `${RealTimeConstraints.id}_${channel}`,
                    index: channelIndices[position]
                };
            });
            Object.keys(newChannels).forEach((channel) => {
                if (!currentChannels.includes(channel)) {
                    delete newChannels[channel];
                    delete waveformRefs.current[channel];
                    delete spectrogramRefs.current[channel];
                }
            });
            prevChannelsRef.current = currentChannels;
            prevChannelIndicesRef.current = channelIndices;
            return newChannels;
        });
    }, []);

    const updatedAtRef = useRef(0);
    const { readyState, sendMessage } = useSocket(
        {
            url: getSocketApiUrl(),
            onData: ({ data }) => {
                const { channel_data, sample_rate, record_time, current_time } = data;

                for (const channel in channel_data) {
                    const samples = channel_data[channel].data_array;
                    waveformRefs.current[channel]?.current?.addData(
                        samples,
                        record_time,
                        current_time,
                        sample_rate
                    );
                    spectrogramRefs.current[channel]?.current?.addData(
                        samples,
                        record_time,
                        current_time,
                        sample_rate
                    );
                }

                updatedAtRef.current = record_time;
                setSampleRate((prev) => (prev === sample_rate ? prev : sample_rate));
                setUpdatedAt((prev) =>
                    Math.floor(prev / 1000) === Math.floor(record_time / 1000) ? prev : record_time
                );

                updateChannels(channel_data);
            }
        },
        true
    );
    useEffect(() => {
        if (readyState === 1) {
            sendMessage(getCredential().token);
            sendMessage('client hello');
        }
    }, [readyState, sendMessage, getCredential]);

    const getInitialLayout = useCallback(
        (id: string, index: number) => {
            const saved = config[id];
            if (saved?.position && saved?.size && saved?.spectrogram) {
                return saved;
            }

            const wide = document.documentElement.clientWidth > 768;
            const cacheKey = wide ? `${id}:w` : `${id}:n`;
            const cached = layoutCacheRef.current.get(cacheKey);
            if (cached) {
                return cached;
            }

            let x = 20;
            let y = 50;
            for (let i = 0; i < id.length; i++) {
                x = (x + id.charCodeAt(i) * (i + 1)) % 450;
                y = (y + id.charCodeAt(i) * (i + 2)) % 450;
            }
            x = Math.max(20, Math.min(x, 500));
            y = Math.max(50, Math.min(y, 500));

            const layout: LayoutConfig = {
                position: { x, y },
                size: {
                    width: wide ? RealTimeConstraints.minWidth * 2 : RealTimeConstraints.minWidth,
                    height: wide ? RealTimeConstraints.minWidth * 2 : RealTimeConstraints.minHeight
                },
                spectrogram: {
                    ...RealTimeConstraints.getDynamicDB(index),
                    colorMap: DEFAULT_SPECTROGRAM_COLOR_MAP
                }
            };
            layoutCacheRef.current.set(cacheKey, layout);
            return layout;
        },
        [config]
    );

    const handleResetLayout = useCallback(() => {
        sendUserConfirm(t('views.RealTime.reset_layout.confirm_message'), {
            title: t('views.RealTime.reset_layout.confirm_title'),
            cancelBtnText: t('views.RealTime.reset_layout.cancel_button'),
            confirmBtnText: t('views.RealTime.reset_layout.confirm_button'),
            onConfirmed: () => {
                Object.values(activeChannels).forEach(({ id }) => resetLayoutConfig(id));
            }
        });
    }, [t, activeChannels, resetLayoutConfig]);

    const handleActiveChart = useCallback((channel: string) => {
        setActiveChart(channel);
    }, []);

    const handleToggleDisplayMode = useCallback(() => {
        setDisplayMode((prevMode) => (prevMode === 'waveform' ? 'spectrogram' : 'waveform'));
    }, []);

    const handleToggleRecording = useCallback(() => {
        const stamp = updatedAtRef.current;
        setRecordingState((prev) => {
            const { isRecording, startTime } = prev;
            if (!isRecording) {
                if (stamp === 0) {
                    return prev;
                }
                sendUserAlert(
                    t('views.RealTime.record_data.start_recording', {
                        startedAt: getTimeString(stamp)
                    })
                );
                return { ...prev, isRecording: true, startTime: stamp };
            }

            const endTime = stamp;
            if (startTime !== endTime) {
                const search = new URLSearchParams({
                    start_time: startTime.toString(),
                    end_time: endTime.toString()
                });
                setRecordList((prevList) => [...prevList, [startTime, endTime, search.toString()]]);
                sendUserAlert(
                    t('views.RealTime.record_data.link_created', {
                        endedAt: getTimeString(endTime),
                        duration: ((endTime - startTime) / 1000).toFixed(1)
                    })
                );
            } else {
                sendUserAlert(t('views.RealTime.record_data.no_data_recorded'), true);
            }

            return { ...prev, isRecording: false, endTime };
        });
    }, [t]);

    const handleRemoveRecord = useCallback(
        (index: number) => {
            setRecordList((prevList) => prevList.filter((_, i) => i !== index));
        },
        [setRecordList]
    );

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
            (
                channel: string,
                index: number,
                minDB: number,
                maxDB: number,
                colorMap: ColorMapName
            ) => {
                setLayoutConfig(channel, {
                    ...getInitialLayout(channel, index),
                    spectrogram: { maxDB, minDB, colorMap }
                });
            },
            [getInitialLayout, setLayoutConfig]
        ),
        500
    );

    return (
        <div className="container mx-auto space-y-6 p-4">
            <Banner
                status={readyState === 1 ? 'ok' : 'warning'}
                message={
                    readyState === 1
                        ? t('views.RealTime.connectivity.connected', {
                              updatedAt: getTimeString(updatedAt)
                          })
                        : t('views.RealTime.connectivity.connecting')
                }
            />
            <div className="flex flex-wrap items-center gap-2">
                <button
                    className="btn btn-sm flex items-center"
                    onClick={() => toggleLock(RealTimeConstraints.id)}
                >
                    {locks[RealTimeConstraints.id] ? (
                        <>
                            <Icon className="flex-shrink-0" path={mdiLock} size={0.7} />
                            <span>{t('views.RealTime.layout_locker.unlock_button')}</span>
                        </>
                    ) : (
                        <>
                            <Icon className="flex-shrink-0" path={mdiLockOpen} size={0.7} />
                            <span>{t('views.RealTime.layout_locker.lock_button')}</span>
                        </>
                    )}
                </button>
                <button className="btn btn-sm flex items-center" onClick={handleResetLayout}>
                    <Icon className="flex-shrink-0" path={mdiLockReset} size={0.7} />
                    <span>{t('views.RealTime.reset_layout.reset_button')}</span>
                </button>
                {sampleRate >= RealTimeConstraints.freqRange[1] * 2 && (
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
                                ? t('views.RealTime.display_mode.spectrogram_mode')
                                : t('views.RealTime.display_mode.waveform_mode')}
                        </span>
                    </button>
                )}
                <button
                    className={`btn btn-sm flex items-center ${recordingState.isRecording ? 'text-pink-800' : ''}`}
                    onClick={handleToggleRecording}
                >
                    <Icon
                        className="flex-shrink-0"
                        path={recordingState.isRecording ? mdiRecordRec : mdiRecordCircle}
                        size={0.7}
                    />
                    <span className={recordingState.isRecording ? 'animate-pulse' : ''}>
                        {recordingState.isRecording
                            ? t('views.RealTime.record_data.stop_button')
                            : t('views.RealTime.record_data.start_button')}
                    </span>
                </button>
                <div className="flex flex-wrap gap-2">
                    <div className="badge badge-soft badge-primary font-medium">
                        {t('views.RealTime.stream_status.sample_rate', { value: sampleRate })}
                    </div>
                    <div className="badge badge-soft badge-secondary font-medium">
                        {t('views.RealTime.stream_status.channels', {
                            num: Object.keys(activeChannels).length
                        })}
                    </div>
                    <div className="badge badge-soft badge-info font-medium">
                        {t('views.RealTime.stream_status.retention', { retention })}
                    </div>
                </div>
            </div>

            {recordList.length > 0 && (
                <div className="flex w-fit flex-col rounded-lg border border-dashed border-gray-300 p-4">
                    {recordList.map(([startTime, endTime, search], index) => (
                        <li className="flex flex-row items-center space-x-2">
                            <Link
                                key={index}
                                className="link link-primary text-sm font-medium"
                                target="_blank"
                                to={{ search, pathname: routerConfig.routes.history.uri }}
                            >
                                {`${getTimeString(startTime)} - ${getTimeString(endTime)}`}
                            </Link>
                            <button
                                className="cursor-pointer text-gray-500 hover:text-gray-700 disabled:cursor-not-allowed disabled:text-gray-300"
                                onClick={() => handleRemoveRecord(index)}
                            >
                                <Icon className="flex-shrink-0" path={mdiClose} size={0.7} />
                            </button>
                        </li>
                    ))}
                </div>
            )}

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
                        const channelInfo = activeChannels[channel];
                        return (
                            <RealtimeChannel
                                key={channel}
                                channel={channel}
                                channelId={channelInfo.id}
                                index={channelInfo.index}
                                layout={getInitialLayout(channelInfo.id, channelInfo.index)}
                                locked={Boolean(locks[RealTimeConstraints.id])}
                                waveformMode={displayMode === 'waveform'}
                                retention={retention}
                                sampleRate={sampleRate}
                                fftExecutor={sharedFFTExecutor}
                                waveformRef={waveformRefs.current[channel]}
                                spectrogramRef={spectrogramRefs.current[channel]}
                                onActive={handleActiveChart}
                                onDragStop={handleDragStop}
                                onResizeStop={handleResizeStop}
                                onSpectrogram={handleSpectrogramUpdate}
                            />
                        );
                    })}
            </div>
        </div>
    );
};

export default RealTime;
