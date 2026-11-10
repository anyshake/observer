import type { ECharts } from 'echarts';
import ReactECharts from 'echarts-for-react';
import { RefObject, useCallback, useEffect, useMemo, useRef } from 'react';

export interface IPlotBounds {
    readonly left: number;
    readonly right: number;
}

export type PlotBoundsHandler = (bounds: IPlotBounds) => void;

const GRID_FINDER = { gridIndex: 0 };

const findEdge = (instance: ECharts, y: number, inside: number, outside: number) => {
    let low = inside;
    let high = outside;
    while (Math.abs(high - low) > 0.5) {
        const middle = (low + high) / 2;
        if (instance.containPixel(GRID_FINDER, [middle, y])) {
            low = middle;
        } else {
            high = middle;
        }
    }
    return low;
};

const getPlotBounds = (instance: ECharts): IPlotBounds | null => {
    const width = instance.getWidth();
    const x = width / 2;
    const y = instance.getHeight() / 2;
    if (!instance.containPixel(GRID_FINDER, [x, y])) {
        return null;
    }
    return {
        left: Math.round(findEdge(instance, y, x, -1)),
        right: Math.round(width - findEdge(instance, y, x, width + 1))
    };
};

export const usePlotBounds = (
    chartRef: RefObject<ReactECharts | null>,
    onPlotBounds?: PlotBoundsHandler
) => {
    const handlerRef = useRef(onPlotBounds);
    handlerRef.current = onPlotBounds;

    const report = useCallback((instance: ECharts) => {
        const bounds = handlerRef.current && getPlotBounds(instance);
        if (bounds) {
            handlerRef.current?.(bounds);
        }
    }, []);

    useEffect(() => {
        if (onPlotBounds && chartRef.current) {
            report(chartRef.current.getEchartsInstance());
        }
    }, [chartRef, onPlotBounds, report]);

    return useMemo(
        () => ({ finished: (_: unknown, instance: ECharts) => report(instance) }),
        [report]
    );
};
