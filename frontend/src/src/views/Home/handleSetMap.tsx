import { Dispatch, SetStateAction } from "react";

import { HolderProps } from "../../components/Holder";
import { MapBoxProps } from "../../components/MapBox";
import { StationUpdates } from "./getStationUpdates";

export const handleSetMap = (
	res: StationUpdates,
	stateFn: Dispatch<
		SetStateAction<{
			mapbox: MapBoxProps;
			holder: HolderProps & { values: Record<string, string> };
		}>
	>
) => {
	if (!res?.data) {
		return;
	}
	const { explorer } = res.data;
	const { longitude, latitude, elevation } = explorer;
	stateFn((prev) => ({
		...prev,
		mapbox: {
			...prev.mapbox,
			center: [latitude, longitude],
			marker: [latitude, longitude]
		},
		holder: {
			...prev.holder,
			values: {
				...prev.holder.values,
				elevation: elevation.toFixed(2),
				latitude: latitude.toFixed(2),
				longitude: longitude.toFixed(2)
			}
		}
	}));
};
