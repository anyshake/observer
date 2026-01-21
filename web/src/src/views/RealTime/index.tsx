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
import { FFTExecutor } from 'spectrogram-js';

import { DequeChart, DequeChartHandle } from '../../components/chart/DequeChart';
import { DequeSpectrogram, DequeSpectrogramHandle } from '../../components/chart/DequeSpectrogram';
import { DraggableBox } from '../../components/ui/DraggableBox';
import { Banner } from '../../components/widget/Banner';
import { RealTimeConstraints } from '../../config/constraints';
import { routerConfig } from '../../config/router';
import { sendUserAlert } from '../../helpers/alert/sendUserAlert';
import { sendUserConfirm } from '../../helpers/alert/sendUserConfirm';
import { getSocketApiUrl } from '../../helpers/app/getSocketApiUrl';
import { useSocket } from '../../helpers/request/useSocket';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { useThrottleFnTrailing } from '../../helpers/utils/useThrottleFnTrailing';
import { useCredentialStore } from '../../stores/credential';
import { useLayoutStore } from '../../stores/layout';
import { useRetentionStore } from '../../stores/retention';

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
    const [activeChart, setActiveChart] = useState<string | null>(null); // Track the active chart

    const updateChannels = useCallback((channelData: Record<string, unknown>) => {
        const currentChannels = Object.keys(channelData);

        if (JSON.stringify(currentChannels) !== JSON.stringify(prevChannelsRef.current)) {
            setActiveChannels((prevChannels) => {
                const newChannels = { ...prevChannels };
                currentChannels.forEach((channel, index) => {
                    if (!newChannels[channel]) {
                        newChannels[channel] = {
                            id: `${RealTimeConstraints.id}_${channel}`,
                            index
                        };
                        waveformRefs.current[channel] =
                            createRef<DequeChartHandle>() as RefObject<DequeChartHandle>;
                        spectrogramRefs.current[channel] =
                            createRef<DequeSpectrogramHandle>() as RefObject<DequeSpectrogramHandle>;
                    }
                });
                Object.keys(newChannels).forEach((channel) => {
                    if (!currentChannels.includes(channel)) {
                        delete newChannels[channel];
                        delete waveformRefs.current[channel];
                        delete spectrogramRefs.current[channel];
                    }
                });
                prevChannelsRef.current = currentChannels;
                return newChannels;
            });
        }
    }, []);

    const { readyState, sendMessage } = useSocket(
        {
            url: getSocketApiUrl(getCredential().token),
            onData: ({ data }) => {
                const { channel_data, sample_rate, record_time, current_time } = data;

                Object.keys(channel_data).forEach((channel) => {
                    const waveformRef = waveformRefs.current[channel]?.current;
                    if (waveformRef) {
                        waveformRef.addData(
                            channel_data[channel].data_array,
                            record_time,
                            current_time,
                            sample_rate
                        );
                    }
                    const spectrogramRef = spectrogramRefs.current[channel]?.current;
                    if (spectrogramRef) {
                        spectrogramRef.addData(
                            channel_data[channel].data_array,
                            record_time,
                            current_time,
                            sample_rate
                        );
                    }
                });

                setSampleRate(sample_rate);
                setUpdatedAt(record_time);

                updateChannels(channel_data);
            }
        },
        true
    );
    useEffect(() => {
        if (readyState === 1) {
            sendMessage('client hello');
        }
    }, [readyState, sendMessage]);

    const getInitialLayout = useCallback(
        (id: string, index: number) => {
            if (config[id]?.position && config[id]?.size && config[id]?.spectrogram) {
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
                            ? RealTimeConstraints.minWidth * 2
                            : RealTimeConstraints.minWidth,
                    height:
                        document.documentElement.clientWidth > 768
                            ? RealTimeConstraints.minWidth * 2
                            : RealTimeConstraints.minHeight
                },
                spectrogram: { ...RealTimeConstraints.getDynamicDB(index) }
            };
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

    const handleToggleDisplayMode = useCallback(() => {
        setDisplayMode((prevMode) => (prevMode === 'waveform' ? 'spectrogram' : 'waveform'));
    }, []);

    const handleToggleRecording = useCallback(() => {
        setRecordingState((prev) => {
            const { isRecording, startTime } = prev;
            if (!isRecording) {
                if (updatedAt === 0) {
                    return prev;
                }
                sendUserAlert(
                    t('views.RealTime.record_data.start_recording', {
                        startedAt: getTimeString(updatedAt)
                    })
                );
                return { ...prev, isRecording: true, startTime: updatedAt };
            }

            const endTime = updatedAt;
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
    }, [t, updatedAt]);

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
                        const initialLayout = getInitialLayout(
                            activeChannels[channel].id,
                            activeChannels[channel].index
                        );
                        return (
                            <DraggableBox
                                key={channel}
                                layout={initialLayout}
                                locked={locks[RealTimeConstraints.id]}
                                constraints={RealTimeConstraints}
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
                                    <DequeChart
                                        minSpanValue={RealTimeConstraints.minSpanValue}
                                        ref={waveformRefs.current[channel]}
                                        lineColor={RealTimeConstraints.lineColor}
                                        maxDuration={retention}
                                        title={channel}
                                        height="100%"
                                        yPosition="right"
                                        zoom={true}
                                        animation={false}
                                    />
                                </div>
                                <div
                                    className={
                                        displayMode === 'spectrogram'
                                            ? 'block h-full w-full'
                                            : 'hidden'
                                    }
                                >
                                    <DequeSpectrogram
                                        title={channel}
                                        duration={retention}
                                        overlap={RealTimeConstraints.overlap}
                                        freqRange={RealTimeConstraints.freqRange}
                                        windowSize={RealTimeConstraints.windowSize}
                                        maxDB={initialLayout.spectrogram.maxDB}
                                        minDB={initialLayout.spectrogram.minDB}
                                        ref={spectrogramRefs.current[channel]}
                                        fftExecutor={sharedFFTExecutor}
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
        </div>
    );
};

export default RealTime;
