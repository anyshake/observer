import ReactECharts from 'echarts-for-react';
import { RefObject, useEffect, useRef } from 'react';

// Reports the time under the pointer and its offset in the chart, or null on leave.
export type HoverTimeHandler = (hover: { time: number; x: number; y: number } | null) => void;

export const useHoverTime = (
    chartRef: RefObject<ReactECharts | null>,
    onHoverTime?: HoverTimeHandler
) => {
    const hoverRef = useRef(onHoverTime);
    hoverRef.current = onHoverTime;
    const hasHover = !!onHoverTime;
    useEffect(() => {
        const instance = chartRef.current?.getEchartsInstance();
        if (!instance || !hasHover) {
            return;
        }
        const zr = instance.getZr();
        const handleMove = ({ offsetX, offsetY }: { offsetX: number; offsetY: number }) => {
            if (!instance.containPixel({ gridIndex: 0 }, [offsetX, offsetY])) {
                hoverRef.current?.(null);
                return;
            }
            const [time] = instance.convertFromPixel({ gridIndex: 0 }, [offsetX, offsetY]);
            hoverRef.current?.({ time, x: offsetX, y: offsetY });
        };
        const handleOut = () => hoverRef.current?.(null);
        zr.on('mousemove', handleMove);
        zr.on('globalout', handleOut);
        return () => {
            zr.off('mousemove', handleMove);
            zr.off('globalout', handleOut);
        };
    }, [chartRef, hasHover]);
};
