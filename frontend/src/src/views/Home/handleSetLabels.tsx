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

	const { errors, received, elapsed } = res.data.explorer;
	stateFn((prev) => {
		prev.errors = { ...prev.errors, value: String(errors) };
		prev.messages = { ...prev.messages, value: String(received) };
		prev.elapsed = { ...prev.elapsed, value: String(elapsed) };

		return prev;
	});
};
