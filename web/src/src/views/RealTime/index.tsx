import {
    mdiBlurLinear,
    mdiChartBellCurve,
    mdiClose,
    mdiLock,
    mdiLockOpen,
    mdiLockReset,
    mdiRecordCircle,
    mdiRecordRec,
    mdiWaveform
} from '@mdi/js';
import Icon from '@mdi/react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { FFTExecutor } from 'spectrogram-js';

import { ResizableWorkspace } from '../../components/ui/ResizableWorkspace';
import { WorkspaceExportButton } from '../../components/ui/WorkspaceExportButton';
import { Banner } from '../../components/widget/Banner';
import { RealTimeConstraints } from '../../config/constraints';
import { routerConfig } from '../../config/router';
import { SpectrumExecutor } from '../../helpers/dsp/SpectrumExecutor';
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
        groundMotion,
        getGroundMotionRecords,
        waveformRefs,
        spectrogramRefs,
        updatedAtRef,
        readyState
    } = useRealtimeStream();
    const {
        locks,
        toggleLock,
        fftModes,
        toggleFftMode,
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
    const sharedSpectrumExecutor = useMemo(() => new SpectrumExecutor(), []);
    useEffect(() => () => sharedSpectrumExecutor.dispose(), [sharedSpectrumExecutor]);
    const fftMode = Boolean(fftModes[RealTimeConstraints.id]);
    const [activeChart, setActiveChart] = useState<string | null>(null); // Track the active chart
    const workspaceRef = useRef<HTMLDivElement>(null);

    const handleActiveChart = useCallback((channel: string) => {
        setActiveChart(channel);
    }, []);

    const handleToggleFftMode = useCallback(
        () => toggleFftMode(RealTimeConstraints.id),
        [toggleFftMode]
    );

    const handleToggleDisplayMode = useCallback(() => {
        setDisplayMode((prevMode) => (prevMode === 'waveform' ? 'spectrogram' : 'waveform'));
    }, []);

    return (
        <div className="app-page">
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
            <div className="app-toolbar">
                <button
                    className="btn btn-sm app-button-toolbar"
                    onClick={() => toggleLock(RealTimeConstraints.id)}
                >
                    {locks[RealTimeConstraints.id] ? (
                        <>
                            <Icon className="app-icon" path={mdiLock} size={0.7} />
                            <span>{t('views.RealTime.layout_locker.unlock_button')}</span>
                        </>
                    ) : (
                        <>
                            <Icon className="app-icon" path={mdiLockOpen} size={0.7} />
                            <span>{t('views.RealTime.layout_locker.lock_button')}</span>
                        </>
                    )}
                </button>
                <button className="btn btn-sm app-button-toolbar" onClick={handleResetLayout}>
                    <Icon className="app-icon" path={mdiLockReset} size={0.7} />
                    <span>{t('views.RealTime.reset_layout.reset_button')}</span>
                </button>
                <button
                    className="btn btn-sm app-button-toolbar"
                    disabled={displayMode !== 'waveform'}
                    onClick={handleToggleFftMode}
                >
                    <Icon className="app-icon" path={mdiChartBellCurve} size={0.7} />
                    <span>
                        {fftMode
                            ? t('views.RealTime.fft_analysis.hide_button')
                            : t('views.RealTime.fft_analysis.show_button')}
                    </span>
                </button>
                {sampleRate >= RealTimeConstraints.freqRange[1] * 2 && (
                    <button
                        className="btn btn-sm app-button-toolbar"
                        onClick={handleToggleDisplayMode}
                    >
                        <Icon
                            className="app-icon"
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
                <WorkspaceExportButton
                    workspaceRef={workspaceRef}
                    fileName="realtime_workspace"
                    disabled={!Object.keys(activeChannels).length}
                />
                <button
                    className={`btn btn-sm app-button-toolbar ${recordingState.isRecording ? 'app-recording' : ''}`}
                    onClick={handleToggleRecording}
                >
                    <Icon
                        className="app-icon"
                        path={recordingState.isRecording ? mdiRecordRec : mdiRecordCircle}
                        size={0.7}
                    />
                    <span className={recordingState.isRecording ? 'app-pulse' : ''}>
                        {recordingState.isRecording
                            ? t('views.RealTime.record_data.stop_button')
                            : t('views.RealTime.record_data.start_button')}
                    </span>
                </button>
                <div className="app-recordings-status">
                    <div className="badge badge-soft badge-primary app-text-medium">
                        {t('views.RealTime.stream_status.sample_rate', { value: sampleRate })}
                    </div>
                    <div className="badge badge-soft badge-primary app-text-medium">
                        {t('views.RealTime.stream_status.channels', {
                            num: Object.keys(activeChannels).length
                        })}
                    </div>
                    <div className="badge badge-soft badge-primary app-text-medium">
                        {t('views.RealTime.stream_status.retention', { retention })}
                    </div>
                </div>
            </div>

            {recordList.length > 0 && (
                <div className="app-recording-files">
                    {recordList.map(([startTime, endTime, search], index) => (
                        <li className="app-recording-file">
                            <Link
                                key={index}
                                className="link link-primary app-recording-link"
                                target="_blank"
                                to={{ search, pathname: routerConfig.routes.history.uri }}
                            >
                                {`${getTimeString(startTime)} - ${getTimeString(endTime)}`}
                            </Link>
                            <button
                                className="app-icon-action"
                                onClick={() => handleRemoveRecord(index)}
                            >
                                <Icon className="app-icon" path={mdiClose} size={0.7} />
                            </button>
                        </li>
                    ))}
                </div>
            )}

            <ResizableWorkspace
                id={RealTimeConstraints.id}
                className="app-chart-workspace"
                containerRef={workspaceRef}
            >
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
                                fftMode={fftMode}
                                retention={retention}
                                sampleRate={sampleRate}
                                fftExecutor={sharedFFTExecutor}
                                spectrumExecutor={sharedSpectrumExecutor}
                                waveformRef={waveformRefs.current[channel]}
                                spectrogramRef={spectrogramRefs.current[channel]}
                                groundMotion={groundMotion[channel]}
                                getGroundMotionRecords={getGroundMotionRecords}
                                onActive={handleActiveChart}
                                onDragStop={handleDragStop}
                                onResizeStop={handleResizeStop}
                                onSpectrogram={handleSpectrogramUpdate}
                            />
                        );
                    })}
            </ResizableWorkspace>
        </div>
    );
};

export default RealTime;
