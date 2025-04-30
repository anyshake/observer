import ReactECharts from 'echarts-for-react';
import { useEffect, useState } from 'react';

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
    pointerColor = '#c23531'
}: IGaugeChart) => {
    const [truncatedValue, setTruncatedValue] = useState(0);

    useEffect(() => {
        setTruncatedValue(Math.round(value * 100) / 100);
    }, [value]);

    return (
        <ReactECharts
            style={{ height: height ? `${height}px` : '100%' }}
            option={{
                animation,
                title: { x: 'center', textStyle: { color: '#666' } },
                series: [
                    {
                        type: 'gauge',
                        max: valueMax,
                        center: ['50%', '58%'],
                        startAngle: 215,
                        endAngle: -35,
                        data: [{ value: truncatedValue, itemStyle: { color: pointerColor } }],
                        axisTick: { distance: -15, length: 5 },
                        splitLine: { distance: -15, length: 10 },
                        pointer: { itemStyle: { color: pointerColor }, length: '60%', width: 4 },
                        detail: {
                            formatter: '{value}%',
                            fontSize: 14,
                            color: '#1e2939'
                        }
                    }
                ]
            }}
        />
    );
};
