import { useCallback, useEffect, useMemo, useRef } from 'react';

import { SpectrumConstraints } from '../../config/constraints';

export type VisibleSamplesHandler = (samples: Float64Array) => void;

export interface ISampleSource {
    readonly points: readonly [number, number | null][];
    readonly extent: readonly [number, number];
}

interface IZoomRange {
    readonly start?: number;
    readonly end?: number;
}

interface IZoomEvent extends IZoomRange {
    readonly batch?: readonly IZoomRange[];
}

export const getPointsExtent = (points: readonly [number, number | null][]) => {
    let min = Infinity;
    let max = -Infinity;
    for (const [time] of points) {
        min = Math.min(min, time);
        max = Math.max(max, time);
    }
    return [min, max] as const;
};

const sliceSamples = ({ points, extent }: ISampleSource, [start, end]: [number, number]) => {
    const span = extent[1] - extent[0];
    const from = extent[0] + (span * start) / 100;
    const to = extent[0] + (span * end) / 100;
    const values: number[] = [];
    for (const [time, value] of points) {
        if (time >= from && time <= to && value !== null && Number.isFinite(value)) {
            values.push(value);
        }
    }
    return Float64Array.from(values);
};

export const useVisibleSamples = (
    getSource: () => ISampleSource | null,
    onVisibleSamples?: VisibleSamplesHandler
) => {
    const percentRef = useRef<[number, number]>([0, 100]);
    const timerRef = useRef<number | undefined>(undefined);
    const getSourceRef = useRef(getSource);
    const handlerRef = useRef(onVisibleSamples);
    getSourceRef.current = getSource;
    handlerRef.current = onVisibleSamples;

    const notify = useCallback(() => {
        window.clearTimeout(timerRef.current);
        if (!handlerRef.current) {
            return;
        }
        timerRef.current = window.setTimeout(() => {
            const source = getSourceRef.current();
            if (source) {
                handlerRef.current?.(sliceSamples(source, percentRef.current));
            }
        }, SpectrumConstraints.debounceDelay);
    }, []);

    const reset = useCallback(() => {
        percentRef.current = [0, 100];
        notify();
    }, [notify]);

    const onEvents = useMemo(
        () => ({
            datazoom: ({ batch, start, end }: IZoomEvent) => {
                const range = batch?.[0] ?? { start, end };
                if (range.start === undefined || range.end === undefined) {
                    return;
                }
                percentRef.current = [range.start, range.end];
                notify();
            }
        }),
        [notify]
    );

    useEffect(() => {
        notify();
    }, [notify, onVisibleSamples]);
    useEffect(() => () => window.clearTimeout(timerRef.current), []);

    return { onEvents, notify, reset };
};
