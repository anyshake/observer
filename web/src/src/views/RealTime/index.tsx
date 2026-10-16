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
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { FFTExecutor } from 'spectrogram-js';

import { Banner } from '../../components/widget/Banner';
import { RealTimeConstraints } from '../../config/constraints';
import { routerConfig } from '../../config/router';
import { getTimeString } from '../../helpers/utils/getTimeString';
import { useRetentionStore } from '../../stores/retention';
import { RealtimeChannel } from './RealtimeChannel';
import { useRealtimeLayout } from './useRealtimeLayout';
import { useRealtimeStream } from './useRealtimeStream';
import { useRecording } from './useRecording';

const RealTime = () => {
    const { t } = useTranslation();

    const { retention } = useRetentionStore();

    const {
        sampleRate,
        updatedAt,
        activeChannels,
        waveformRefs,
        spectrogramRefs,
        updatedAtRef,
        readyState
    } = useRealtimeStream();
    const {
        locks,
        toggleLock,
        getInitialLayout,
        handleResetLayout,
        handleDragStop,
        handleResizeStop,
        handleSpectrogramUpdate
    } = useRealtimeLayout(activeChannels);
    const { recordingState, recordList, handleToggleRecording, handleRemoveRecord } =
        useRecording(updatedAtRef);
    const [displayMode, setDisplayMode] = useState<'waveform' | 'spectrogram'>('waveform');
    useEffect(() => {
        // Make sure sample rate supports spectrogram display
        if (sampleRate < RealTimeConstraints.freqRange[1] * 2) {
            setDisplayMode('waveform');
        }
    }, [sampleRate]);

    const sharedFFTExecutor = useMemo(() => new FFTExecutor(RealTimeConstraints.fftSize), []);
    const [activeChart, setActiveChart] = useState<string | null>(null); // Track the active chart

    const handleActiveChart = useCallback((channel: string) => {
        setActiveChart(channel);
    }, []);

    const handleToggleDisplayMode = useCallback(() => {
        setDisplayMode((prevMode) => (prevMode === 'waveform' ? 'spectrogram' : 'waveform'));
    }, []);

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
