import ReactECharts from 'echarts-for-react';
import { memo, useEffect, useMemo, useRef } from 'react';

import { getChartAxis, getChartTitle } from '../../styles/chartTheme';
import { useThemeColors } from '../../styles/useThemeColors';
import { getLineSeries, type LineChartData } from './lineSeries';

interface ILineChart {
    readonly hidden?: boolean;
    readonly title?: string;
    readonly animation?: boolean;
    readonly lineColor?: string;
    readonly lineWidth?: number;
    readonly height?: number | string;
    readonly zoom?: boolean;
    readonly data: LineChartData;
    readonly yMax?: number;
    readonly yMin?: number;
    readonly yInterval?: number;
    readonly yPosition?: 'left' | 'right';
    readonly minSpanValue?: number;
}

export const LineChart = memo(
    ({
        hidden,
        title,
        animation = true,
        lineColor,
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
        const colors = useThemeColors();
        const seriesColor = lineColor ?? colors.line;
        const axis = useMemo(() => getChartAxis(colors), [colors]);

        const option = useMemo(
            () => ({
                animation,
                backgroundColor: 'transparent',
                textStyle: { color: colors.label },
                title: getChartTitle(title, colors),
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
                    max: yMax,
                    min: yMin,
                    interval: yInterval,
                    position: yPosition,
                    scale: true
                },
                series: getLineSeries(Array.isArray(data) ? data : [], seriesColor, lineWidth),
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
                colors,
                axis,
                seriesColor,
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
                className={`app-line-chart ${hidden ? 'app-hidden' : ''}`}
                style={hidden ? undefined : { height }}
                opts={{ renderer: 'canvas' }}
                notMerge={true}
                lazyUpdate={true}
            />
        );
    }
);
