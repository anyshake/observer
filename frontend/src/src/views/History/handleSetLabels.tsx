import { Dispatch, SetStateAction } from "react";

import { LabelProps } from "../../components/Label";
import { apiConfig } from "../../config/api";
import store from "../../config/store";
import { getAccelerationArr } from "../../helpers/seismic/getAccelerationArr";
import { getVelocityArr } from "../../helpers/seismic/getVelocityArr";
import { getVoltageArr } from "../../helpers/seismic/getVoltageArr";

export const handleSetLabels = (
	res:
		| typeof apiConfig.endpoints.history.model.response.common
		| typeof apiConfig.endpoints.history.model.response.error,
	stateFn: Dispatch<
		SetStateAction<Record<string, LabelProps & { values?: Record<string, string> }>>
	>
) => {
	if (!res?.data) {
		return;
	}

	const { adc } = store.getState().adc;
	const { geophone } = store.getState().geophone;

	stateFn((prev) => {
		Object.keys(prev).forEach((key) => {
			if (!res.data.every((obj) => key in obj)) {
				return;
			}
			const channelDataArr = res.data.map(
				(obj) => obj[key as keyof typeof obj]
			) as number[][];
			const voltageDataArr = channelDataArr.map((arr) => {
				return getVoltageArr(arr, adc.resolution, adc.fullscale);
			});
			const velocityDataArr = voltageDataArr.map((arr) => {
				const sensitivity = geophone.sensitivity / 100;
				return getVelocityArr(arr, sensitivity);
			});
			const accelerationDataArr = velocityDataArr.map((arr) => {
				const channelDataSpanMS = 1000 / arr.length;
				return getAccelerationArr(arr, channelDataSpanMS);
			});

			const pgv = velocityDataArr
				.flat()
				.reduce((a, b) => Math.max(Math.abs(a), Math.abs(b)), 0);
			const pga = accelerationDataArr
				.flat()
				.reduce((a, b) => Math.max(Math.abs(a), Math.abs(b)), 0);

			prev[key] = {
				...prev[key],
				values: {
					pgv: pgv.toFixed(5),
					pga: pga.toFixed(5)
				},
				value: `views.history.labels.${key}_detail.value`
			};
		});

		return prev;
	});
};
