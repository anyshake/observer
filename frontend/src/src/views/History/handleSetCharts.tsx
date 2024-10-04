import { wrap } from "comlink";
import { HighchartsReactRefObject } from "highcharts-react-official";
import { Dispatch, RefObject, SetStateAction } from "react";

import { ChartProps } from "../../components/Chart";
import { HolderProps } from "../../components/Holder";
import { LabelProps } from "../../components/Label";
import { apiConfig } from "../../config/api";
import { CircularQueue2D } from "../../helpers/utils/CircularQueue2D";
import HandleSetChartsWorker, { api } from "../../workers/handleSetCharts.worker";

const handleSetChartsWorkerApi = wrap<typeof api>(new HandleSetChartsWorker());

export const handleSetCharts = (
	res:
		| typeof apiConfig.endpoints.history.model.response.common
		| typeof apiConfig.endpoints.history.model.response.error,
	chartStateFn: Dispatch<
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
					holder: HolderProps;
				}
			>
		>
	>,
	labelStateFn: Dispatch<
		SetStateAction<Record<string, LabelProps & { values?: Record<string, string> }>>
	>
) => {
	if (!res?.data?.length) {
		return;
	}

	chartStateFn((prev) => {
		const respData = res.data;
		const respDataDuration = respData.length;
		const respSampleRate = respData[0].sample_rate;

		Object.keys(prev).forEach(async (key) => {
			// Validate if the buffer size is the same as the sample rate
			const [dataDuration, sampleRate] = prev[key].chart.buffer.getShape();
			if (dataDuration !== respDataDuration || sampleRate !== respSampleRate + 1) {
				prev[key].chart.buffer = new CircularQueue2D(
					respDataDuration,
					// plus one for timestamp
					respSampleRate + 1
				);
			}

			// Put response data to circular buffer
			respData.forEach(({ timestamp, ...data }) => {
				const axisData = data[key as keyof typeof data] as number[];
				prev[key].chart.buffer.write(new Float64Array([timestamp, ...axisData]));
			});

			// Get filter settings and apply to chart
			const { enabled: filterEnabled, lowCorner, highCorner } = prev[key].chart.filter;
			const { lowFreqCorner, highFreqCorner } = {
				lowFreqCorner: lowCorner ?? 0.1,
				highFreqCorner: highCorner ?? 10
			};
			prev[key].chart = {
				...prev[key].chart,
				title: filterEnabled ? `Band pass [${lowFreqCorner}-${highFreqCorner} Hz]` : ""
			};

			// Get filtered values and apply to chart data
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
										respSampleRate,
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

			// Get label axis values and apply to label state
			const { max, min } = await handleSetChartsWorkerApi.getLabelAxisValues(
				prev[key].chart.buffer.readAll()
			);
			labelStateFn((_prev) => {
				const _prevCopy = { ..._prev };
				_prevCopy[key] = { ..._prev[key], values: { max, min } };
				return _prevCopy;
			});
		});

		return prev;
	});
};
