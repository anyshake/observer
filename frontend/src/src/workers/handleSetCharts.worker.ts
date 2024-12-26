import { expose } from "comlink";

import { FilterPassband, getFilteredCounts } from "../helpers/seismic/getFilteredCounts";

export default {} as typeof Worker & { new (): Worker };

export const api = {
	getChartAxisData: (
		bufferData: Float64Array,
		sampleRate: number,
		filterEnabled: boolean,
		poles: number,
		lowFreqCorner: number,
		highFreqCorner: number
	) => {
		const timestamp = bufferData[0];
		const channelData = Array.from(bufferData).slice(1);

		if (filterEnabled) {
			return Array.from(
				getFilteredCounts(Float32Array.from(channelData), {
					poles,
					lowFreqCorner,
					highFreqCorner,
					sampleRate,
					passbandType: FilterPassband.BAND_PASS
				})
			).map((value, index) => [timestamp + (1000 / sampleRate) * index, value]);
		}

		return channelData.map((value, index) => [timestamp + (1000 / sampleRate) * index, value]);
	},
	getLabelAxisValues: (bufferData: Float64Array[]) => {
		const channelData = [];
		for (const data of bufferData) {
			channelData.push(...Array.from(data).slice(1));
		}
		const max = Math.max(...channelData);
		const min = Math.min(...channelData);
		return { max: max.toFixed(0), min: min.toFixed(0) };
	}
};

expose(api);
