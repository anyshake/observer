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
					chart: ChartProps & {
						ref: RefObject<HighchartsReactRefObject>;
					};
					holder: HolderProps & { values: Record<string, string> };
				}
			>
		>
	>
) => {
	if (!res?.data) {
		return;
	}
	stateFn((prev) => {
		const { data } = res;
		const { timestamp } = data;
		Object.keys(prev).forEach((key) => {
			if (!(key in data)) {
				return;
			}

			if (Object.prototype.hasOwnProperty.call(data[key as keyof typeof data], "percent")) {
				// Get percentage value by key in state
				const { percent } = data[key as keyof typeof data] as {
					percent: number;
				};
				const { current: chart } = prev[key].chart.ref;
				if (chart) {
					// Append new data to buffer and remove expired data
					const initTimestamp = chart.chart.series[0].data.length
						? chart.chart.series[0].data[0].x
						: timestamp;
					chart.chart.series[0].addPoint(
						[timestamp, percent],
						true,
						timestamp - initTimestamp >= RETENTION_THRESHOLD_MS
					);
				}

				prev[key] = {
					...prev[key],
					holder: {
						...prev[key].holder,
						values: { usage: percent.toFixed(2) }
					}
				};
			}
		});

		return prev;
	});
};
