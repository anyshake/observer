import { wrap } from "comlink";
import { HighchartsReactRefObject } from "highcharts-react-official";
import { Dispatch, RefObject, SetStateAction } from "react";

import { ChartProps } from "../../components/Chart";
import { HolderProps } from "../../components/Holder";
import { CircularQueue2D } from "../../helpers/utils/CircularQueue2D";
import HandleSetChartsWorker, { api } from "../../workers/handleSetCharts.worker";
import { SocketUpdates } from "./getSocketUpdates";

const handleSetChartsWorkerApi = wrap<typeof api>(new HandleSetChartsWorker());

export const handleSetCharts = async (
	res: SocketUpdates,
	stateFn: Dispatch<
		SetStateAction<
			Record<
				string,
				{
					chart: ChartProps & {
						buffer: CircularQueue2D;
						ref: RefObject<HighchartsReactRefObject>;
						filter: {
							enabled: boolean;
							lowCorner?: number;
							highCorner?: number;
						};
					};
					holder: HolderProps & { values: Record<string, string> };
				}
			>
		>
	>
) => {
	stateFn((prev) => {
		const { timestamp, sample_rate, ...channels } = res;
		Object.keys(channels).forEach(async (key) => {
			// Validate if the buffer size is the same as the sample rate
			const [retention, sampleRate] = prev[key].chart.buffer.getShape();
			if (sampleRate !== sample_rate + 1) {
				prev[key].chart.buffer = new CircularQueue2D(
					retention,
					// plus one for timestamp
					sample_rate + 1
				);
			}
			// Put new data to circular buffer
			const channelData = new Float64Array([
				timestamp,
				...channels[key as keyof typeof channels]
			]);
			prev[key].chart.buffer.write(channelData);

			// Get filter settings and apply to chart
			const { enabled: filterEnabled, lowCorner, highCorner } = prev[key].chart.filter;
			const { lowFreqCorner, highFreqCorner } = {
				lowFreqCorner: lowCorner ?? 0.1,
				highFreqCorner: highCorner ?? 10
			};
			prev[key].chart = {
				...prev[key].chart,
				title: filterEnabled
					? `Band pass [${lowFreqCorner.toFixed(1)}-${highFreqCorner.toFixed(1)} Hz]`
					: ""
			};

			// Update chart from buffer, using web worker
			const { current: chartObj } = prev[key].chart.ref;
			if (chartObj) {
				const chartData = (
					await Promise.all(
						prev[key].chart.buffer
							.readAll()
							.map(
								async (data) =>
									await handleSetChartsWorkerApi.getChartAxisData(
										data,
										sample_rate,
										filterEnabled,
										4,
										lowFreqCorner,
										highFreqCorner
									)
							)
					)
				).flat();
				chartObj.chart.series[0].setData(chartData, true, false, false);
			}

			// Update label in holder component
			const { max, min } = await handleSetChartsWorkerApi.getLabelAxisValues(
				prev[key].chart.buffer.readAll()
			);
			prev[key].holder.values = { max, min };
		});

		return prev;
	});
};
