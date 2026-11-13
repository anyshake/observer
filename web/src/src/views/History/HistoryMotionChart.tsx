import { memo, useCallback } from 'react';

import type { IGroundMotionRecord } from '../../components/chart/groundMotion';
import { GroundMotionHover } from '../../components/chart/GroundMotionHover';
import { LineChart } from '../../components/chart/LineChart';
import type { LineChartData } from '../../components/chart/lineSeries';
import type { PlotBoundsHandler } from '../../components/chart/usePlotBounds';
import type { VisibleSamplesHandler } from '../../components/chart/useVisibleSamples';
import { HistoryConstraints } from '../../config/constraints';

interface IHistoryMotionChart {
    readonly data: LineChartData;
    readonly groundMotion?: IGroundMotionRecord[];
    readonly onVisibleSamples?: VisibleSamplesHandler;
    readonly onPlotBounds?: PlotBoundsHandler;
}

// LineChart with an overlay showing the PGA/PGV/PGD of the second under the pointer.
export const HistoryMotionChart = memo(function HistoryMotionChart({
    data,
    groundMotion,
    onVisibleSamples,
    onPlotBounds
}: IHistoryMotionChart) {
    const getRecords = useCallback(() => groundMotion, [groundMotion]);
    return (
        <GroundMotionHover getRecords={getRecords}>
            {(onHoverTime) => (
                <LineChart
                    minSpanValue={HistoryConstraints.minSpanValue}
                    height="100%"
                    lineWidth={1}
                    yPosition="right"
                    zoom={true}
                    animation={false}
                    data={data}
                    onVisibleSamples={onVisibleSamples}
                    onPlotBounds={onPlotBounds}
                    onHoverTime={groundMotion?.length ? onHoverTime : undefined}
                />
            )}
        </GroundMotionHover>
    );
});
