import ReactECharts from 'echarts-for-react';
import { memo, useEffect, useMemo, useRef } from 'react';

import { getChartAxis } from '../../styles/chartTheme';
import { useThemeColors } from '../../styles/useThemeColors';
import { getPressureTooltip } from './getPressureTooltip';
import { getLineSeries, type LineChartData } from './lineSeries';
import { HoverTimeHandler, useHoverTime } from './useHoverTime';
import { PlotBoundsHandler, usePlotBounds } from './usePlotBounds';
import { getPointsExtent, useVisibleSamples, VisibleSamplesHandler } from './useVisibleSamples';

interface ILineChart {
    readonly hidden?: boolean;
    readonly animation?: boolean;
    readonly lineColor?: string;
    readonly lineWidth?: number;
    readonly height?: number | string;
    readonly zoom?: boolean;
    readonly data: LineChartData;
    readonly unit?: string;
    readonly yMax?: number;
    readonly yMin?: number;
    readonly yInterval?: number;
    readonly yPosition?: 'left' | 'right';
    readonly minSpanValue?: number;
    readonly onVisibleSamples?: VisibleSamplesHandler;
    readonly onPlotBounds?: PlotBoundsHandler;
    // Reports the time under the pointer and its offset in the chart, or null on leave.
    readonly onHoverTime?: HoverTimeHandler;
}

export type { HoverTimeHandler } from './useHoverTime';

export const LineChart = memo(
    ({
        hidden,
        animation = true,
        lineColor,
        lineWidth = 2,
        height = 250,
        zoom = false,
        data,
        unit,
        yMax,
        yMin,
        yInterval,
        minSpanValue,
        yPosition = 'left',
        onVisibleSamples,
        onPlotBounds,
        onHoverTime
    }: ILineChart) => {
        const chartRef = useRef<ReactECharts>(null);
        const colors = useThemeColors();
        const seriesColor = lineColor ?? colors.line;
        const axis = useMemo(() => getChartAxis(colors), [colors]);
        const { onEvents, reset } = useVisibleSamples(() => {
            if (!Array.isArray(data[0])) {
                return null;
            }
            const points = data as [number, number | null][];
            return { points, extent: getPointsExtent(points) };
        }, onVisibleSamples);
        const plotEvents = usePlotBounds(chartRef, onPlotBounds);
        const events = useMemo(() => ({ ...onEvents, ...plotEvents }), [onEvents, plotEvents]);

        const option = useMemo(
            () => ({
                animation,
                backgroundColor: 'transparent',
                textStyle: { color: colors.label },
                xAxis: {
                    ...axis,
                    type: 'time',
                    axisLabel: {
                        ...axis.axisLabel,
                        hideOverlap: true,
                        formatter: (value: number) => {
                            const date = new Date(value);
                            const hours = date.getHours();
                            const minutes = date.getMinutes();
                            const seconds = date.getSeconds();
                            return `{normal|${hours}:${minutes < 10 ? '0' + minutes : minutes}:${seconds < 10 ? '0' + seconds : seconds}}`;
                        }
                    }
                },
                yAxis: {
                    ...axis,
                    type: 'value',
                    name: unit,
                    nameTextStyle: { color: colors.label, align: yPosition },
                    max: yMax,
                    min: yMin,
                    interval: yInterval,
                    position: yPosition,
                    scale: true
                },
                series: getLineSeries(Array.isArray(data) ? data : [], seriesColor, lineWidth),
                grid: { top: unit ? 28 : '5%', bottom: '3%', containLabel: true },
                tooltip: getPressureTooltip(unit),
                dataZoom: zoom
                    ? [
                          {
                              type: 'inside',
                              xAxisIndex: [0],
                              minValueSpan: minSpanValue,
                              zoomOnMouseWheel: true,
                              moveOnMouseMove: true
                          }
                      ]
                    : []
            }),
            [
                animation,
                unit,
                yMax,
                yMin,
                yInterval,
                yPosition,
                colors,
                axis,
                seriesColor,
                lineWidth,
                data,
                zoom,
                minSpanValue
            ]
        );

        useHoverTime(chartRef, onHoverTime);

        useEffect(() => {
            if (chartRef.current) {
                const instance = chartRef.current.getEchartsInstance();
                instance.setOption(option, true);
            }
            reset();
        }, [option, reset]);

        return (
            <ReactECharts
                ref={chartRef}
                option={option}
                className={`app-line-chart ${zoom ? 'app-draggable-cancel' : ''} ${hidden ? 'app-hidden' : ''}`}
                style={hidden ? undefined : { height }}
                opts={{ renderer: 'canvas' }}
                onEvents={events}
                notMerge={true}
                lazyUpdate={true}
            />
        );
    }
);
