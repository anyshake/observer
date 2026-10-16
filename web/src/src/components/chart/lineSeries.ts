export type LineChartData = [number, number | null][] | (number | null)[];

export const getLineSeries = (data: LineChartData, lineColor: string, lineWidth: number) => {
    const createSeries = (segment: LineChartData, index: number) => ({
        id: `line-segment-${index}`,
        type: 'line' as const,
        sampling: 'lttb' as const,
        showSymbol: false,
        lineStyle: { color: lineColor, width: lineWidth },
        connectNulls: false,
        data: segment
    });

    if (!Array.isArray(data[0])) {
        // Splitting scalar data would reset its implicit indices. Keep gaps without sampling.
        const hasGap = data.some((value) => value === null || !Number.isFinite(value));
        return [{ ...createSeries(data, 0), sampling: hasGap ? undefined : ('lttb' as const) }];
    }

    const segments: LineChartData[] = [];
    let segmentStart = 0;
    let inGap = false;

    for (let i = 0; i < data.length; i++) {
        const point = data[i];
        const value = Array.isArray(point) ? point[1] : point;
        if (value === null || !Number.isFinite(value)) {
            inGap = true;
            continue;
        }

        if (inGap && i > segmentStart) {
            // Separate series cannot reconnect when LTTB drops a gap marker.
            // Keep trailing nulls so their timestamps still contribute to the axis extent.
            segments.push(data.slice(segmentStart, i));
            segmentStart = i;
        }
        inGap = false;
    }

    segments.push(segmentStart === 0 ? data : data.slice(segmentStart));
    return segments.map(createSeries);
};
