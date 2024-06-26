import { Dispatch, SetStateAction } from "react";

import { LabelProps } from "../../components/Label";
import { StationUpdates } from "./getStationUpdates";

export const handleSetLabels = (
	res: StationUpdates,
	stateFn: Dispatch<SetStateAction<Record<string, LabelProps>>>
) => {
	if (!res?.data) {
		return;
	}
	const { status } = res.data;
	stateFn((prev) => {
		Object.keys(status).forEach((key) => {
			if (key in prev) {
				const newValue = status[key as keyof typeof status];
				prev[key] = { ...prev[key], value: String(newValue) };
			}
		});

		return prev;
	});
};
