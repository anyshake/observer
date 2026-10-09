import { memo, RefObject, useCallback } from 'react';
import type { ColorMapName } from 'spectrogram-js';
import { FFTExecutor } from 'spectrogram-js';

import { DequeChart, DequeChartHandle } from '../../components/chart/DequeChart';
import { DequeSpectrogram, DequeSpectrogramHandle } from '../../components/chart/DequeSpectrogram';
import { DraggableBox } from '../../components/ui/DraggableBox';
import { RealTimeConstraints } from '../../config/constraints';
import { LayoutConfig } from '../../stores/layout';

interface IRealTimeChannel {
    readonly channel: string;
    readonly channelId: string;
    readonly index: number;
    readonly layout: LayoutConfig;
    readonly locked: boolean;
    readonly waveformMode: boolean;
    readonly retention: number;
    readonly sampleRate: number;
    readonly fftExecutor: FFTExecutor;
    readonly waveformRef?: RefObject<DequeChartHandle>;
    readonly spectrogramRef?: RefObject<DequeSpectrogramHandle>;
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
    retention,
    sampleRate,
    fftExecutor,
    waveformRef,
    spectrogramRef,
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

    return (
        <DraggableBox
            layout={layout}
            locked={locked}
            constraints={RealTimeConstraints}
            onDragStart={handleDragStart}
            onDragStop={handleDragStop}
            onResizeStop={handleResizeStop}
        >
            <div className={waveformMode ? 'block h-full w-full' : 'hidden'}>
                <DequeChart
                    minSpanValue={RealTimeConstraints.minSpanValue}
                    ref={waveformRef}
                    lineColor={RealTimeConstraints.lineColor}
                    maxDuration={retention}
                    title={channel}
                    height="100%"
                    yPosition="right"
                    zoom={true}
                    animation={false}
                    paused={!waveformMode}
                />
            </div>
            <div className={waveformMode ? 'hidden' : 'block h-full w-full'}>
                <DequeSpectrogram
                    title={channel}
                    duration={retention}
                    overlap={RealTimeConstraints.overlap}
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
        </DraggableBox>
    );
});
