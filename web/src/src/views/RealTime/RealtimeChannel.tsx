import { mdiChartLine, mdiGradientVertical } from '@mdi/js';
import { memo, RefObject, useCallback } from 'react';
import type { ColorMapName } from 'spectrogram-js';
import { FFTExecutor } from 'spectrogram-js';

import { ChartPanel } from '../../components/chart/ChartPanel';
import { DequeChart, DequeChartHandle } from '../../components/chart/DequeChart';
import { DequeSpectrogram, DequeSpectrogramHandle } from '../../components/chart/DequeSpectrogram';
import type { IGroundMotion, IGroundMotionRecord } from '../../components/chart/groundMotion';
import { GroundMotionHover } from '../../components/chart/GroundMotionHover';
import { GroundMotionValues } from '../../components/chart/GroundMotionValues';
import { PressureValues } from '../../components/chart/PressureValues';
import { SpectrumSplit } from '../../components/chart/SpectrumSplit';
import { DraggableBox } from '../../components/ui/DraggableBox';
import { RealTimeConstraints } from '../../config/constraints';
import { SpectrumExecutor } from '../../helpers/dsp/SpectrumExecutor';
import { LayoutConfig } from '../../stores/layout';

interface IRealTimeChannel {
    readonly channel: string;
    readonly channelId: string;
    readonly index: number;
    readonly layout: LayoutConfig;
    readonly locked: boolean;
    readonly waveformMode: boolean;
    readonly fftMode: boolean;
    readonly retention: number;
    readonly sampleRate: number;
    readonly fftExecutor: FFTExecutor;
    readonly spectrumExecutor: SpectrumExecutor;
    readonly waveformRef?: RefObject<DequeChartHandle>;
    readonly spectrogramRef?: RefObject<DequeSpectrogramHandle>;
    readonly groundMotion?: IGroundMotion;
    readonly unit?: string;
    readonly latestValue?: number;
    readonly getGroundMotionRecords: (channel: string) => IGroundMotionRecord[] | undefined;
    readonly onActive: (channel: string) => void;
    readonly onDragStop: (channel: string, index: number, x: number, y: number) => void;
    readonly onResizeStop: (channel: string, index: number, width: number, height: number) => void;
    readonly onSpectrogram: (
        channel: string,
        index: number,
        minDB: number,
        maxDB: number,
        colorMap: ColorMapName
    ) => void;
}

export const RealtimeChannel = memo(function RealtimeChannel({
    channel,
    channelId,
    index,
    layout,
    locked,
    waveformMode,
    fftMode,
    retention,
    sampleRate,
    fftExecutor,
    spectrumExecutor,
    waveformRef,
    spectrogramRef,
    groundMotion,
    unit,
    latestValue,
    getGroundMotionRecords,
    onActive,
    onDragStop,
    onResizeStop,
    onSpectrogram
}: IRealTimeChannel) {
    const handleDragStart = useCallback(() => onActive(channel), [channel, onActive]);
    const handleDragStop = useCallback(
        (x: number, y: number) => onDragStop(channelId, index, x, y),
        [channelId, index, onDragStop]
    );
    const handleResizeStop = useCallback(
        (width: number, height: number) => onResizeStop(channelId, index, width, height),
        [channelId, index, onResizeStop]
    );
    const handleSpectrogram = useCallback(
        (minDB: number, maxDB: number, colorMap: ColorMapName) =>
            onSpectrogram(channelId, index, minDB, maxDB, colorMap),
        [channelId, index, onSpectrogram]
    );

    const getRecords = useCallback(
        () => getGroundMotionRecords(channel),
        [channel, getGroundMotionRecords]
    );

    return (
        <DraggableBox
            layout={layout}
            locked={locked}
            constraints={RealTimeConstraints}
            onDragStart={handleDragStart}
            onDragStop={handleDragStop}
            onResizeStop={handleResizeStop}
        >
            <ChartPanel
                title={channel}
                icon={waveformMode ? mdiChartLine : mdiGradientVertical}
                extra={
                    unit ? (
                        <PressureValues value={latestValue} unit={unit} />
                    ) : groundMotion ? (
                        <GroundMotionValues motion={groundMotion} />
                    ) : undefined
                }
            >
                <div className={waveformMode ? 'app-chart-visible' : 'app-hidden'}>
                    <SpectrumSplit
                        enabled={fftMode && waveformMode}
                        sampleRate={sampleRate}
                        executor={spectrumExecutor}
                        unit={unit}
                    >
                        {(onVisibleSamples, onPlotBounds) => (
                            <GroundMotionHover getRecords={getRecords}>
                                {(onHoverTime) => (
                                    <DequeChart
                                        key={unit || 'counts'}
                                        minSpanValue={RealTimeConstraints.minSpanValue}
                                        ref={waveformRef}
                                        maxDuration={retention}
                                        height="100%"
                                        lineWidth={1}
                                        yPosition="right"
                                        zoom={true}
                                        animation={false}
                                        paused={!waveformMode}
                                        unit={unit}
                                        onVisibleSamples={onVisibleSamples}
                                        onPlotBounds={onPlotBounds}
                                        onHoverTime={groundMotion ? onHoverTime : undefined}
                                    />
                                )}
                            </GroundMotionHover>
                        )}
                    </SpectrumSplit>
                </div>
                <div className={waveformMode ? 'app-hidden' : 'app-chart-visible'}>
                    <DequeSpectrogram
                        duration={retention}
                        overlap={RealTimeConstraints.overlap}
                        gapThreshold={RealTimeConstraints.gapThreshold}
                        freqRange={RealTimeConstraints.freqRange}
                        windowSize={RealTimeConstraints.windowSize}
                        maxDB={layout.spectrogram.maxDB}
                        minDB={layout.spectrogram.minDB}
                        colorMap={layout.spectrogram.colorMap}
                        ref={spectrogramRef}
                        fftExecutor={fftExecutor}
                        sampleRate={sampleRate}
                        paused={waveformMode}
                        onSpectrogramUpdate={handleSpectrogram}
                    />
                </div>
            </ChartPanel>
        </DraggableBox>
    );
});
