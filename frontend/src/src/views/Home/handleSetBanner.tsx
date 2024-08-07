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
	const {  name } = station;
    const { device_id } = explorer;
	const { arch, distro } = os;
	const { percent } = disk;
	stateFn({
		title: "views.home.banner.success.title",
		content: "views.home.banner.success.content",
		type: "success",
		values: {
			uptime: String(os.uptime),
			station: name,
			os: distro,
			arch,
			uuid: String(device_id),
			disk: (100 - percent).toFixed(2)
		}
	});
};
