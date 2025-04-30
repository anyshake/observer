import ReactECharts from 'echarts-for-react';
import { memo, useEffect, useMemo, useRef } from 'react';

interface ILineChart {
    readonly title?: string;
    readonly animation?: boolean;
    readonly lineColor?: string;
    readonly lineWidth?: number;
    readonly height?: number | string;
    readonly zoom?: boolean;
    readonly data: [number, number | null][] | (number | null)[];
    readonly yMax?: number;
    readonly yMin?: number;
    readonly yInterval?: number;
    readonly yPosition?: 'left' | 'right';
    readonly minSpanValue?: number;
}

export const LineChart = memo(
    ({
        title,
        animation = true,
        lineColor = '#7e59ff',
        lineWidth = 2,
        height = 250,
        zoom = false,
        data,
        yMax,
        yMin,
        yInterval,
        minSpanValue,
        yPosition = 'left'
    }: ILineChart) => {
        const chartRef = useRef<ReactECharts>(null);

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
                        lineStyle: {
                            color: lineColor,
                            width: lineWidth
                        },
                        connectNulls: false,
                        data: Array.isArray(data) ? data : []
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
                title,
                animation,
                yMax,
                yMin,
                yInterval,
                yPosition,
                lineColor,
                lineWidth,
                data,
                zoom,
                minSpanValue
            ]
        );

        useEffect(() => {
            if (chartRef.current) {
                const instance = chartRef.current.getEchartsInstance();
                instance.setOption(option, true);
            }
        }, [option]);

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
);
