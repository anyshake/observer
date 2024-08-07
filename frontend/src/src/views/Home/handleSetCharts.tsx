import { HighchartsReactRefObject } from "highcharts-react-official";
import { Dispatch, RefObject, SetStateAction } from "react";

import { ChartProps } from "../../components/Chart";
import { HolderProps } from "../../components/Holder";
import { StationUpdates } from "./getStationUpdates";

const RETENTION_THRESHOLD_MS = 1000 * 60 * 5;

export const handleSetCharts = (
	res: StationUpdates,
	stateFn: Dispatch<
		SetStateAction<
			Record<
				string,
				{
					chart: ChartProps & { ref: RefObject<HighchartsReactRefObject> };
					holder: HolderProps & { values: Record<string, string> };
				}
			>
		>
	>
) => {
	if (!res?.data) {
		return;
	}

	const { timestamp } = res.data.os;
	stateFn((prev) => {
		// Set CPU usage chart
		const { percent: cpuPercent } = res.data.cpu;
		const { current: cpuChart } = prev.cpu.chart.ref;
		if (cpuChart) {
			const initTimestamp = cpuChart.chart.series[0].data.length
				? cpuChart.chart.series[0].data[0].x
				: timestamp;
			cpuChart.chart.series[0].addPoint(
				[timestamp, cpuPercent],
				true,
				timestamp - initTimestamp >= RETENTION_THRESHOLD_MS
			);
		}

		// Set memory usage chart
		const { percent: memoryPercent } = res.data.memory;
		const { current: memoryChart } = prev.memory.chart.ref;
		if (memoryChart) {
			const initTimestamp = memoryChart.chart.series[0].data.length
				? memoryChart.chart.series[0].data[0].x
				: timestamp;
			memoryChart.chart.series[0].addPoint(
				[timestamp, memoryPercent],
				true,
				timestamp - initTimestamp >= RETENTION_THRESHOLD_MS
			);
		}

		return prev;
	});
};
