import { Dispatch, SetStateAction } from "react";

import { BannerProps } from "../../components/Banner";
import { StationUpdates } from "./getStationUpdates";

export const handleSetBanner = (
	res: StationUpdates,
	stateFn: Dispatch<SetStateAction<BannerProps & { values?: Record<string, string> }>>
) => {
	if (!res?.data) {
		stateFn({
			type: "error",
			title: "views.home.banner.error.title",
			content: "views.home.banner.error.content"
		});
		return;
	}
	const { station, explorer, os, disk } = res.data;
	stateFn({
		title: "views.home.banner.success.title",
		content: "views.home.banner.success.content",
		type: "success",
		values: {
			disk: (100 - disk.percent).toFixed(2),
			serial: `0x${explorer.device_id}`,
			uptime: String(os.uptime),
			station: station.name,
			arch: os.arch,
			os: os.os
		}
	});
};
