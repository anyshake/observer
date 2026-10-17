import ReactECharts from 'echarts-for-react';
import { useEffect, useState } from 'react';

import { useThemeColors } from '../../styles/useThemeColors';

interface IGaugeChart {
    readonly animation?: boolean;
    readonly height?: number;
    readonly value: number;
    readonly valueMax?: number;
    readonly pointerColor?: string;
}

export const GaugeChart = ({
    animation,
    height,
    value,
    valueMax = 100,
    pointerColor
}: IGaugeChart) => {
    const [truncatedValue, setTruncatedValue] = useState(0);
    const colors = useThemeColors();
    const gaugeColor = pointerColor ?? colors.line;

    useEffect(() => {
        setTruncatedValue(Math.round(value * 100) / 100);
    }, [value]);

    return (
        <ReactECharts
            style={{ height: height ? `${height}px` : '100%' }}
            option={{
                animation,
                backgroundColor: 'transparent',
                textStyle: { color: colors.label },
                title: { x: 'center', textStyle: { color: colors.label } },
                series: [
                    {
                        type: 'gauge',
                        max: valueMax,
                        center: ['50%', '58%'],
                        startAngle: 215,
                        endAngle: -35,
                        data: [{ value: truncatedValue, itemStyle: { color: gaugeColor } }],
                        axisLine: { lineStyle: { color: [[1, colors.grid]] } },
                        axisLabel: { color: colors.label },
                        axisTick: { distance: -15, length: 5, lineStyle: { color: colors.axis } },
                        splitLine: { distance: -15, length: 10, lineStyle: { color: colors.axis } },
                        anchor: { itemStyle: { color: gaugeColor } },
                        pointer: { itemStyle: { color: gaugeColor }, length: '60%', width: 4 },
                        detail: {
                            formatter: '{value}%',
                            fontSize: 14,
                            color: colors.text
                        }
                    }
                ]
            }}
        />
    );
};
