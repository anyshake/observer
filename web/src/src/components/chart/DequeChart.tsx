import ReactECharts from 'echarts-for-react';
import {
    forwardRef,
    memo,
    useCallback,
    useEffect,
    useImperativeHandle,
    useMemo,
    useRef
} from 'react';

import TimeSeriesBuffer from '../../helpers/storage/TimeSeriesBuffer';
import { getChartAxis, getChartTitle } from '../../styles/chartTheme';
import { useThemeColors } from '../../styles/useThemeColors';
import { getLineSeries } from './lineSeries';

export interface DequeChartHandle {
    addData: (
        values: number[],
        recordTime: number,
        currentTime: number,
        sampleRate: number
    ) => void;
}

interface IDequeChart {
    readonly title?: string;
    readonly animation?: boolean;
    readonly lineColor?: string;
    readonly lineWidth?: number;
    readonly height?: number | string;
    readonly zoom?: boolean;
    readonly yMax?: number;
    readonly yMin?: number;
    readonly yInterval?: number;
    readonly yPosition?: 'left' | 'right';
    readonly minSpanValue?: number;
    readonly maxDuration?: number; // Buffer duration in seconds
    readonly paused?: boolean;
}

export const DequeChart = memo(
    forwardRef<DequeChartHandle, IDequeChart>(
        (
            {
                title,
                animation = true,
                lineColor,
                lineWidth = 2,
                height = 250,
                zoom = false,
                yMax,
                yMin,
                yInterval,
                yPosition = 'left',
                maxDuration = 60,
                minSpanValue,
                paused = false
            }: IDequeChart,
            ref
        ) => {
            const chartRef = useRef<ReactECharts>(null);
            const colors = useThemeColors();
            const seriesColor = lineColor ?? colors.line;
            const axis = useMemo(() => getChartAxis(colors), [colors]);
            const seriesStyleRef = useRef({ color: seriesColor, width: lineWidth });
            seriesStyleRef.current = { color: seriesColor, width: lineWidth };
            const bufferRef = useRef(new TimeSeriesBuffer(maxDuration));
            const rafRef = useRef<number | null>(null);
            const needsUpdateRef = useRef(false);
            const pausedRef = useRef(paused);
            const maxDurationRef = useRef(maxDuration);
            const resizeOnDrawRef = useRef(false);
            const scheduleRef = useRef<() => void>(() => {});
            pausedRef.current = paused;
            maxDurationRef.current = maxDuration;

            scheduleRef.current = () => {
                if (rafRef.current !== null) {
                    return;
                }
                rafRef.current = requestAnimationFrame(() => {
                    rafRef.current = null;
                    if (pausedRef.current || !needsUpdateRef.current) {
                        return;
                    }
                    if (!chartRef.current) {
                        scheduleRef.current();
                        return;
                    }

                    needsUpdateRef.current = false;
                    const instance = chartRef.current.getEchartsInstance();
                    if (resizeOnDrawRef.current) {
                        resizeOnDrawRef.current = false;
                        instance.resize();
                    }
                    const data = bufferRef.current.getData();
                    const endTime = bufferRef.current.getEndTime();
                    const startTime = endTime - maxDurationRef.current * 1000;
                    instance.setOption(
                        {
                            series: getLineSeries(
                                data,
                                seriesStyleRef.current.color,
                                seriesStyleRef.current.width
                            ),
                            xAxis: { min: startTime, max: endTime }
                        },
                        { replaceMerge: ['series'] }
                    );
                    if (needsUpdateRef.current && !pausedRef.current) {
                        scheduleRef.current();
                    }
                });
            };

            const addData = useCallback(
                (values: number[], recordTime: number, currentTime: number, sampleRate: number) => {
                    bufferRef.current.addData(values, recordTime, currentTime, sampleRate);
                    needsUpdateRef.current = true;
                    if (!pausedRef.current) {
                        scheduleRef.current();
                    }
                },
                []
            );

            useImperativeHandle(ref, () => ({
                addData
            }));

            const wasPausedRef = useRef(paused);
            useEffect(
                () => () => {
                    if (rafRef.current !== null) {
                        cancelAnimationFrame(rafRef.current);
                        rafRef.current = null;
                    }
                },
                []
            );
            useEffect(() => {
                if (paused) {
                    wasPausedRef.current = true;
                    if (rafRef.current !== null) {
                        cancelAnimationFrame(rafRef.current);
                        rafRef.current = null;
                    }
                    return;
                }

                const resumed = wasPausedRef.current;
                wasPausedRef.current = false;
                if (!resumed) {
                    return;
                }

                resizeOnDrawRef.current = true;
                needsUpdateRef.current = true;
                scheduleRef.current();
                return () => {
                    if (rafRef.current !== null) {
                        cancelAnimationFrame(rafRef.current);
                        rafRef.current = null;
                    }
                };
            }, [paused]);

            const option = useMemo(() => {
                // Theme changes repaint buffered samples without resetting the live stream.
                const data = bufferRef.current.getData();
                const endTime = data.length ? bufferRef.current.getEndTime() : Date.now() + 1000;
                return {
                    animation,
                    backgroundColor: 'transparent',
                    textStyle: { color: colors.label },
                    title: getChartTitle(title, colors),
                    xAxis: {
                        ...axis,
                        type: 'time',
                        min: data.length ? endTime - maxDurationRef.current * 1000 : endTime - 1000,
                        max: endTime,
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
                        max: yMax,
                        min: yMin,
                        interval: yInterval,
                        position: yPosition,
                        scale: true
                    },
                    series: getLineSeries(data, seriesColor, lineWidth),
                    grid: { top: '5%', bottom: '3%', containLabel: true },
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
                };
            }, [
                animation,
                title,
                yMax,
                yMin,
                yInterval,
                yPosition,
                colors,
                axis,
                seriesColor,
                lineWidth,
                zoom,
                minSpanValue
            ]);

            return (
                <ReactECharts
                    ref={chartRef}
                    option={option}
                    className="app-line-chart"
                    style={{ height }}
                    opts={{ renderer: 'canvas' }}
                    notMerge={true}
                    lazyUpdate={true}
                />
            );
        }
    )
);
