import { Dispatch, SetStateAction } from "react";

import { BannerProps } from "../../components/Banner";
import { getTimeString } from "../../helpers/utils/getTimeString";
import { SocketUpdates } from "./getSocketUpdates";

export const handleSetBanner = (
	res: SocketUpdates,
	adcResolution: number,
	sensorType: string,
	stateFn: Dispatch<SetStateAction<BannerProps & { values?: Record<string, string> }>>
) => {
	stateFn((prev) => ({
		...prev,
		type: "success",
		title: "views.realtime.banner.success.label",
		content: "views.realtime.banner.success.text",
		values: {
			sampleRate: String(res.sample_rate),
			adc_resolution: String(adcResolution),
			sensor_type: sensorType,
			time: getTimeString(res.timestamp)
		}
	}));
};
