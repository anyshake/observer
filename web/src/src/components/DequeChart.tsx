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

import TimeSeriesBuffer from '../helpers/storage/TimeSeriesBuffer';

export interface DequeChartHandle {
    addData: (values: number[], timestamp: number, sampleRate: number) => void;
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
}

export const DequeChart = memo(
    forwardRef<DequeChartHandle, IDequeChart>(
        (
            {
                title,
                animation = true,
                lineColor = '#7e59ff',
                lineWidth = 2,
                height = 250,
                zoom = false,
                yMax,
                yMin,
                yInterval,
                yPosition = 'left',
                maxDuration = 60,
                minSpanValue
            }: IDequeChart,
            ref
        ) => {
            const chartRef = useRef<ReactECharts>(null);
            const bufferRef = useRef(new TimeSeriesBuffer(maxDuration));
            const rafRef = useRef<number | null>(null);
            const needsUpdateRef = useRef(false);

            const addData = useCallback(
                (values: number[], timestamp: number, sampleRate: number) => {
                    bufferRef.current.addData(values, timestamp, sampleRate);
                    needsUpdateRef.current = true;
                },
                []
            );

            useImperativeHandle(ref, () => ({
                addData
            }));

            const updateChart = useCallback(() => {
                if (needsUpdateRef.current && chartRef.current) {
                    needsUpdateRef.current = false;
                    const instance = chartRef.current.getEchartsInstance();
                    const data = bufferRef.current.getData();
                    const endTime = bufferRef.current.getEndTime();
                    const startTime = endTime - maxDuration * 1000;
                    instance.setOption({
                        series: [{ data }],
                        xAxis: { min: startTime, max: endTime }
                    });
                }
                rafRef.current = requestAnimationFrame(updateChart);
            }, [maxDuration]);

            useEffect(() => {
                rafRef.current = requestAnimationFrame(updateChart);
                return () => {
                    if (rafRef.current) {
                        cancelAnimationFrame(rafRef.current);
                    }
                };
            }, [updateChart]);

            const option = useMemo(
                () => ({
                    animation,
                    title: title
                        ? {
                              text: title,
                              left: 'left',
                              textStyle: {
                                  color: '#4A4A4A',
                                  fontSize: 13,
                                  fontWeight: 'bold'
                              },
                              backgroundColor: '#f0f0f0',
                              padding: [6, 12],
                              borderRadius: 3
                          }
                        : {},
                    xAxis: {
                        type: 'time',
                        min: Date.now(),
                        max: Date.now() + 1000,
                        axisLabel: {
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
                        type: 'value',
                        max: yMax,
                        min: yMin,
                        interval: yInterval,
                        position: yPosition,
                        scale: true
                    },
                    series: [
                        {
                            type: 'line',
                            sampling: 'lttb',
                            showSymbol: false,
                            lineStyle: { color: lineColor, width: lineWidth },
                            connectNulls: false,
                            data: []
                        }
                    ],
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
                }),
                [
                    animation,
                    title,
                    yMax,
                    yMin,
                    yInterval,
                    yPosition,
                    lineColor,
                    lineWidth,
                    zoom,
                    minSpanValue
                ]
            );

            return (
                <ReactECharts
                    ref={chartRef}
                    option={option}
                    style={{
                        height: typeof height === 'number' ? `${height}px` : height,
                        width: '100%'
                    }}
                    opts={{ renderer: 'canvas' }}
                    notMerge={true}
                    lazyUpdate={true}
                />
            );
        }
    )
);
