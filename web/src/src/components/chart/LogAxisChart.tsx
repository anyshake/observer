import ReactECharts from 'echarts-for-react';
import { memo, useMemo } from 'react';

import { getChartAxis } from '../../styles/chartTheme';
import { useThemeColors } from '../../styles/useThemeColors';
import { formatChartNumber } from './formatChartNumber';
import type { IPlotBounds } from './usePlotBounds';

interface ILogAxisChart {
    readonly height?: number | string;
    readonly data: [number, number | null][];
    readonly xName: string;
    readonly xMin?: number;
    readonly yName?: string;
    readonly yLog?: boolean;
    readonly yPosition?: 'left' | 'right';
    readonly bounds?: IPlotBounds | null;
}

export const LogAxisChart = memo(
    ({
        height = 300,
        data,
        xName,
        xMin,
        yName,
        yLog = false,
        yPosition = 'left',
        bounds
    }: ILogAxisChart) => {
        const colors = useThemeColors();
        const axis = useMemo(() => getChartAxis(colors), [colors]);

        const option = useMemo(
            () => ({
                animation: false,
                backgroundColor: 'transparent',
                textStyle: { color: colors.label },
                grid: bounds
                    ? { top: yName ? 36 : 16, bottom: 56, left: bounds.left, right: bounds.right }
                    : { top: yName ? 36 : 16, right: 24, bottom: 48, left: 12, containLabel: true },
                xAxis: {
                    ...axis,
                    type: 'log',
                    logBase: 10,
                    min: xMin ?? 'dataMin',
                    max: 'dataMax',
                    name: xName,
                    nameLocation: 'middle',
                    nameGap: 28,
                    axisLabel: {
                        ...axis.axisLabel,
                        hideOverlap: true,
                        formatter: formatChartNumber
                    },
                    splitLine: { ...axis.splitLine, show: true }
                },
                yAxis: {
                    ...axis,
                    type: yLog ? 'log' : 'value',
                    logBase: 10,
                    scale: !yLog,
                    position: yPosition,
                    name: yName,
                    nameLocation: 'end',
                    nameTextStyle: { color: colors.label, align: yPosition },
                    axisLabel: { ...axis.axisLabel, formatter: formatChartNumber }
                },
                series: [
                    {
                        type: 'line',
                        showSymbol: false,
                        connectNulls: false,
                        lineStyle: { color: colors.line, width: 1.5 },
                        itemStyle: { color: colors.line },
                        data
                    }
                ]
            }),
            [axis, bounds, colors, data, xMin, xName, yLog, yName, yPosition]
        );

        return (
            <ReactECharts
                option={option}
                className="app-line-chart"
                style={{ height }}
                opts={{ renderer: 'canvas' }}
                notMerge={true}
                lazyUpdate={true}
            />
        );
    }
);
